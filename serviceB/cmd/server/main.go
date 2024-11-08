package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/felipehrs/goexpert-labs-otel-serciceB/internal/weather/handler"
	"github.com/felipehrs/goexpert-labs-otel-serciceB/internal/weather/usecase"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"
)

func initProvider(serviceName, collectorURL string, zipkinURL string) (func(context.Context) error, error) {
	ctx := context.Background()
	client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(collectorURL), otlptracehttp.WithInsecure())
	exporter, err := otlptrace.New(ctx, client)

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %v", err)
	}

	zipkinExporter, err := zipkin.New(zipkinURL + "/api/v2/spans")
	if err != nil {
		return nil, fmt.Errorf("failed to create zipkin exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithBatcher(zipkinExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp.Shutdown, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider("SERVICE_B", "otel-collector:4317", "http://zipkin:9411")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	r := gin.Default()
	r.Use(otelgin.Middleware("Service A"))

	usecase := usecase.NewWeatherUsecase()
	handler := handler.NewWeatherHandler(usecase)

	r.GET("/weather/:zipcode", handler.Handle)

	go func() {
		r.Run(":8081")
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
