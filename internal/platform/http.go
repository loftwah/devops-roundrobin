package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

type StatusResponse struct {
	Service      string            `json:"service"`
	Environment  string            `json:"environment"`
	Status       string            `json:"status"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Version      string            `json:"version,omitempty"`
	Commit       string            `json:"commit,omitempty"`
	Time         time.Time         `json:"time"`
}

type RequestMetrics struct {
	RequestCount   *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
}

func NewRequestMetrics(serviceName string) *RequestMetrics {
	return &RequestMetrics{
		RequestCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "round_robin",
			Subsystem: serviceName,
			Name:      "http_requests_total",
			Help:      "Total HTTP requests handled by the service.",
		}, []string{"method", "route", "status"}),
		RequestLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "round_robin",
			Subsystem: serviceName,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latency by route.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "route", "status"}),
	}
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func HealthHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if cfg.FailHealth {
			JSON(w, http.StatusInternalServerError, StatusResponse{
				Service:     cfg.ServiceName,
				Environment: cfg.Environment,
				Status:      "forced-health-failure",
				Time:        time.Now().UTC(),
			})
			return
		}

		JSON(w, http.StatusOK, StatusResponse{
			Service:     cfg.ServiceName,
			Environment: cfg.Environment,
			Status:      "ok",
			Time:        time.Now().UTC(),
		})
	}
}

func ReadyHandler(cfg Config, db *pgxpool.Pool, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if cfg.FailReady {
			JSON(w, http.StatusServiceUnavailable, StatusResponse{
				Service:     cfg.ServiceName,
				Environment: cfg.Environment,
				Status:      "forced-ready-failure",
				Time:        time.Now().UTC(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ReadyTimeout)
		defer cancel()

		dependencies := map[string]string{
			"postgres": "unknown",
			"redis":    "unknown",
		}

		statusCode := http.StatusOK
		status := "ready"

		if err := db.Ping(ctx); err != nil {
			dependencies["postgres"] = fmt.Sprintf("down: %v", err)
			status = "not-ready"
			statusCode = http.StatusServiceUnavailable
		} else {
			dependencies["postgres"] = "up"
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			dependencies["redis"] = fmt.Sprintf("down: %v", err)
			status = "not-ready"
			statusCode = http.StatusServiceUnavailable
		} else {
			dependencies["redis"] = "up"
		}

		JSON(w, statusCode, StatusResponse{
			Service:      cfg.ServiceName,
			Environment:  cfg.Environment,
			Status:       status,
			Dependencies: dependencies,
			Time:         time.Now().UTC(),
		})
	}
}

func WithRequestLogging(logger *slog.Logger, metrics *RequestMetrics, routeName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		status := fmt.Sprintf("%d", recorder.statusCode)
		metrics.RequestCount.WithLabelValues(r.Method, routeName, status).Inc()
		metrics.RequestLatency.WithLabelValues(r.Method, routeName, status).Observe(time.Since(start).Seconds())

		logger.Info("request completed",
			slog.String("route", routeName),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", recorder.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
