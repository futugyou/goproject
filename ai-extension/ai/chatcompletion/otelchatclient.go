package chatcompletion

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "gen-ai-tracer"
	metricName = "gen-ai-metrics"
)

var (
	tracer = otel.Tracer(tracerName)
	meter  = otel.Meter(metricName)
)
var latencyHist, _ = meter.Float64Histogram(
	"gen_ai.client.operation.duration",
	metric.WithUnit("token"),
	metric.WithDescription("Measures the duration of a GenAI operation"),
	metric.WithExplicitBucketBoundaries([]float64{0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64, 1.28, 2.56, 5.12, 10.24, 20.48, 40.96, 81.92}...),
)

var tokenHist, _ = meter.Int64Histogram(
	"gen_ai.client.token.usage",
	metric.WithUnit("token"),
	metric.WithDescription("Measures number of input and output tokens used"),
	metric.WithExplicitBucketBoundaries([]float64{1, 4, 16, 64, 256, 1_024, 4_096, 16_384, 65_536, 262_144, 1_048_576, 4_194_304, 16_777_216, 67_108_864}...),
)

type OpenTelemetryChatClient struct {
	chatcompletion.DelegatingChatClient
	Metadata *chatcompletion.ChatClientMetadata
}

func NewOpenTelemetryChatClient(innerClient chatcompletion.IChatClient, metadata *chatcompletion.ChatClientMetadata) *OpenTelemetryChatClient {
	return &OpenTelemetryChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{InnerClient: innerClient},
		Metadata:             metadata,
	}
}

func (c *OpenTelemetryChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	system := "gen-ai"
	serverAddress := "localhost"
	serverPort := "0"
	modelId := ""

	if c.Metadata != nil {
		if c.Metadata.ModelId != nil {
			modelId = *c.Metadata.ModelId
		}
		if c.Metadata.ProviderName != nil {
			system = *c.Metadata.ProviderName
		}
		if c.Metadata.ProviderUri != nil {
			if u, err := url.Parse(*c.Metadata.ProviderUri); err == nil {
				serverAddress = u.Hostname()
				serverPort = u.Port()
			}
		}
	}

	ctx, span := CreateAndConfigureSpan(ctx, options, system, serverAddress, modelId, serverPort)
	startTime := time.Now()
	defer span.End()

	response, err := c.InnerClient.GetResponse(ctx, chatMessages, options)
	TraceResponse(ctx, span, response, err, startTime)
	return response, err
}

func (c *OpenTelemetryChatClient) GetStreamingResponse(
	ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	system := "gen-ai"
	serverAddress := "localhost"
	serverPort := "0"
	modelId := ""

	if c.Metadata != nil {
		if c.Metadata.ModelId != nil {
			modelId = *c.Metadata.ModelId
		}
		if c.Metadata.ProviderName != nil {
			system = *c.Metadata.ProviderName
		}
		if c.Metadata.ProviderUri != nil {
			if u, err := url.Parse(*c.Metadata.ProviderUri); err == nil {
				serverAddress = u.Hostname()
				serverPort = u.Port()
			}
		}
	}

	// create OpenTelemetry span
	ctx, span := CreateAndConfigureSpan(ctx, options, system, serverAddress, modelId, serverPort)
	startTime := time.Now()

	// ploxy InnerClient stream response
	responseChan := c.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
	outputChan := make(chan chatcompletion.ChatStreamingResponse)

	go func() {
		defer close(outputChan)
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %v", r)
				TraceResponse(ctx, span, nil, err, startTime)
				span.End()
			}
		}()

		var updates []chatcompletion.ChatResponseUpdate
		var responseError error

		for response := range responseChan {
			if response.Err != nil {
				responseError = response.Err
			}

			if responseError != nil {
				break
			}

			updates = append(updates, *response.Update)
			select {
			case outputChan <- response:
			case <-ctx.Done():
				responseError = ctx.Err()
			}
		}

		// end OpenTelemetry tracing
		response := chatcompletion.ToChatResponse(updates)
		TraceResponse(ctx, span, &response, responseError, startTime)
		span.End()
	}()

	return outputChan
}

func CreateAndConfigureSpan(
	ctx context.Context,
	options *chatcompletion.ChatOptions,
	system string,
	serverAddress string,
	serverPort string,
	modelId string,
) (context.Context, trace.Span) {
	// add Span name
	if options != nil && options.ModelId != nil {
		modelId = *options.ModelId
	}
	spanName := "gen-ai-chat"
	if modelId != "" {
		spanName = fmt.Sprintf("%s %s", spanName, modelId)
	}

	ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))

	// add base Tags
	span.SetAttributes(
		attribute.String("gen_ai.operation.name", "chat"),
		attribute.String("gen_ai.request.model", modelId),
		attribute.String("gen_ai.system_name", system),
	)

	if serverAddress != "" {
		span.SetAttributes(
			attribute.String("server.address", serverAddress),
			attribute.String("server.port", serverPort),
		)
	}

	// handle chatcompletion.ChatOptions
	if options != nil {
		if options.FrequencyPenalty != nil {
			span.SetAttributes(attribute.Float64("gen_ai.request.frequency_penalty", *options.FrequencyPenalty))
		}
		if options.MaxOutputTokens != nil {
			span.SetAttributes(attribute.Int("gen_ai.request.max_tokens", *options.MaxOutputTokens))
		}
		if options.PresencePenalty != nil {
			span.SetAttributes(attribute.Float64("gen_ai.request.presence_penalty", *options.PresencePenalty))
		}
		if options.Seed != nil {
			span.SetAttributes(attribute.Int64("gen_ai.request.seed", *options.Seed))
		}
		if len(options.StopSequences) > 0 {
			span.SetAttributes(attribute.String("gen_ai.request.stop_sequences", fmt.Sprintf("[%q]", strings.Join(options.StopSequences, "\", \""))))
		}
		if options.Temperature != nil {
			span.SetAttributes(attribute.Float64("gen_ai.request.temperature", *options.Temperature))
		}
		if options.TopK != nil {
			span.SetAttributes(attribute.Int("gen_ai.request.top_k", *options.TopK))
		}
		if options.TopP != nil {
			span.SetAttributes(attribute.Float64("gen_ai.request.top_p", *options.TopP))
		}

		// handle ResponseFormat
		if options.ResponseFormat != nil {
			format := *options.ResponseFormat
			span.SetAttributes(attribute.String(fmt.Sprintf("gen_ai.request.%s.response_format", system), string(format)))
		}

		// log AdditionalProperties
		for key, value := range options.AdditionalProperties {
			span.SetAttributes(attribute.String(fmt.Sprintf("gen_ai.request.%s.%s", system, key), fmt.Sprintf("%v", value)))
		}
	}

	return ctx, span
}

// TraceResponse log request/reponse
func TraceResponse(ctx context.Context, span trace.Span, response *chatcompletion.ChatResponse, err error, startTime time.Time) {
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
		if response.FinishReason != nil {
			reason := *response.FinishReason
			span.SetAttributes(attribute.String("gen_ai.response.finish_reason", fmt.Sprintf("[\"%s\"]", strings.ToLower(string(reason)))))
		}
		if response.ResponseId != nil {
			span.SetAttributes(attribute.String("gen_ai.response.id", *response.ResponseId))
		}
		if response.ModelId != nil {
			span.SetAttributes(attribute.String("gen_ai.response.model", *response.ModelId))
		}
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
