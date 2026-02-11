package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvOrDefault returns the environment variable value, or the fallback if unset or empty.
func EnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ParseBrokers splits a comma-separated broker list, trimming whitespace.
func ParseBrokers(value string) []string {
	parts := strings.Split(value, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}

// ParseBatchSize reads BATCH_SIZE from the environment.
// Default: 50. Range: 1-1000.
func ParseBatchSize() (int, error) {
	s := os.Getenv("BATCH_SIZE")
	if s == "" {
		return 50, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 1000 {
		return 0, errors.New("invalid BATCH_SIZE: must be 1-1000")
	}
	return n, nil
}

// ParseBatchFlushInterval reads BATCH_FLUSH_INTERVAL from the environment.
// Default: 500ms. Must be a positive duration.
func ParseBatchFlushInterval() (time.Duration, error) {
	s := EnvOrDefault("BATCH_FLUSH_INTERVAL", "500ms")
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, errors.New("invalid BATCH_FLUSH_INTERVAL: must be a positive duration")
	}
	return d, nil
}

// ParseShutdownTimeout reads SHUTDOWN_TIMEOUT from the environment.
// Default: 10s. Must be a positive duration.
func ParseShutdownTimeout() (time.Duration, error) {
	s := EnvOrDefault("SHUTDOWN_TIMEOUT", "10s")
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, errors.New("invalid SHUTDOWN_TIMEOUT")
	}
	return d, nil
}
