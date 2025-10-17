// Package attr defines standard attribute key names for structured logging.
// These constants ensure consistent field naming across the application and
// follow OpenTelemetry semantic conventions where applicable.
package attr

// Standard attribute key names for slog.Attr.
// These keys follow naming conventions and best practices for structured logging.
const (
	// Address represents network addresses (IP addresses, hostnames, etc.)
	Address = "address"

	// DurationMs represents request duration in milliseconds
	DurationMs = "duration_ms"

	// Error represents error objects or error messages
	Error = "error"

	// Headers represents grouped HTTP headers
	Headers = "headers"

	// Method represents HTTP methods, RPC methods, or other operation types
	Method = "method"

	// RemoteAddr represents client IP addresses from headers
	RemoteAddr = "remote_addr"

	// Request represents request objects or request identifiers
	Request = "request"

	// RPC represents RPC procedure names
	RPC = "rpc"

	// SpanID represents OpenTelemetry span identifiers.
	// Follows OpenTelemetry semantic conventions: https://opentelemetry.io/docs/specs/semconv/general/naming/
	SpanID = "span_id"

	// Status represents response status codes or error states
	Status = "status"

	// TraceID represents OpenTelemetry trace identifiers.
	// Follows OpenTelemetry semantic conventions: https://opentelemetry.io/docs/specs/semconv/general/naming/
	TraceID = "trace_id"

	// UserAgent represents client User-Agent header values
	UserAgent = "user_agent"
)
