package main

import (
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strconv"
	"valsea_coding_challenge/cmd"
	"valsea_coding_challenge/util"
)

func main() {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Failed to parse port: %s", portStr)
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	logLevel, err := zapcore.ParseLevel(logLevelStr)

	cmd.StartServer(&util.Config{
		Port:     port,
		LogLevel: logLevel,
	}, nil)
}
