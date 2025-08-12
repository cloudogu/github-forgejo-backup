package main

import (
	"log/slog"
	"os"
)

var logs = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}))

type cronLogger struct {
	logger  interface{ Printf(string, ...interface{}) }
	logInfo bool
}

func (l cronLogger) Info(msg string, keysAndValues ...interface{}) {
	logs.Info(msg, keysAndValues...)
}

func (l cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error")
	keysAndValues = append(keysAndValues, err.Error())
	logs.Error(msg, keysAndValues...)
}
