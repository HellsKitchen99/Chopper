package usecase

import (
	"chopper/internal/domain"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockAlertRepository struct {
	GetLastSevenDaysFn func(ctx context.Context, userId uuid.UUID) ([]domain.Day, error)
	// переданные аргументы
	getLastSevenDaysFnIsCalled bool
	userId                     uuid.UUID
}

func (m *MockAlertRepository) GetLastSevenDays(ctx context.Context, userId uuid.UUID) ([]domain.Day, error) {
	m.getLastSevenDaysFnIsCalled = true
	m.userId = userId
	if m.GetLastSevenDaysFn != nil {
		return m.GetLastSevenDaysFn(ctx, userId)
	}
	return nil, nil
}

// Тест GetLastSevenDays - Успех (Есть возврат как алерта так и сообщения о том что все хорошо)
func TestGetLastSevenDaysSuccess(t *testing.T) {
	// preparing
	daysAlert := []domain.Day{
		domain.Day{
			Date:       time.Date(2025, 1, 3, 0, 0, 0, 0, time.Now().Location()),
			Mood:       4,
			SleepHours: 9.0,
			Load:       6,
		},
		domain.Day{
			Date:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.Now().Location()),
			Mood:       4,
			SleepHours: 9.0,
			Load:       6,
		},
		domain.Day{
			Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.Now().Location()),
			Mood:       4,
			SleepHours: 9.0,
			Load:       6,
		},
	}
	daysNotAlert := []domain.Day{
		domain.Day{
			Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.Now().Location()),
			Mood:       5,
			SleepHours: 5.5,
			Load:       5,
		},
		domain.Day{
			Date:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.Now().Location()),
			Mood:       5,
			SleepHours: 5.5,
			Load:       5,
		},
		domain.Day{
			Date:       time.Date(2025, 1, 4, 0, 0, 0, 0, time.Now().Location()),
			Mood:       5,
			SleepHours: 5.5,
			Load:       5,
		},
	}
	mockAlertRepositoryAlert := &MockAlertRepository{
		GetLastSevenDaysFn: func(ctx context.Context, userId uuid.UUID) ([]domain.Day, error) {
			return daysAlert, nil
		},
	}
	mockAlertRepositoryNotAlert := &MockAlertRepository{
		GetLastSevenDaysFn: func(ctx context.Context, userId uuid.UUID) ([]domain.Day, error) {
			return daysNotAlert, nil
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	tests := []struct {
		name                string
		ctx                 context.Context
		userId              uuid.UUID
		mockAlertRepository *MockAlertRepository
		expectedResponse    string
		expectedError       error
		expectedIsCalled    bool
		expectedUserId      uuid.UUID
	}{
		{
			name:                "alert",
			ctx:                 ctx,
			userId:              userId,
			mockAlertRepository: mockAlertRepositoryAlert,
			expectedResponse:    "За последние дни низкий уровень настроения и большая загрузка",
			expectedError:       nil,
			expectedIsCalled:    true,
			expectedUserId:      userId,
		},
		{
			name:                "not alert",
			ctx:                 ctx,
			userId:              userId,
			mockAlertRepository: mockAlertRepositoryNotAlert,
			expectedResponse:    "Все хорошо",
			expectedError:       nil,
			expectedIsCalled:    true,
			expectedUserId:      userId,
		},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alertService := NewAlertServcie(test.mockAlertRepository)
			response, err := alertService.GetLastSevenDays(test.ctx, test.userId)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error was - %v", test.expectedError)
			}
			if response != test.expectedResponse {
				fmt.Println(response, test.expectedResponse)
				t.Errorf("expected response was - %v", test.expectedResponse)
			}
			if test.mockAlertRepository.getLastSevenDaysFnIsCalled != test.expectedIsCalled {
				t.Errorf("get last seven days was not called")
			}
			if test.mockAlertRepository.userId != test.expectedUserId {
				fmt.Println(test.mockAlertRepository.userId, test.expectedUserId)
				t.Errorf("expected userId was - %v", test.expectedUserId)
			}
		})
	}
}

// Тест GetLastSevenDays - Провал (Err)
func TestGetLastSevenDaysErr(t *testing.T) {
	// preparing
	needError := errors.New("need error")
	mockAlertRepository := &MockAlertRepository{
		GetLastSevenDaysFn: func(ctx context.Context, userId uuid.UUID) ([]domain.Day, error) {
			return nil, needError
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	alertService := NewAlertServcie(mockAlertRepository)
	expectedError := needError

	// test
	response, err := alertService.GetLastSevenDays(ctx, userId)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error was - %v", expectedError)
	}
	if response != "" {
		t.Errorf("expected response was empty")
	}
	if !mockAlertRepository.getLastSevenDaysFnIsCalled {
		t.Errorf("get last seven days was not called")
	}
	if mockAlertRepository.userId != userId {
		t.Errorf("userId was expected - %v", userId)
	}
}

// Тест isSleepHoursLow
func TestIsSleepHoursLow(t *testing.T) {
	// preparing
	tests := []struct {
		name            string
		sleepHoursOne   float64
		sleepHoursTwo   float64
		sleepHoursThree float64
		expectedOk      bool
	}{
		{
			name:            "all bad",
			sleepHoursOne:   4.4,
			sleepHoursTwo:   3.3,
			sleepHoursThree: 2.2,
			expectedOk:      true,
		},
		{
			name:            "two bad",
			sleepHoursOne:   8.0,
			sleepHoursTwo:   5.5,
			sleepHoursThree: 4.4,
			expectedOk:      true,
		},
		{
			name:            "one bad",
			sleepHoursOne:   8.0,
			sleepHoursTwo:   7.7,
			sleepHoursThree: 5.5,
			expectedOk:      false,
		},
		{
			name:            "all good",
			sleepHoursOne:   8.0,
			sleepHoursTwo:   7.7,
			sleepHoursThree: 7.6,
			expectedOk:      false,
		},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := isSleepHoursLow(test.sleepHoursOne, test.sleepHoursTwo, test.sleepHoursThree)
			if ok != test.expectedOk {
				t.Errorf("expected ok was - %v", test.expectedOk)
			}
		})
	}
}

// Тест isLoadHigh
func TestIsLoadHigh(t *testing.T) {
	// preparing
	tests := []struct {
		name       string
		loadOne    int16
		loadTwo    int16
		loadThree  int16
		expectedOk bool
	}{
		{
			name:       "all bad",
			loadOne:    8,
			loadTwo:    7,
			loadThree:  6,
			expectedOk: true,
		},
		{
			name:       "two bad",
			loadOne:    5,
			loadTwo:    4,
			loadThree:  6,
			expectedOk: true,
		},
		{
			name:       "one bad",
			loadOne:    5,
			loadTwo:    4,
			loadThree:  3,
			expectedOk: false,
		},
		{
			name:       "all good",
			loadOne:    4,
			loadTwo:    3,
			loadThree:  2,
			expectedOk: false,
		},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := isLoadHigh(test.loadOne, test.loadTwo, test.loadThree)
			if ok != test.expectedOk {
				t.Errorf("expected ok was - %v", test.expectedOk)
			}
		})
	}
}
