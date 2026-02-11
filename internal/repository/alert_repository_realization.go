package repository

import (
	"chopper/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepositoryRealization struct {
	pool *pgxpool.Pool
}

func NewAlertRepositoryRealization(pool *pgxpool.Pool) *AlertRepositoryRealization {
	return &AlertRepositoryRealization{
		pool: pool,
	}
}

func (a *AlertRepositoryRealization) GetLastSevenDays(ctx context.Context, userId uuid.UUID) ([]domain.Day, error) {
	sql := "SELECT date, mood, sleep_hours, load FROM DailyEntries WHERE user_id = $1 ORDER BY date DESC LIMIT 7"
	rows, err := a.pool.Query(ctx, sql, userId)
	if err != nil {
		return []domain.Day{}, err
	}
	defer rows.Close()
	days := []domain.Day{}
	for rows.Next() {
		var day domain.Day
		if err := rows.Scan(&day.Date, &day.Mood, &day.SleepHours, &day.Load); err != nil {
			return []domain.Day{}, err
		}
		days = append(days, day)
	}
	if err := rows.Err(); err != nil {
		return []domain.Day{}, err
	}
	return days, nil
}
