package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/deanlofts/devops-roundrobin/internal/build"
	"github.com/deanlofts/devops-roundrobin/internal/platform"
)

func main() {
	cfg := platform.LoadConfig("worker", ":8081")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(cfg.LogLevel),
	}))

	ctx := context.Background()

	db, err := platform.OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	redisClient := platform.OpenRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer redisClient.Close()

	requestMetrics := platform.NewRequestMetrics(cfg.ServiceName)
	jobProcessedCount := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "round_robin",
		Subsystem: "worker",
		Name:      "jobs_processed_total",
		Help:      "Total jobs processed by the worker.",
	})
	jobFailedCount := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "round_robin",
		Subsystem: "worker",
		Name:      "jobs_failed_total",
		Help:      "Total jobs that failed processing.",
	})

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", platform.WithRequestLogging(logger, requestMetrics, "health", platform.HealthHandler(cfg)))
	mux.Handle("/ready", platform.WithRequestLogging(logger, requestMetrics, "ready", platform.ReadyHandler(cfg, db, redisClient)))
	mux.Handle("/", platform.WithRequestLogging(logger, requestMetrics, "root", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		platform.JSON(w, http.StatusOK, map[string]any{
			"service":     cfg.ServiceName,
			"environment": cfg.Environment,
			"version":     build.Version,
			"commit":      build.Commit,
			"built_at":    build.Date,
			"message":     "Round Robin worker is ready to consume jobs.",
			"time":        time.Now().UTC(),
		})
	})))

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("starting worker",
		slog.String("address", cfg.HTTPAddress),
		slog.String("service", cfg.ServiceName),
		slog.String("environment", cfg.Environment),
	)

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	go runWorker(workerCtx, logger, cfg, db, redisClient, jobProcessedCount, jobFailedCount)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("worker http server failed", slog.Any("error", err))
			os.Exit(1)
		}
	case <-signalCtx.Done():
		logger.Info("shutdown signal received")
	}

	workerCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGracePeriod)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("worker stopped cleanly")
}

func runWorker(
	ctx context.Context,
	logger *slog.Logger,
	cfg platform.Config,
	db *pgxpool.Pool,
	redisClient *redis.Client,
	jobProcessedCount prometheus.Counter,
	jobFailedCount prometheus.Counter,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		results, err := redisClient.BLPop(ctx, 5*time.Second, cfg.JobQueue).Result()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			if errors.Is(err, redis.Nil) {
				continue
			}
			logger.Error("failed to fetch job from redis", slog.Any("error", err))
			jobFailedCount.Inc()
			time.Sleep(2 * time.Second)
			continue
		}

		if len(results) < 2 {
			continue
		}

		var job platform.Job
		if err := json.Unmarshal([]byte(results[1]), &job); err != nil {
			logger.Error("failed to decode job payload", slog.Any("error", err))
			jobFailedCount.Inc()
			continue
		}

		if err := platform.StoreProcessedJob(ctx, db, job, cfg.ServiceName); err != nil {
			logger.Error("failed to store processed job", slog.Any("error", err), slog.String("job_id", job.ID))
			jobFailedCount.Inc()
			continue
		}

		jobProcessedCount.Inc()
		logger.Info("processed job", slog.String("job_id", job.ID), slog.String("payload", job.Payload))
	}
}

func parseLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
