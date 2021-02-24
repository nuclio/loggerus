# Loggerus

A Nuclio.Logger wrapper for logrus

Usage example:

```golang
package myapp

import "github.com/nuclio/loggerus"

// Enrich with who
enrichWhoField := true

// Color output
color := true

// Create a text logger instance
logger, _ := loggerus.NewTextLoggerus("app-logger", logrus.DebugLevel, os.Stdout, enrichWhoField, color)

// Log
logger.Debug("Hello from Loggerus")

// Structured Log
logger.DebugWith("Hello from Loggerus", "someValue", 123)

// Structured Log + context object
logger.DebugWithCtx(ctx, "Hello from Loggerus", "someValue", 123)
```