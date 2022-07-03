package logger

import (
	"io"

	"com.thanos/pkg/config"
	"github.com/sirupsen/logrus"
)

// Logger entry
type Logger struct {
	*logrus.Entry
}

// NewLogger creates a new logger
func NewLogger(cfg config.Logger, options ...Option) *Logger {
	l := logrus.New()
	l.SetLevel(logrus.Level(cfg.LogLevel))
	l.SetReportCaller(false)
	l.SetFormatter(&logrus.JSONFormatter{})

	entry := l.WithFields(logrus.Fields{
		"application":           cfg.AppName,
		"package.build_version": cfg.AppVersion,
	})

	logger := Logger{
		Entry: entry,
	}

	return &logger
}

// DisableOutput discards log output (mainly for tests)
func DisableOutput() Option {
	return func(args *Logger) {
		args.Logger.SetOutput(io.Discard)
	}
}
