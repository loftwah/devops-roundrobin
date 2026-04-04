package platform

import (
	"testing"
	"time"
)

func TestLoadConfigUsesFallbacks(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_HTTP_ADDR", "")
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("POSTGRES_PORT", "")
	t.Setenv("POSTGRES_USER", "")
	t.Setenv("POSTGRES_PASSWORD", "")
	t.Setenv("POSTGRES_DB", "")
	t.Setenv("REDIS_HOST", "")
	t.Setenv("REDIS_PORT", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("JOB_QUEUE_NAME", "")
	t.Setenv("READY_TIMEOUT", "")

	cfg := LoadConfig("app", ":8080")

	if cfg.HTTPAddress != ":8080" {
		t.Fatalf("expected default app address, got %q", cfg.HTTPAddress)
	}

	expectedDatabaseURL := "postgres://roundrobin:roundrobin@localhost:5432/roundrobin?sslmode=disable"
	if cfg.DatabaseURL != expectedDatabaseURL {
		t.Fatalf("expected database URL %q, got %q", expectedDatabaseURL, cfg.DatabaseURL)
	}

	if cfg.RedisAddr != "localhost:6379" {
		t.Fatalf("expected default redis address, got %q", cfg.RedisAddr)
	}

	if cfg.JobQueue != "round-robin-jobs" {
		t.Fatalf("expected default queue name, got %q", cfg.JobQueue)
	}

	if cfg.ReadyTimeout != 2*time.Second {
		t.Fatalf("expected default ready timeout, got %s", cfg.ReadyTimeout)
	}
}

func TestLoadConfigReadsWorkerAddress(t *testing.T) {
	t.Setenv("WORKER_HTTP_ADDR", ":9999")

	cfg := LoadConfig("worker", ":8081")

	if cfg.HTTPAddress != ":9999" {
		t.Fatalf("expected worker address override, got %q", cfg.HTTPAddress)
	}
}
