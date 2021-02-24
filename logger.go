package loggerus

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nuclio/logger"
	"github.com/sirupsen/logrus"
)

type Loggerus struct {
	logrus *logrus.Logger
	name   string
	output io.Writer
}

// Creates a logger pre-configured for commands
func NewJSONLoggerus(name string, level logrus.Level, output io.Writer) (*Loggerus, error) {

	// default timestamp formatting, and local timezone - defaults
	loggerJSONFormatter, err := newJSONFormatter("", "")
	if err != nil {
		return nil, err
	}

	return createLoggerus(name, level, output, loggerJSONFormatter)
}

func NewTextLoggerus(name string, level logrus.Level, output io.Writer, enrichWhoField bool, color bool) (*Loggerus, error) {
	loggerTextFormatter, err := newTextFormatter(0, enrichWhoField, color)
	if err != nil {
		return nil, err
	}

	return createLoggerus(name, level, output, loggerTextFormatter)
}

func NewLoggerusForTests(name string) (*Loggerus, error) {
	var loggerLevel logrus.Level

	if isVerboseTesting() {
		loggerLevel = logrus.DebugLevel
	} else {
		loggerLevel = logrus.InfoLevel
	}

	loggerRedactor := NewRedactor(os.Stdout)
	loggerRedactor.Disable()

	return NewTextLoggerus(name, loggerLevel, loggerRedactor, true, true)
}

func createLoggerus(name string, level logrus.Level, output io.Writer, formatter logrus.Formatter) (*Loggerus, error) {
	newLoggerus := Loggerus{
		logrus: logrus.New(),
		name:   name,
		output: output,
	}

	newLoggerus.logrus.SetOutput(output)
	newLoggerus.logrus.SetLevel(level)
	newLoggerus.logrus.SetFormatter(formatter)

	return &newLoggerus, nil
}

// Error emits an unstructured error log
func (l *Loggerus) Error(format interface{}, vars ...interface{}) {
	l.logrus.Errorf(format.(string), vars...)
}

// Warn emits an unstructured warning log
func (l *Loggerus) Warn(format interface{}, vars ...interface{}) {
	l.logrus.Warnf(format.(string), vars...)
}

// Info emits an unstructured informational log
func (l *Loggerus) Info(format interface{}, vars ...interface{}) {
	l.logrus.Infof(format.(string), vars...)
}

// Debug emits an unstructured debug log
func (l *Loggerus) Debug(format interface{}, vars ...interface{}) {
	l.logrus.Debugf(format.(string), vars...)
}

// ErrorCtx emits an unstructured error log with context
func (l *Loggerus) ErrorCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.Error(l.getFormatWithContext(ctx, format), vars...)
}

// WarnCtx emits an unstructured warning log with context
func (l *Loggerus) WarnCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.Warn(l.getFormatWithContext(ctx, format), vars...)
}

// InfoCtx emits an unstructured informational log with context
func (l *Loggerus) InfoCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.Info(l.getFormatWithContext(ctx, format), vars...)
}

// DebugCtx emits an unstructured debug log with context
func (l *Loggerus) DebugCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.Debug(l.getFormatWithContext(ctx, format), vars...)
}

// ErrorWith emits a structured error log
func (l *Loggerus) ErrorWith(format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFields(vars)).Error(format)
}

// WarnWith emits a structured warning log
func (l *Loggerus) WarnWith(format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFields(vars)).Warn(format)
}

// InfoWith emits a structured info log
func (l *Loggerus) InfoWith(format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFields(vars)).Info(format)
}

// DebugWith emits a structured debug log
func (l *Loggerus) DebugWith(format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFields(vars)).Debug(format)
}

// ErrorWithCtx emits a structured error log with context
func (l *Loggerus) ErrorWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFieldsWithCtx(ctx, vars)).Error(format)
}

// WarnWithCtx emits a structured warning log with context
func (l *Loggerus) WarnWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFieldsWithCtx(ctx, vars)).Warn(format)
}

// InfoWithCtx emits a structured info log with context
func (l *Loggerus) InfoWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFieldsWithCtx(ctx, vars)).Info(format)
}

// DebugWithCtx emits a structured debug log with context
func (l *Loggerus) DebugWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	l.logrus.WithFields(l.varsToFieldsWithCtx(ctx, vars)).Debug(format)
}

// Flush flushes buffered logs, if applicable
func (l *Loggerus) Flush() {
}

// GetChild returns a child logger, if underlying logger supports hierarchal logging
func (l *Loggerus) GetChild(name string) logger.Logger {
	childLogger, _ := createLoggerus(l.name+"."+name,
		l.logrus.Level,
		l.logrus.Out,
		l.logrus.Formatter)

	return childLogger
}

func (l *Loggerus) GetRedactor() *Redactor {
	logRedactor, ok := l.logrus.Out.(*Redactor)
	if ok {
		return logRedactor
	}
	return nil
}

func (l *Loggerus) GetOutput() io.Writer {
	return l.output
}

func (l *Loggerus) GetLogrus() *logrus.Logger {
	return l.logrus
}

func (l *Loggerus) getFormatWithContext(ctx context.Context, format interface{}) string {
	formatString := format.(string)

	// get request ID from context
	requestID := ctx.Value("RequestID")

	// if not set, don't add it to vars
	if requestID == nil || requestID == "" {
		return formatString
	}

	return formatString + fmt.Sprintf(" (requestID: %s)", requestID)
}

func (l *Loggerus) varsToFields(vars []interface{}) logrus.Fields {
	fields := logrus.Fields{}

	// enrich with who
	fields["who"] = l.name

	for varIndex := 0; varIndex < len(vars); varIndex += 2 {
		fields[vars[varIndex].(string)] = vars[varIndex+1]
	}

	return fields
}

func (l *Loggerus) varsToFieldsWithCtx(ctx context.Context, vars []interface{}) logrus.Fields {

	if ctx != nil {
		vars = append(vars, "ctx")
		vars = append(vars, ctx)

		// special treatment - 1st class fields
		for _, key := range []string{"RequestID", "SystemID"} {
			if value, ok := ctx.Value(key).(string); ok {
				vars = append(vars, key, value)
			}
		}
	}

	return l.varsToFields(vars)
}

// use this instead of testing.Verbose since we don't want to get testing flags in our code
func isVerboseTesting() bool {
	for _, arg := range os.Args {
		if arg == "-test.v=true" || arg == "-test.v" {
			return true
		}
	}
	return false
}
