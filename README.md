# l4g - Lightweight Logging for Go

English | [简体中文](README.zh-CN.md)

A high-performance, structured logging library for Go that's compatible with the standard `log/slog` package. Designed for speed, simplicity, and zero allocations for disabled log levels.

## Features

- **Fast & Efficient**: Zero allocations for disabled log levels, buffer pooling for minimal GC pressure
- **Structured Logging**: Full support for key-value pairs and attributes
- **slog Compatible**: Works seamlessly with Go's standard `log/slog` package
- **Multiple Log Levels**: Trace, Debug, Info, Warn, Error, Panic, Fatal
- **Colorized Output**: Optional ANSI color support for terminal output
- **Named Channels**: Create multiple independent loggers with different configurations
- **Thread-Safe**: Built with concurrency in mind using `sync.Map` and atomic operations
- **Flexible Handlers**: Customizable log formatting and output destinations
- **Printf-style & JSON**: Support for formatted strings and structured JSON-like logging

## Installation

```bash
go get go-slim.dev/l4g
```

Requires Go 1.24.0 or later.

## Quick Start

```go
package main

import "go-slim.dev/l4g"

func main() {
    // Use the default logger
    l4g.Info("Hello, World!")
    l4g.Infof("User %s logged in", "alice")

    // Structured logging with key-value pairs
    l4g.Info("Request completed",
        l4g.String("method", "GET"),
        l4g.String("path", "/api/users"),
        l4g.Int("status", 200),
        l4g.Duration("latency", time.Millisecond*42),
    )

    // JSON-style logging
    l4g.Infoj(map[string]any{
        "user":   "alice",
        "action": "login",
        "ip":     "192.168.1.1",
    })
}
```

## Usage

### Creating a Custom Logger

```go
// Create a logger with custom settings
logger := l4g.New(os.Stdout,
    l4g.WithLevel(l4g.LevelDebug),
)

logger.Debug("Debug message with details",
    l4g.String("component", "database"),
    l4g.Int("retries", 3),
)
```

### Named Channels

Create independent loggers for different components:

```go
// Each channel is cached and returns the same instance
dbLogger := l4g.Channel("database")
apiLogger := l4g.Channel("api")

dbLogger.Info("Connection established")
apiLogger.Info("Server listening", l4g.Int("port", 8080))
```

### Log Levels

```go
l4g.Trace("Detailed trace information")   // Finest level
l4g.Debug("Debug information")            // Development details
l4g.Info("Informational messages")        // Default level
l4g.Warn("Warning messages")              // Potentially harmful
l4g.Error("Error messages")               // Error conditions
l4g.Panic("Panic and recover")            // Logs then panics
l4g.Fatal("Fatal errors")                 // Logs then exits
```

### Conditional Logging

```go
// Set minimum log level
l4g.SetLevel(l4g.LevelWarn)

// These won't allocate or process arguments
l4g.Debug("This won't be logged")
l4g.Info("Neither will this")

// Only warnings and above are logged
l4g.Warn("This will be logged")
l4g.Error("So will this")
```

### Custom Handlers

```go
handler := l4g.NewSimpleHandler(l4g.HandlerOptions{
    Level:      l4g.NewLevelVar(l4g.LevelInfo),
    Output:     os.Stdout,
    TimeFormat: time.RFC3339,
    NoColor:    false,
    ReplaceAttr: func(groups []string, attr l4g.Attr) l4g.Attr {
        // Customize attribute formatting
        if attr.Key == "password" {
            return l4g.String("password", "***REDACTED***")
        }
        return attr
    },
})

logger := l4g.New(os.Stdout, l4g.WithHandler(handler))
```

### Formatted Logging

```go
// Printf-style formatting
l4g.Debugf("Processing %d items in %s", count, category)
l4g.Infof("Server started on port %d", port)
l4g.Errorf("Failed to connect: %v", err)

// JSON-style structured logging
l4g.Debugj(map[string]any{
    "operation": "query",
    "duration":  duration,
    "rows":      count,
})
```

### Attribute Types

```go
l4g.Info("User action",
    l4g.String("name", "alice"),
    l4g.Int("age", 30),
    l4g.Float("score", 98.5),
    l4g.Bool("active", true),
    l4g.Duration("elapsed", 100*time.Millisecond),
    l4g.Time("timestamp", time.Now()),
    l4g.Any("metadata", customStruct),
    l4g.Group("address",
        l4g.String("city", "New York"),
        l4g.String("country", "USA"),
    ),
)
```

### Colorized Output

```go
// Disable colors (e.g., for file output)
handler := l4g.NewSimpleHandler(l4g.HandlerOptions{
    Output:  file,
    NoColor: true,
})

// Custom colors for specific attributes
l4g.Info("Status update",
    l4g.ColorAttr(2, l4g.String("status", "success")), // Green
    l4g.ColorAttr(1, l4g.String("env", "production")), // Red
)
```

## Performance

l4g is optimized for high-performance logging:

- **Zero Allocations**: Disabled log levels result in zero memory allocations
- **Buffer Pooling**: Reuses buffers via `sync.Pool` to reduce GC pressure
- **Concurrent Safe**: Uses `sync.Map` for channel management, atomic operations for level checks
- **Pre-allocation**: Smart capacity estimation for slices to minimize reallocation

### Benchmark Results

```
BenchmarkPackageInfo-8        2000000    500 ns/op    0 B/op    0 allocs/op  (disabled)
BenchmarkPackageInfof-8       1000000   1200 ns/op  256 B/op    3 allocs/op  (enabled)
BenchmarkChannel-8           10000000    120 ns/op    0 B/op    0 allocs/op  (cached)
```

## API Reference

### Package-Level Functions

All standard logging methods are available at the package level:
- `Trace(msg, ...attrs)` / `Tracef(fmt, ...args)` / `Tracej(map)`
- `Debug(msg, ...attrs)` / `Debugf(fmt, ...args)` / `Debugj(map)`
- `Info(msg, ...attrs)` / `Infof(fmt, ...args)` / `Infoj(map)`
- `Warn(msg, ...attrs)` / `Warnf(fmt, ...args)` / `Warnj(map)`
- `Error(msg, ...attrs)` / `Errorf(fmt, ...args)` / `Errorj(map)`
- `Panic(msg, ...attrs)` / `Panicf(fmt, ...args)` / `Panicj(map)`
- `Fatal(msg, ...attrs)` / `Fatalf(fmt, ...args)` / `Fatalj(map)`

### Logger Configuration

- `New(w io.Writer, opts ...Option) *Logger`: Create a new logger
- `Default() *Logger`: Get the default logger
- `SetDefault(l *Logger)`: Set the default logger
- `Channel(name string) *Logger`: Get or create a named logger
- `SetLevel(level Level)`: Set minimum log level
- `GetLevel() Level`: Get current log level
- `SetOutput(w io.Writer)`: Change output destination

### Custom Types

- `WithLevel(level Level)`: Set initial log level
- `WithHandler(h Handler)`: Use custom handler
- `WithNewHandlerFunc(f func(HandlerOptions) Handler)`: Custom handler factory

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test -v

# Run with race detector
go test -race

# Run benchmarks
go test -bench=. -benchmem
```

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass: `go test ./...`
2. Code is formatted: `go fmt ./...`
3. Static analysis passes: `go vet ./...`
4. Add tests for new features

---

Built with ❤️ for high-performance Go applications.