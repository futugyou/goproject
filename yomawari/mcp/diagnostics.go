package mcp

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	tracerName = "Experimental.ModelContextProtocol"
	metricName = "Experimental.ModelContextProtocol"
)

var (
	Tracer                       = otel.Tracer(tracerName)
	meter                        = otel.Meter(metricName)
	ShortSecondsBucketBoundaries = []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10}
	LongSecondsBucketBoundaries  = []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 30, 60, 120, 300}
)

func CreateDurationHistogram(name string, description string, longBuckets bool) metric.Float64Histogram {
	boundaries := ShortSecondsBucketBoundaries
	if longBuckets {
		boundaries = LongSecondsBucketBoundaries
	}
	m, _ := meter.Float64Histogram(
		name,
		metric.WithUnit("s"),
		metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(boundaries...),
	)
	return m
}
