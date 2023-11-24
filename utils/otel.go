package utils

import (
	"context"
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

func InitTracing(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	// Create exporter.
	// projectID := "tilda-trial-dev"
	// exporter, err := texporter.New(texporter.WithProjectID(cfg.ProjectID))
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(cfg.ProjectID))
	if err != nil {
		return nil, err
	}
	// if err != nil {
	// 	log.Fatalf("texporter.New: %v", err)
	// }

	// Identify your application using resource detection
	res, err := resource.New(ctx,
		// Use the GCP resource detector to detect information about the GCP platform
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// probabilistic sampling.
	// Example:
	//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	// defer tp.Shutdown(ctx) // flushes any pending spans, and closes connections.
	otel.SetTracerProvider(tp)
	return tp, nil
}

func Tracer() trace.Tracer {
	return otel.Tracer("main")
}
