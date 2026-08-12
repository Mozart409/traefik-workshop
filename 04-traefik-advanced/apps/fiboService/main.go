package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func Fibonacci(n int) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}
	var n2, n1 uint64 = 0, 1
	for i := int(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	return n2 + n1, nil
}

// initTracer sets up the OpenTelemetry tracer provider with OTLP HTTP exporter.
// It reads configuration from environment variables:
// - OTEL_EXPORTER_OTLP_ENDPOINT: Jaeger collector endpoint (default: http://localhost:4318)
// - OTEL_SERVICE_NAME: Service name (default: fiboservice)
func initTracer(ctx context.Context) (*trace.TracerProvider, error) {
	// Create OTLP HTTP exporter. Configuration is read from environment variables.
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Build a resource with service name.
	// Merge with resource.Default() to pick up environment variables like OTEL_SERVICE_NAME.
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "fiboservice"
	}

	// NewSchemaless avoids a schema-URL conflict: resource.Default() carries the
	// SDK's own schema version, and merging it with a differently-versioned
	// semconv resource is a hard error.
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create the tracer provider with the OTLP exporter and resource.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return tp, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize tracer provider.
	tp, err := initTracer(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize tracer: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}()

	// Set the global tracer provider.
	otel.SetTracerProvider(tp)

	// Set the global propagator to handle W3C traceparent headers from Traefik.
	// This is critical: without this, spans will appear as orphan traces instead of
	// nesting under Traefik's span.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	app := fiber.New()

	// Register the OTel Fiber middleware to trace HTTP requests.
	app.Use(otelfiber.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World from Fiber!")
	})

	app.Get("/fibonacci/:n", func(c *fiber.Ctx) error {
		n, err := c.ParamsInt("n")
		if err != nil {
			return c.Status(400).SendString("Invalid parameter")
		}

		// Start a child span around the Fibonacci computation.
		tracer := otel.Tracer("fiboService")
		_, span := tracer.Start(c.UserContext(), "Fibonacci")
		defer span.End()

		// Add attributes to the span for input and output.
		span.SetAttributes(attribute.Int("fibonacci.input", n))

		fib, err := Fibonacci(n)
		if err != nil {
			return c.Status(500).SendString("Internal server error")
		}

		// Record the result as a span attribute.
		span.SetAttributes(attribute.Int64("fibonacci.result", int64(fib)))

		return c.SendString(fmt.Sprintf("%d", fib))
	})

	app.Listen(":3333")
}
