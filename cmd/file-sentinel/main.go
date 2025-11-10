package main

import (
	"file-sentinel/cmd/file-sentinel/repository"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	targetDirPath = "/mnt/harman"
	minuteCycle   = 1
)

func init() {
	envDirPath := os.Getenv("TARGET_DIR_PATH")
	envMinuteCycle := os.Getenv("MINUTE_CYCLE")
	envLogLevel := os.Getenv("LOG_LEVEL")

	if envDirPath != "" {
		targetDirPath = envDirPath
	}
	if envMinuteCycle != "" {
		cycle, err := strconv.Atoi(strings.TrimSpace(envMinuteCycle))
		minuteCycle = cycle
		if err != nil {
			slog.Info("Invalid MINUTE_CYCLE value: %v. Using default value 1.")
			minuteCycle = 1
		}
	}

	if envLogLevel != "" {
		envLogLevel = strings.ToUpper(strings.TrimSpace(envLogLevel))
		switch envLogLevel {
		case "DEBUG":
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			slog.SetDefault(logger)
		case "INFO":
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			slog.SetDefault(logger)
		case "ERROR":
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
			slog.SetDefault(logger)
		default:
			slog.Info("Invalid LOG_LEVEL value: %v. Using default INFO level.")
		}
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)
	}
}

func main() {
	slog.Info("Starting File Sentinel...")
	slog.Info("env TARGET_DIR_PATH value", "TARGET_DIR_PATH", os.Getenv("TARGET_DIR_PATH"))
	slog.Info("env MINUTE_CYCLE value", "MINUTE_CYCLE", os.Getenv("MINUTE_CYCLE"))
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		slog.Error("Error creating repository", "error", err)
		return
	}

	fileSentinel := NewFileSentinel(repo, targetDirPath, minuteCycle)
	err = fileSentinel.Start()
	if err != nil {
		slog.Error("Error starting file sentinel", "error", err)
		return
	}
}
