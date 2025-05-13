package embeddings

import (
	"context"
	"fmt"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type OpenTelemetryEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	embeddings.DelegatingEmbeddingGenerator[TInput, TEmbedding]
	Metadata      *embeddings.EmbeddingGeneratorMetadata
	system        string
	serverAddress string
	serverPort    string
	modelId       string
}

const (
	tracerName = "gen-ai-tracer"
	metricName = "gen-ai-metrics"
)

var (
	tracer = otel.Tracer(tracerName)
	meter  = otel.Meter(metricName)
)
var latencyHist, _ = meter.Float64Histogram(
	"gen_ai.embeddings.operation.duration",
	metric.WithUnit("token"),
	metric.WithDescription("Measures the duration of a GenAI operation"),
	metric.WithExplicitBucketBoundaries([]float64{0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64, 1.28, 2.56, 5.12, 10.24, 20.48, 40.96, 81.92}...),
)

var tokenHist, _ = meter.Int64Histogram(
	"gen_ai.embeddings.token.usage",
	metric.WithUnit("token"),
	metric.WithDescription("Measures number of input and output tokens used"),
	metric.WithExplicitBucketBoundaries([]float64{1, 4, 16, 64, 256, 1_024, 4_096, 16_384, 65_536, 262_144, 1_048_576, 4_194_304, 16_777_216, 67_108_864}...),
)

func NewOpenTelemetryEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
	metadata *embeddings.EmbeddingGeneratorMetadata,
) *OpenTelemetryEmbeddingGenerator[TInput, TEmbedding] {
	system := "gen-ai"
	serverAddress := "localhost"
	serverPort := "0"
	modelId := ""

	if metadata != nil {
		if metadata.DefaultModelId != nil {
			modelId = *metadata.DefaultModelId
		}
		if metadata.ProviderName != nil {
			system = *metadata.ProviderName
		}
		if metadata.ProviderUri != nil {
			serverAddress = metadata.ProviderUri.Hostname()
			serverPort = metadata.ProviderUri.Port()
		}
	}
	return &OpenTelemetryEmbeddingGenerator[TInput, TEmbedding]{
		DelegatingEmbeddingGenerator: *embeddings.NewDelegatingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
		Metadata:                     metadata,
		system:                       system,
		serverAddress:                serverAddress,
		serverPort:                   serverPort,
		modelId:                      modelId,
	}
}

func (c *OpenTelemetryEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[TEmbedding], error) {
	ctx, span := CreateAndConfigureSpan(ctx, options, c.system, c.serverAddress, c.modelId, c.serverPort)
	startTime := time.Now()

	response, err := c.DelegatingEmbeddingGenerator.Generate(ctx, values, options)

	TraceResponse(ctx, span, response, err, startTime)
	return response, err
}

func CreateAndConfigureSpan(
	ctx context.Context,
	options *embeddings.EmbeddingGenerationOptions,
	system string,
	serverAddress string,
	serverPort string,
	modelId string,
) (context.Context, trace.Span) {
	// add Span name
	if options != nil && options.ModelId != nil {
		modelId = *options.ModelId
	}
	spanName := "gen-ai-embeddings"
	if modelId != "" {
		spanName = fmt.Sprintf("%s %s", spanName, modelId)
	}

	ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))

	// add base Tags
	span.SetAttributes(
		attribute.String("gen_ai.operation.name", "embeddings"),
		attribute.String("gen_ai.request.model", modelId),
		attribute.String("gen_ai.system_name", system),
	)

	if serverAddress != "" {
		span.SetAttributes(
			attribute.String("server.address", serverAddress),
			attribute.String("server.port", serverPort),
		)
	}

	// handle embeddings.EmbeddingGenerationOptions
	if options != nil {
		if options.Dimensions != nil {
			span.SetAttributes(attribute.Int64("gen_ai.request.embedding.dimensions", *options.Dimensions))
		}

		// log AdditionalProperties
		for key, value := range options.AdditionalProperties {
			span.SetAttributes(attribute.String(fmt.Sprintf("gen_ai.request.%s.%s", system, key), fmt.Sprintf("%v", value)))
		}
	}

	return ctx, span
}

// TraceResponse log request/reponse
func TraceResponse[TEmbedding embeddings.IEmbedding](ctx context.Context, span trace.Span, response *embeddings.GeneratedEmbeddings[TEmbedding], err error, startTime time.Time) {
	// log time metric
	duration := time.Since(startTime).Seconds()
	latencyHist.Record(ctx, duration)

	// log Token sum
	if response != nil && response.Usage != nil {
		if response.Usage.InputTokenCount != nil && *response.Usage.InputTokenCount > 0 {
			tokenHist.Record(ctx, *response.Usage.InputTokenCount, metric.WithAttributes(attribute.String("gen_ai.token.type", "input")))
		}
		if response.Usage.OutputTokenCount != nil && *response.Usage.OutputTokenCount > 0 {
			tokenHist.Record(ctx, *response.Usage.OutputTokenCount, metric.WithAttributes(attribute.String("gen_ai.token.type", "output")))
		}
	}

	// if err set status
	if err != nil {
		span.SetAttributes(attribute.String("error.type", fmt.Sprintf("%T", err)))
		span.SetStatus(1, err.Error())
		span.RecordError(err)
	}

	// log Response Tags
	if response != nil {
		if response.Usage != nil {
			if response.Usage.InputTokenCount != nil {
				span.SetAttributes(
					attribute.Int64("gen_ai.response.input_tokens", *response.Usage.InputTokenCount),
				)
			}
			if response.Usage.OutputTokenCount != nil {
				span.SetAttributes(
					attribute.Int64("gen_ai.response.output_tokens", *response.Usage.OutputTokenCount),
				)
			}
		}

		// log AdditionalProperties
		for key, value := range response.AdditionalProperties {
			span.SetAttributes(attribute.String(fmt.Sprintf("gen_ai.response.%s.%s", tracerName, key), fmt.Sprintf("%v", value)))
		}
	}

	span.End()
}
