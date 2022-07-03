package logger

import "github.com/sirupsen/logrus"

type Option func(*Logger)

func EnableReportCaller() Option {
	return func(args *Logger) {
		args.Logger.SetReportCaller(true)
	}
}

func SetLevel(level uint8) Option {
	return func(args *Logger) {
		args.Logger.SetLevel(logrus.Level(level))
	}
}
