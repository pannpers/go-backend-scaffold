package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace" // Import the OTEL trace package
)

// Following https://opentelemetry.io/docs/specs/semconv/general/naming/.
const (
	// traceIDKey is the context key for trace ID.
	traceIDKey string = "trace_id"
	// spanIDKey is the context key for span ID.
	spanIDKey string = "span_id"
)

// Logger is a structured logger using slog.
type Logger struct {
	logger *slog.Logger
}

// New creates a new Logger with the given options.
func New(opts ...Option) *Logger {
	// Start with default options.
	o := defaultOptions()

	// Apply all options.
	for _, opt := range opts {
		opt(o)
	}

	handlerOpts := &slog.HandlerOptions{
		Level:       o.level,
		ReplaceAttr: o.replaceAttrFunc,
	}

	var handler slog.Handler
	switch o.format {
	case FormatText:
		handler = slog.NewTextHandler(o.writer, handlerOpts)
	case FormatJSON:
		fallthrough
	default:
		handler = slog.NewJSONHandler(o.writer, handlerOpts)
	}

	logger := slog.New(handler)

	return &Logger{
		logger: logger,
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelError, msg, args...)
}

// With returns a logger with the given attributes.
func (l *Logger) With(args ...slog.Attr) *Logger {
	slogArgs := make([]any, len(args))
	for i, v := range args {
		slogArgs[i] = v
	}
	return &Logger{
		logger: l.logger.With(slogArgs...),
	}
}

// log is the internal logging method that handles context.
func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...slog.Attr) {
	// Extract trace and span IDs from context.
	contextAttrs := fromContext(ctx)

	// Combine context attributes with provided args.
	allArgs := make([]slog.Attr, 0, len(contextAttrs)+len(args))
	allArgs = append(allArgs, contextAttrs...)
	allArgs = append(allArgs, args...)

	// Log with the combined attributes.
	l.logger.LogAttrs(ctx, level, msg, allArgs...)
}

// fromContext extracts trace and span IDs from context using OpenTelemetry.
func fromContext(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	// Get the SpanContext from the context
	spanContext := trace.SpanFromContext(ctx).SpanContext()

	// If the SpanContext is not valid, there are no trace IDs to add.
	if !spanContext.IsValid() {
		return attrs
	}
	// If spanContext.IsValid() is true, then TraceID and SpanID are also valid and non-zero.
	attrs = append(attrs, slog.String(traceIDKey, spanContext.TraceID().String()))
	attrs = append(attrs, slog.String(spanIDKey, spanContext.SpanID().String()))

	return attrs
}
