package main

import (
	"context"
	"fmt"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Sampler: CustomSampler{Fallback: sdktrace.AlwaysSample()},

		// Pretty print spans to stdout. For debugging purposes.
		PrettyPrint: true,
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	tracer := otel.Tracer("app_or_package_name")

	{
		ctx, trace1 := tracer.Start(ctx, "trace1")

		_, span := tracer.Start(ctx, "child1")
		span.End()

		trace1.End()
		fmt.Printf("trace1: %s\n", upclient.TraceURL(span))
	}

	{
		ctx, trace2 := tracer.Start(ctx, "trace2")

		_, span := tracer.Start(ctx, "child1")
		span.End()

		trace2.End()
		fmt.Printf("trace2: %s\n", upclient.TraceURL(span))
	}
}

// CustomSampler drops some traces based on their name and uses fallback sampler otherwise.
type CustomSampler struct {
	Fallback sdktrace.Sampler
}

func (s CustomSampler) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if params.Name == "trace2" {
		// Drop traces with such name.
		return sdktrace.SamplingResult{
			Decision: sdktrace.Drop,
		}
	}

	// For the rest use fallback sampler.
	return s.Fallback.ShouldSample(params)
}

func (s CustomSampler) Description() string {
	return s.Fallback.Description()
}
