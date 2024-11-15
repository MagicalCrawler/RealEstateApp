package utils

import (
	"log/slog"
	"os"
	"sync"
)

var mainLoggerProducer func() *slog.Logger = sync.OnceValue(func() *slog.Logger {
	return NewLogger("main")
})

func NewLogger(category string) *slog.Logger {
	logLevel := slog.LevelInfo
	logLevelConfigValue := GetConfig("LOG_LEVEL")
	if err := logLevel.UnmarshalText([]byte(logLevelConfigValue)); err != nil {
		slog.Warn("LOG_LEVEL not set or has wrong value", slog.Attr{Key: "err", Value: slog.AnyValue(err)})
	}

	var logFile *os.File = os.Stdout
	logPath := GetConfig("LOG_PATH")
	os.MkdirAll(logPath, os.ModePerm)
	logFile, err := os.OpenFile(logPath+"/"+category+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("could not open log file", slog.Attr{Key: "err", Value: slog.AnyValue(err)})
	}
	return slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: logLevel}))
}

func MainLogger() *slog.Logger {
	return mainLoggerProducer()
}
