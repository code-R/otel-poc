package utils

import (
	"context"
	"fmt"
	"log"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ProjectID      string
	ServiceName    string
	ServiceVersion string
}

func InitTracing(ctx context.Context, cfg Config) (func(), error) {
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(cfg.ProjectID))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return func() {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down trace provider: %+v", err)
		}
	}, nil
}

func Tracer() trace.Tracer {
	return otel.Tracer("main")
}
