package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initResources() *resource.Resource {
	res, _ := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("backend-service"),
		),
	)
	return res
}

func initTracer(res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

func initMetrics(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

func main() {
	res := initResources()

	// Initialize Tracing
	tp, _ := initTracer(res)
	defer tp.Shutdown(context.Background())

	// Initialize Metrics
	mp, _ := initMetrics(res)
	defer mp.Shutdown(context.Background())

	// Define Business Meter and Counter
	meter := otel.Meter("backend-operations")
	jobCounter, _ := meter.Int64Counter("processed_jobs_total",
		metric.WithDescription("Total number of business jobs successfully processed"),
	)

	// Routing
	mux := http.NewServeMux()

	// New "Chaos" Route for Senior-Level testing
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a 20% failure rate to test our Grafana alerts later
		if rand.Intn(5) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Critical Business Failure"))
			return
		}

		// Increment our business metric on success
		jobCounter.Add(r.Context(), 1)
		w.Write([]byte("Business Job Completed Successfully"))
	})

	// Wrap handler for automatic OTel instrumentation
	otelHandler := otelhttp.NewHandler(mux, "http-server")

	log.Println("Backend active on :8080. Test the chaos route at http://localhost:8080/work")
	http.ListenAndServe(":8080", otelHandler)
}
