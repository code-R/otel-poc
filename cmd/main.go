package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"otel-poc/utils"

	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(semconv.HTTPRouteKey.String("hello"))

	_, span = utils.Tracer().Start(ctx, "my-internal-operation")
	defer span.End()

	span.RecordError(errors.New("ooooooops"), trace.WithAttributes(
		attribute.String("failed here", "just testing"),
	))
	span.SetStatus(codes.Error, "Error occurred")
	time.Sleep(30 * time.Millisecond)

	span.AddEvent("writing response", trace.WithAttributes(
		attribute.String("body", "hello world"),
	))
	w.Write([]byte("Hello, world!"))
}

func installPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

func main() {
	installPropagators()
	ctx := context.Background()
	gcpProject := "tilda-trial-dev"
	// fmt.Println(common.Version)
	fmt.Println(utils.ServiceName)
	tp, err := utils.InitTracing(ctx, utils.Config{
		ProjectID:      gcpProject,
		ServiceName:    utils.ServiceName,
		ServiceVersion: utils.Version,
	})

	if err != nil {
		fmt.Println(err)
	}
	defer tp.Shutdown(ctx)

	// define router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware(utils.ServiceName, otelchi.WithChiRoutes(r)))
	r.HandleFunc("/hello", helloHandler)
	// r.HandleFunc("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write(([]byte)("Hello, world..!"))
	// }))

	// otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello")
	// http.Handle("/hello", otelHandler)
	// http.ListenAndServe(":9000", nil)
	err = http.ListenAndServe("localhost:9000", r)
	if err != nil {
		log.Fatalf("unable to execute server due: %v", err)
	}
}
