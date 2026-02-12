package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DailyNotesRepositoryRealization struct {
	pool *pgxpool.Pool
}

func NewDailyNotesRepositoryRealization(pool *pgxpool.Pool) *DailyNotesRepositoryRealization {
	return &DailyNotesRepositoryRealization{
		pool: pool,
	}
}

func (d *DailyNotesRepositoryRealization) CreateNote(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
	sql := "INSERT INTO DailyEntries (id, user_id, date, mood, sleep_hours, load) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := d.pool.Exec(ctx, sql, id, userId, date, mood, sleepHours, load)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}
	return nil
}

func (d *DailyNotesRepositoryRealization) ChangeMood(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
	sql := "UPDATE DailyEntries SET mood = $1 WHERE user_id = $2 AND date = $3"
	tag, err := d.pool.Exec(ctx, sql, mood, userId, date)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrDailyEntryNotFound
	}
	return nil
}

func (d *DailyNotesRepositoryRealization) ChangeSleepHours(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
	sql := "UPDATE DailyEntries SET sleep_hours = $1 WHERE user_id = $2 AND date = $3"
	tag, err := d.pool.Exec(ctx, sql, sleepHours, userId, date)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrDailyEntryNotFound
	}
	return nil
}

func (d *DailyNotesRepositoryRealization) ChangeLoad(ctx context.Context, userId uuid.UUID, date time.Time, load int16) error {
	sql := "UPDATE DailyEntries SET load = $1 WHERE user_id = $2 AND date = $3"
	tag, err := d.pool.Exec(ctx, sql, load, userId, date)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrDailyEntryNotFound
	}
	return nil
}
