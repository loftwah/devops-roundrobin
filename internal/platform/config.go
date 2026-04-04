package platform

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName         string
	Environment         string
	LogLevel            string
	HTTPAddress         string
	DatabaseURL         string
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	JobQueue            string
	FailHealth          bool
	FailReady           bool
	ReadyTimeout        time.Duration
	ShutdownGracePeriod time.Duration
}

func LoadConfig(serviceName string, defaultHTTPAddress string) Config {
	postgresHost := envOrDefault("POSTGRES_HOST", "localhost")
	postgresPort := envOrDefault("POSTGRES_PORT", "5432")
	postgresUser := envOrDefault("POSTGRES_USER", "roundrobin")
	postgresPassword := envOrDefault("POSTGRES_PASSWORD", "roundrobin")
	postgresDB := envOrDefault("POSTGRES_DB", "roundrobin")

	return Config{
		ServiceName:         serviceName,
		Environment:         envOrDefault("APP_ENV", "local"),
		LogLevel:            envOrDefault("LOG_LEVEL", "INFO"),
		HTTPAddress:         envOrDefault(defaultHTTPEnvKey(serviceName), defaultHTTPAddress),
		DatabaseURL:         envOrDefault("DATABASE_URL", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)),
		RedisAddr:           envOrDefault("REDIS_ADDR", fmt.Sprintf("%s:%s", envOrDefault("REDIS_HOST", "localhost"), envOrDefault("REDIS_PORT", "6379"))),
		RedisPassword:       os.Getenv("REDIS_PASSWORD"),
		RedisDB:             envInt("REDIS_DB", 0),
		JobQueue:            envOrDefault("JOB_QUEUE_NAME", "round-robin-jobs"),
		FailHealth:          envBool("FAIL_HEALTH", false),
		FailReady:           envBool("FAIL_READY", false),
		ReadyTimeout:        envDuration("READY_TIMEOUT", 2*time.Second),
		ShutdownGracePeriod: envDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func defaultHTTPEnvKey(serviceName string) string {
	switch serviceName {
	case "worker":
		return "WORKER_HTTP_ADDR"
	default:
		return "APP_HTTP_ADDR"
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
