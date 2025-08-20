package logs

import (
	"log/slog"
	"os"
)

var slogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}))

type CronLogger struct {
	logger  interface{ Printf(string, ...interface{}) }
	logInfo bool
}

func (l CronLogger) Info(msg string, keysAndValues ...interface{}) {
	slogger.Info(msg, keysAndValues...)
}

func (l CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error")
	keysAndValues = append(keysAndValues, err.Error())
	slogger.Error(msg, keysAndValues...)
}

func Info(msg string, data ...interface{}) {
	slogger.Info(msg, data...)
}

func Fatal(msg string, data ...interface{}) {
	slogger.Error(msg, data...)
	os.Exit(1)
}
