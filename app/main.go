package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/deanlofts/devops-roundrobin/internal/build"
	"github.com/deanlofts/devops-roundrobin/internal/platform"
)

func main() {
	cfg := platform.LoadConfig("app", ":8080")
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
	jobEnqueueCount := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "round_robin",
		Subsystem: "app",
		Name:      "jobs_enqueued_total",
		Help:      "Total jobs accepted by the app service.",
	})

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", platform.WithRequestLogging(logger, requestMetrics, "health", platform.HealthHandler(cfg)))
	mux.Handle("/ready", platform.WithRequestLogging(logger, requestMetrics, "ready", platform.ReadyHandler(cfg, db, redisClient)))

	mux.Handle("/", platform.WithRequestLogging(logger, requestMetrics, "root", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platform.JSON(w, http.StatusOK, map[string]any{
			"service":     cfg.ServiceName,
			"environment": cfg.Environment,
			"version":     build.Version,
			"commit":      build.Commit,
			"built_at":    build.Date,
			"message":     "DevOps Round Robin platform app is running.",
			"queue":       cfg.JobQueue,
			"time":        time.Now().UTC(),
		})
	})))

	mux.Handle("/jobs", platform.WithRequestLogging(logger, requestMetrics, "jobs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobs, err := platform.FetchRecentJobs(r.Context(), db, 20)
			if err != nil {
				logger.Error("failed to fetch jobs", slog.Any("error", err))
				platform.JSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
				return
			}
			platform.JSON(w, http.StatusOK, map[string]any{"jobs": jobs})
		case http.MethodPost:
			payload := "hello from the app"

			if r.Body != nil {
				defer r.Body.Close()
				body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
				if err != nil {
					platform.JSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read body"})
					return
				}

				if len(body) > 0 {
					var request struct {
						Payload string `json:"payload"`
					}
					if err := json.Unmarshal(body, &request); err != nil {
						platform.JSON(w, http.StatusBadRequest, map[string]string{"error": "body must be valid JSON"})
						return
					}
					if request.Payload != "" {
						payload = request.Payload
					}
				}
			}

			job, err := platform.EnqueueJob(r.Context(), redisClient, cfg.JobQueue, payload)
			if err != nil {
				logger.Error("failed to enqueue job", slog.Any("error", err))
				platform.JSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
				return
			}

			jobEnqueueCount.Inc()
			platform.JSON(w, http.StatusAccepted, map[string]any{
				"message": "job queued",
				"job":     job,
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("starting app server",
		slog.String("address", cfg.HTTPAddress),
		slog.String("service", cfg.ServiceName),
		slog.String("environment", cfg.Environment),
		slog.String("version", build.Version),
		slog.String("commit", build.Commit),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("app server failed", slog.Any("error", err))
			os.Exit(1)
		}
	case <-signalCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGracePeriod)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("app server stopped cleanly")
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
