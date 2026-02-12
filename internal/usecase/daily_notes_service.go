package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type DailyNotesService struct {
	dailyNotesRepository DailyNotesRepository
}

func NewDailyNotesService(dailyNotesRepository DailyNotesRepository) *DailyNotesService {
	return &DailyNotesService{
		dailyNotesRepository: dailyNotesRepository,
	}
}

func (d *DailyNotesService) CreateNote(ctx context.Context, userId uuid.UUID, dailyNoteFromFront domain.DailyNoteFromFront) error {
	id := uuid.New()
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	mood := dailyNoteFromFront.Mood
	sleepHours := dailyNoteFromFront.SleepHours
	if sleepHours > 9.9 || sleepHours < 0.0 {
		return ErrWrongSleepHourValue
	}
	load := dailyNoteFromFront.Load
	if err := d.dailyNotesRepository.CreateNote(ctx, id, userId, date, mood, sleepHours, load); err != nil && errors.Is(err, repository.ErrUniqueViolation) {
		return ErrNoteAlreadyExists
	} else if err != nil {
		return err
	}
	return nil
}

func (d *DailyNotesService) ChangeMood(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) (string, error) {
	if err := d.dailyNotesRepository.ChangeMood(ctx, userId, date, mood); err != nil {
		if errors.Is(err, repository.ErrDailyEntryNotFound) {
			return "", ErrNoteNotExists
		}
		return "", err
	}
	return "mood успешно изменен", nil
}

func (d *DailyNotesService) ChangeSleepHours(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) (string, error) {
	if err := d.dailyNotesRepository.ChangeSleepHours(ctx, userId, date, sleepHours); err != nil {
		if errors.Is(err, repository.ErrDailyEntryNotFound) {
			return "", ErrNoteNotExists
		}
		return "", err
	}
	return "sleep hours успешно изменен", nil
}

func (d *DailyNotesService) ChangeLoad(ctx context.Context, userId uuid.UUID, date time.Time, load int16) (string, error) {
	if err := d.dailyNotesRepository.ChangeLoad(ctx, userId, date, load); err != nil {
		if errors.Is(err, repository.ErrDailyEntryNotFound) {
			return "", ErrNoteNotExists
		}
		return "", err
	}
	return "load успешно изменен", nil
}
