package loggerus

import (
	"context"

	"github.com/nuclio/logger"
)

// a logger that multiplexes logs towards multiple loggers
type MuxLogger struct {
	loggers []logger.Logger
}

func NewMuxLogger(loggers ...logger.Logger) (*MuxLogger, error) {
	return &MuxLogger{loggers: loggers}, nil
}

func (ml *MuxLogger) SetLoggers(loggers ...logger.Logger) {
	ml.loggers = loggers
}

func (ml *MuxLogger) GetLoggers() []logger.Logger {
	return ml.loggers
}

func (ml *MuxLogger) Error(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.Error(format, vars...)
	}
}

func (ml *MuxLogger) ErrorCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.ErrorCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) Warn(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.Warn(format, vars...)
	}
}

func (ml *MuxLogger) WarnCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.WarnCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) Info(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.Info(format, vars...)
	}
}

func (ml *MuxLogger) InfoCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.InfoCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) Debug(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.Debug(format, vars...)
	}
}

func (ml *MuxLogger) DebugCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.DebugCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) ErrorWith(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.ErrorWith(format, vars...)
	}
}

func (ml *MuxLogger) ErrorWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.ErrorWithCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) WarnWith(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.WarnWith(format, vars...)
	}
}

func (ml *MuxLogger) WarnWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.WarnWithCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) InfoWith(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.InfoWith(format, vars...)
	}
}

func (ml *MuxLogger) InfoWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.InfoWithCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) DebugWith(format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.DebugWith(format, vars...)
	}
}

func (ml *MuxLogger) DebugWithCtx(ctx context.Context, format interface{}, vars ...interface{}) {
	for _, loggerInstance := range ml.loggers {
		loggerInstance.DebugWithCtx(ctx, format, vars...)
	}
}

func (ml *MuxLogger) Flush() {
}

func (ml *MuxLogger) GetChild(name string) logger.Logger {
	childLoggers := []logger.Logger{}
	for _, loggerInstance := range ml.loggers {
		childLogger := loggerInstance.GetChild(name)
		childLoggers = append(childLoggers, childLogger)
	}
	childMuxLogger, _ := NewMuxLogger(childLoggers...)

	return childMuxLogger
}
