package usecase

import (
	"chopper/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type AlertService struct {
	alertRepository AlertRepository
}

func NewAlertServcie(alertRepository AlertRepository) *AlertService {
	return &AlertService{
		alertRepository: alertRepository,
	}
}

func (a *AlertService) GetLastSevenDays(ctx context.Context, userId uuid.UUID) (string, error) {
	notes, err := a.alertRepository.GetLastSevenDays(ctx, userId)
	if err != nil {
		return "", err
	}
	alert, ok := isAlert(notes)
	if ok {
		return alert, nil
	}
	return "Все хорошо", nil
}

func isAlert(days []domain.Day) (string, bool) {
	if len(days) < 3 {
		return "", false
	}
	for i := 0; i+2 < len(days); i++ {
		zero := days[i]
		one := days[i+1]
		two := days[i+2]
		first := isNextDay(two.Date, one.Date)
		second := isNextDay(one.Date, zero.Date)
		if !first || !second {
			continue
		}
		moodOne := zero.Mood
		moodTwo := one.Mood
		moodThree := two.Mood
		badMood := isMoodBad(moodOne, moodTwo, moodThree)
		sleepHoursOne := zero.SleepHours
		sleepHoursTwo := one.SleepHours
		sleepHoursThird := two.SleepHours
		lowSleepHours := isSleepHoursLow(sleepHoursOne, sleepHoursTwo, sleepHoursThird)
		loadOne := zero.Load
		loadTwo := one.Load
		loadThree := two.Load
		isLoadHigh := isLoadHigh(loadOne, loadTwo, loadThree)
		if badMood && lowSleepHours {
			return "За последние дни низкий уровень настроения и мало сна", true
		}
		if badMood && isLoadHigh {
			return "За последние дни низкий уровень настроения и большая загрузка", true
		}
		if lowSleepHours && isLoadHigh {
			return "За последние дни мало сна и большая загрузка", true
		}
	}
	return "", false
}

func isNextDay(one, two time.Time) bool {
	y1, m1, d1 := one.Year(), one.Month(), one.Day()
	y2, m2, d2 := two.Year(), two.Month(), two.Day()
	if time.Date(y1, m1, d1+1, 0, 0, 0, 0, time.UTC).Equal(time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)) {
		return true
	}
	return false
}

func isMoodBad(moodOne, moodTwo, moodThree int16) bool {
	if moodOne <= 5 && moodTwo <= 5 {
		return true
	}
	if moodOne <= 5 && moodThree <= 5 {
		return true
	}
	if moodTwo <= 5 && moodThree <= 5 {
		return true
	}
	return false
}

func isSleepHoursLow(sleepHoursOne, sleepHoursTwo, sleepHoursThree float64) bool {
	hours := []float64{sleepHoursOne, sleepHoursTwo, sleepHoursThree}
	bad := 0
	for _, hour := range hours {
		if hour <= 7.0 {
			bad++
		}
	}
	if bad >= 2 {
		return true
	}
	return false
}

func isLoadHigh(loadOne, loadTwo, loadThree int16) bool {
	loads := []int16{loadOne, loadTwo, loadThree}
	bad := 0
	for _, load := range loads {
		if load >= 5 {
			bad++
		}
	}
	if bad >= 2 {
		return true
	}
	return false
}
