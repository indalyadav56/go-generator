package tracer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	ErrEmptyServiceName    = errors.New("service name cannot be empty")
	ErrEmptyServiceVersion = errors.New("service version cannot be empty")
	ErrEmptyCollectorURL   = errors.New("collector URL cannot be empty")
	ErrInitTimeout         = errors.New("tracer initialization timed out")
)

type Config struct {
	ServiceName    string        `json:"service_name" yaml:"service_name"`
	ServiceVersion string        `json:"service_version" yaml:"service_version"`
	Environment    string        `json:"environment" yaml:"environment"`
	CollectorURL   string        `json:"collector_url" yaml:"collector_url"`
	Insecure       bool          `json:"insecure" yaml:"insecure"`
	InitTimeout    time.Duration `json:"init_timeout" yaml:"init_timeout"`
}

type Tracer interface {
	TracerProvider() *sdktrace.TracerProvider
	Shutdown(ctx context.Context) error
}

func DefaultConfig() Config {
	return Config{
		Environment: "development",
		Insecure:    false,
		InitTimeout: 5 * time.Second,
	}
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	if c.ServiceName == "" {
		return ErrEmptyServiceName
	}
	if c.ServiceVersion == "" {
		return ErrEmptyServiceVersion
	}
	if c.CollectorURL == "" {
		return ErrEmptyCollectorURL
	}
	return nil
}

// Provider wraps the OpenTelemetry TracerProvider
type Provider struct {
	tp *sdktrace.TracerProvider
}

func New(cfg Config) (*Provider, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.InitTimeout)
	defer cancel()

	tp, err := initTracerProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Provider{tp: tp}, nil
}

// Shutdown gracefully shuts down the tracer provider
func (p *Provider) Shutdown(ctx context.Context) error {
	if err := p.tp.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down tracer provider: %w", err)
	}
	return nil
}

// TracerProvider returns the underlying OpenTelemetry TracerProvider
func (p *Provider) TracerProvider() *sdktrace.TracerProvider {
	return p.tp
}

// initTracerProvider initializes the OpenTelemetry TracerProvider
func initTracerProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	exporter, err := createExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global TracerProvider and TextMapPropagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// createExporter creates and configures the OTLP exporter
func createExporter(ctx context.Context, cfg Config) (*otlptrace.Exporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.CollectorURL),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	return exporter, nil
}

// createResource creates the OpenTelemetry resource with service information
func createResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	return res, nil
}
