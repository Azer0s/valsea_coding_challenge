package util

import "go.uber.org/zap/zapcore"

type Config struct {
	Port     int
	LogLevel zapcore.Level
}
