package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Job struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type ProcessedJob struct {
	ID          string    `json:"id"`
	Payload     string    `json:"payload"`
	ProcessedBy string    `json:"processed_by"`
	ProcessedAt time.Time `json:"processed_at"`
}

const schema = `
create table if not exists processed_jobs (
  id text primary key,
  payload text not null,
  processed_by text not null,
  processed_at timestamptz not null default now()
);
`

func OpenPostgres(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	return pool, nil
}

func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}
	return nil
}

func OpenRedis(redisAddr, password string, database int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		Password:        password,
		DB:              database,
		Protocol:        2,
		DisableIdentity: true,
	})
}

func EnqueueJob(ctx context.Context, redisClient *redis.Client, queueName string, payload string) (Job, error) {
	job := Job{
		ID:        uuid.NewString(),
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return Job{}, fmt.Errorf("marshal job: %w", err)
	}

	if err := redisClient.RPush(ctx, queueName, body).Err(); err != nil {
		return Job{}, fmt.Errorf("enqueue job: %w", err)
	}

	return job, nil
}

func FetchRecentJobs(ctx context.Context, db *pgxpool.Pool, limit int) ([]ProcessedJob, error) {
	if err := EnsureSchema(ctx, db); err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, `
select id, payload, processed_by, processed_at
from processed_jobs
order by processed_at desc
limit $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("query jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]ProcessedJob, 0, limit)
	for rows.Next() {
		var job ProcessedJob
		if err := rows.Scan(&job.ID, &job.Payload, &job.ProcessedBy, &job.ProcessedAt); err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate jobs: %w", err)
	}

	return jobs, nil
}

func StoreProcessedJob(ctx context.Context, db *pgxpool.Pool, job Job, processedBy string) error {
	if err := EnsureSchema(ctx, db); err != nil {
		return err
	}

	_, err := db.Exec(ctx, `
insert into processed_jobs (id, payload, processed_by, processed_at)
values ($1, $2, $3, now())
on conflict (id) do update
set payload = excluded.payload,
    processed_by = excluded.processed_by,
    processed_at = excluded.processed_at
`, job.ID, job.Payload, processedBy)
	if err != nil {
		return fmt.Errorf("store processed job: %w", err)
	}
	return nil
}
