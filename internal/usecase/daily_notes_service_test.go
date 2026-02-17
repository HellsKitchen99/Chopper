package usecase

import (
	"chopper/internal/domain"
	"chopper/internal/repository"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Моки
// Мок репозитория
type MockDailyNotesRepository struct {
	CreateNoteFn func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error
	// переданные аргументы
	createNoteFnIsCalled bool
	createNoteId         uuid.UUID
	createNoteUserId     uuid.UUID
	createNoteDate       time.Time
	createNoteMood       int16
	createNoteSleepHours float64
	createNoteLoad       int16

	ChangeMoodFn func(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error
	// переданные аргументы
	changeMoodFnIsCalled bool
	changeMoodUserId     uuid.UUID
	changeMoodDate       time.Time
	changeMoodMood       int16

	ChangeSleepHoursFn func(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error
	//переданные аргументы
	changeSleepHoursFnIsCalled bool
	changeSleepHoursUserId     uuid.UUID
	changeSleepHoursDate       time.Time
	changeSleepHoursSleepHours float64

	ChangeLoadFn func(ctx context.Context, userId uuid.UUID, date time.Time, load int16) error
	// переданные аргументы
	changeLoadFnIsCalled bool
	changeLoadUserId     uuid.UUID
	changeLoadDate       time.Time
	changeLoadLoad       int16
}

func (m *MockDailyNotesRepository) CreateNote(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
	m.createNoteFnIsCalled = true
	m.createNoteId = id
	m.createNoteUserId = userId
	m.createNoteDate = date
	m.createNoteMood = mood
	m.createNoteSleepHours = sleepHours
	m.createNoteLoad = load
	if m.CreateNoteFn != nil {
		return m.CreateNoteFn(ctx, id, userId, date, mood, sleepHours, load)
	}
	return nil
}

func (m *MockDailyNotesRepository) ChangeMood(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
	m.changeMoodFnIsCalled = true
	m.changeMoodUserId = userId
	m.changeMoodDate = date
	m.changeMoodMood = mood
	if m.ChangeMoodFn != nil {
		return m.ChangeMoodFn(ctx, userId, date, mood)
	}
	return nil
}

func (m *MockDailyNotesRepository) ChangeSleepHours(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
	m.changeSleepHoursFnIsCalled = true
	m.changeSleepHoursUserId = userId
	m.changeSleepHoursDate = date
	m.changeSleepHoursSleepHours = sleepHours
	if m.ChangeSleepHoursFn != nil {
		return m.ChangeSleepHoursFn(ctx, userId, date, sleepHours)
	}
	return nil
}

func (m *MockDailyNotesRepository) ChangeLoad(ctx context.Context, userId uuid.UUID, date time.Time, load int16) error {
	m.changeLoadFnIsCalled = true
	m.changeLoadUserId = userId
	m.changeLoadDate = date
	m.changeLoadLoad = load
	if m.ChangeLoadFn != nil {
		return m.ChangeLoadFn(ctx, userId, date, load)
	}
	return nil
}

// Мок генератора uuid
type MockUUIDGenerator struct {
	NewIdFn  func() uuid.UUID
	isCalled bool
	uuid     uuid.UUID
}

func (m *MockUUIDGenerator) NewId() uuid.UUID {
	m.isCalled = true
	if m.NewIdFn != nil {
		uuid := m.NewIdFn()
		m.uuid = uuid
		return uuid
	}
	return uuid.UUID{}
}

// Тест CreateNote - Успех
func TestCreateNoteSuccess(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return nil
		},
	}
	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	dailyNotesSevice := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	ctx := context.Background()
	userId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFront := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 5.5,
		Load:       5,
	}
	expectedId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedUserId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	expectedMood, expectedSleepHours, expectedLoad := int16(5), 5.5, int16(5)
	now := time.Now()
	expectedDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// test
	err := dailyNotesSevice.CreateNote(ctx, userId, dailyNoteFromFront)

	// assert
	if err != nil {
		t.Errorf("ошибки не ожидалось")
	}
	if !mockIdGenerator.isCalled {
		t.Errorf("генератор id не был вызван")
	}
	if mockIdGenerator.uuid != expectedId {
		t.Errorf("expected id - %v", expectedId)
	}
	if !mockDailyNotesRepository.createNoteFnIsCalled {
		t.Errorf("create note не был вызван")
	}
	if mockDailyNotesRepository.createNoteId != expectedId {
		t.Errorf("expected id - %v", expectedId)
	}
	if mockDailyNotesRepository.createNoteUserId != expectedUserId {
		t.Errorf("expected userId - %v", expectedUserId)
	}
	if mockDailyNotesRepository.createNoteDate != expectedDate {
		fmt.Println(mockDailyNotesRepository.createNoteDate, expectedDate)
		t.Errorf("expected date - %v", expectedDate)
	}
	if mockDailyNotesRepository.createNoteMood != expectedMood {
		t.Errorf("expected mood - %v", expectedMood)
	}
	if mockDailyNotesRepository.createNoteSleepHours != expectedSleepHours {
		t.Errorf("expected sleep hours - %v", expectedSleepHours)
	}
	if mockDailyNotesRepository.createNoteLoad != expectedLoad {
		t.Errorf("expected load - %v", expectedLoad)
	}
}

// Тест CreateNote - Провал (Невалидный mood)
func TestCreateNoteFailureInvalidMood(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return nil
		},
	}
	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	ctx := context.Background()
	userId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFrontNegativeMood := domain.DailyNoteFromFront{
		Mood:       -1,
		SleepHours: 5.5,
		Load:       5,
	}
	dailyNoteFromFrontOverMood := domain.DailyNoteFromFront{
		Mood:       11,
		SleepHours: 5.5,
		Load:       5,
	}
	tests := []struct {
		name               string
		ctx                context.Context
		userId             uuid.UUID
		dailyNoteFromFront domain.DailyNoteFromFront
		expectedError      error
	}{
		{name: "negative mood", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontNegativeMood, expectedError: ErrWrongMoodValue},
		{name: "over mood", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontOverMood, expectedError: ErrWrongMoodValue},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := dailyNotesService.CreateNote(test.ctx, test.userId, test.dailyNoteFromFront)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error - %v", test.expectedError)
			}
			if !mockIdGenerator.isCalled {
				t.Errorf("генератор id не был вызван")
			}
			if mockDailyNotesRepository.createNoteFnIsCalled {
				t.Errorf("create note был вызван")
			}
		})
	}
}

// Тест CreateNote - Провал (Невалидный sleepHours)
func TestCreateNoteInvalidSleepHours(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return nil
		},
	}

	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFrontNegativeSleepHours := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: -1.1,
		Load:       5,
	}
	dailyNoteFromFrontOverSleepHours := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 11.1,
		Load:       5,
	}
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	tests := []struct {
		name               string
		ctx                context.Context
		userId             uuid.UUID
		dailyNoteFromFront domain.DailyNoteFromFront
		expectedError      error
	}{
		{name: "negative sleep hours", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontNegativeSleepHours, expectedError: ErrWrongSleepHourValue},
		{name: "over sleep hours", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontOverSleepHours, expectedError: ErrWrongSleepHourValue},
	}

	// test ++ assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := dailyNotesService.CreateNote(test.ctx, test.userId, test.dailyNoteFromFront)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error - %v", test.expectedError)
			}
			if !mockIdGenerator.isCalled {
				t.Errorf("генератор id не был вызван")
			}
			if mockDailyNotesRepository.createNoteFnIsCalled {
				t.Error("create note был вызван")
			}
		})
	}
}

// Тест CreateNote - Провал (Невалидный load)
func TestCreateNoteInvalidLoad(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return nil
		},
	}
	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFrontNegativeLoad := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 5.5,
		Load:       -1,
	}
	dailyNoteFromFrontOverLoad := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 5.5,
		Load:       11,
	}
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	tests := []struct {
		name               string
		ctx                context.Context
		userId             uuid.UUID
		dailyNoteFromFront domain.DailyNoteFromFront
		expectedError      error
	}{
		{name: "negative load", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontNegativeLoad, expectedError: ErrWrongLoadValue},
		{name: "negative load", ctx: ctx, userId: userId, dailyNoteFromFront: dailyNoteFromFrontOverLoad, expectedError: ErrWrongLoadValue},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := dailyNotesService.CreateNote(test.ctx, test.userId, test.dailyNoteFromFront)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error - %v", test.expectedError)
			}
			if !mockIdGenerator.isCalled {
				t.Errorf("генератор id не был вызван")
			}
			if mockDailyNotesRepository.createNoteFnIsCalled {
				t.Errorf("create note был вызван")
			}
		})
	}

}

// Тест CreateNote - Провал (ErrNoteAlreadyExists)
func TestCreateNoteErrNoteAlreadyExists(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return repository.ErrUniqueViolation
		},
	}
	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFront := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 5.5,
		Load:       5,
	}
	dailyNoteService := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	expectedUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	expectedMood, expectedSleepHours, expectedLoad := int16(5), 5.5, int16(5)
	expectedError := ErrNoteAlreadyExists

	// test
	err := dailyNoteService.CreateNote(ctx, userId, dailyNoteFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error - %v", expectedError)
	}
	if !mockIdGenerator.isCalled {
		t.Error("генератор id не был вызван")
	}
	if mockIdGenerator.uuid != expectedUUID {
		t.Errorf("expected uuid - %v", expectedUUID)
	}
	if !mockDailyNotesRepository.createNoteFnIsCalled {
		t.Errorf("create note не был вызван")
	}
	if mockDailyNotesRepository.createNoteId != expectedUUID {
		t.Errorf("expected id - %v", expectedUUID)
	}
	if mockDailyNotesRepository.createNoteUserId != userId {
		t.Errorf("expected user id - %v", userId)
	}
	if mockDailyNotesRepository.createNoteDate != expectedDate {
		t.Errorf("expected date - %v", expectedDate)
	}
	if mockDailyNotesRepository.createNoteMood != expectedMood {
		t.Errorf("expected mood - %v", expectedMood)
	}
	if mockDailyNotesRepository.createNoteSleepHours != expectedSleepHours {
		t.Errorf("expected sleep hours - %v", expectedSleepHours)
	}
	if mockDailyNotesRepository.createNoteLoad != expectedLoad {
		t.Errorf("expected load - %v", expectedLoad)
	}
}

// Тест CreateNote - Провал (err)
func TestCreateNoteErr(t *testing.T) {
	needError := errors.New("need error")
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		CreateNoteFn: func(ctx context.Context, id, userId uuid.UUID, date time.Time, mood int16, sleepHours float64, load int16) error {
			return needError
		},
	}
	mockIdGenerator := &MockUUIDGenerator{
		NewIdFn: func() uuid.UUID {
			return uuid.MustParse("11111111-1111-1111-1111-111111111111")
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("22222222-2222-2222-2222-222222222222")
	dailyNoteFromFront := domain.DailyNoteFromFront{
		Mood:       5,
		SleepHours: 5.5,
		Load:       5,
	}
	dailyNoteService := NewDailyNotesService(mockDailyNotesRepository, mockIdGenerator)
	expectedUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	expectedMood, expectedSleepHours, expectedLoad := int16(5), 5.5, int16(5)
	expectedError := needError

	// test
	err := dailyNoteService.CreateNote(ctx, userId, dailyNoteFromFront)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error - %v", expectedError)
	}
	if !mockIdGenerator.isCalled {
		t.Error("генератор id не был вызван")
	}
	if mockIdGenerator.uuid != expectedUUID {
		t.Errorf("expected uuid - %v", expectedUUID)
	}
	if !mockDailyNotesRepository.createNoteFnIsCalled {
		t.Errorf("create note не был вызван")
	}
	if mockDailyNotesRepository.createNoteId != expectedUUID {
		t.Errorf("expected id - %v", expectedUUID)
	}
	if mockDailyNotesRepository.createNoteUserId != userId {
		t.Errorf("expected user id - %v", userId)
	}
	if mockDailyNotesRepository.createNoteDate != expectedDate {
		t.Errorf("expected date - %v", expectedDate)
	}
	if mockDailyNotesRepository.createNoteMood != expectedMood {
		t.Errorf("expected mood - %v", expectedMood)
	}
	if mockDailyNotesRepository.createNoteSleepHours != expectedSleepHours {
		t.Errorf("expected sleep hours - %v", expectedSleepHours)
	}
	if mockDailyNotesRepository.createNoteLoad != expectedLoad {
		t.Errorf("expected load - %v", expectedLoad)
	}
}

// Тест ChangeMood - Успех
func TestChangeMoodSuccess(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeMoodFn: func(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
			return nil
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	mood := int16(5)
	dailyNoteService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedResponse := "mood успешно изменен"
	/*
		changeMoodFnIsCalled bool
		changeMoodUserId     uuid.UUID
		changeMoodDate       time.Time
		changeMoodMood       int16
	*/

	// test
	response, err := dailyNoteService.ChangeMood(ctx, userId, date, mood)

	// assert
	if err != nil {
		t.Errorf("ошибки не ожидалось")
	}
	if response != expectedResponse {
		t.Errorf("expected response - %v", expectedResponse)
	}
	if !mockDailyNotesRepository.changeMoodFnIsCalled {
		t.Errorf("change mood не был вызван")
	}
	if mockDailyNotesRepository.changeMoodUserId != userId {
		t.Errorf("expected userId - %v", userId)
	}
	if mockDailyNotesRepository.changeMoodDate != date {
		t.Errorf("expected date - %v", date)
	}
	if mockDailyNotesRepository.changeMoodMood != mood {
		t.Errorf("expected mood - %v", mood)
	}
}

// Тест ChangeMood - Провал (ErrWrongMoodValue)
func TestChangeMoodErrWrongMoodValue(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeMoodFn: func(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
			return nil
		},
	}
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	ctx := context.Background()
	userId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	moodNegative := int16(-1)
	moodOver := int16(11)
	expectedError := ErrWrongMoodValue
	tests := []struct {
		name          string
		ctx           context.Context
		userId        uuid.UUID
		date          time.Time
		mood          int16
		expectedError error
	}{
		{name: "negative mood", ctx: ctx, userId: userId, date: date, mood: moodNegative, expectedError: expectedError},
		{name: "over mood", ctx: ctx, userId: userId, date: date, mood: moodOver, expectedError: expectedError},
	}

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dailyNotesService.ChangeMood(test.ctx, test.userId, test.date, test.mood)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error - %v", expectedError)
			}
			if result != "" {
				t.Errorf("ожидался пустой результат")
			}
			if mockDailyNotesRepository.changeMoodFnIsCalled {
				t.Errorf("change mood был вызван")
			}
		})
	}
}

// Тест ChangeMood - Провал (ErrNoteNotExists)
func TestChangeMoodErrNoteNotExists(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeMoodFn: func(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
			return repository.ErrDailyEntryNotFound
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	mood := int16(5)
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedError := ErrNoteNotExists

	// test
	result, err := dailyNotesService.ChangeMood(ctx, userId, date, mood)

	// assert
	if !errors.Is(err, expectedError) {
		fmt.Println(err, expectedError)
		t.Errorf("expected error - %v", expectedError)
	}
	if result != "" {
		t.Errorf("result was expected empty")
	}
	if !mockDailyNotesRepository.changeMoodFnIsCalled {
		t.Errorf("change mood was not called")
	}
	if mockDailyNotesRepository.changeMoodUserId != userId {
		t.Errorf("expected userId - %v", userId)
	}
	if mockDailyNotesRepository.changeMoodDate != date {
		t.Errorf("expected date - %v", date)
	}
	if mockDailyNotesRepository.changeMoodMood != mood {
		t.Errorf("expected mood - %v", mood)
	}
}

// Тест ChangeMood - Провал (Err)
func TestChangeMoodErr(t *testing.T) {
	// preparing
	needError := errors.New("need error")
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeMoodFn: func(ctx context.Context, userId uuid.UUID, date time.Time, mood int16) error {
			return needError
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	mood := int16(5)
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedError := needError

	// test
	result, err := dailyNotesService.ChangeMood(ctx, userId, date, mood)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error - %v", expectedError)
	}
	if result != "" {
		t.Errorf("result was expected empty")
	}
	if !mockDailyNotesRepository.changeMoodFnIsCalled {
		t.Error("change mood was not called")
	}
	if mockDailyNotesRepository.changeMoodUserId != userId {
		t.Errorf("expected userId - %v", userId)
	}
	if mockDailyNotesRepository.changeMoodDate != date {
		t.Errorf("expected date - %v", date)
	}
	if mockDailyNotesRepository.changeMoodMood != mood {
		t.Errorf("expected mood - %v", mood)
	}
}

// Тест ChangeSleepHours - Успех
func TestChangeSleepHoursSuccess(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeSleepHoursFn: func(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
			return nil
		},
	}
	ctx := context.Background()
	userId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sleepHours := 5.5
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedResponse := "sleep hours успешно изменен"

	// test
	response, err := dailyNotesService.ChangeSleepHours(ctx, userId, date, sleepHours)

	// assert
	if err != nil {
		t.Errorf("error was expected nil")
	}
	if response != expectedResponse {
		t.Errorf("response was expected - %v", expectedResponse)
	}
	if !mockDailyNotesRepository.changeSleepHoursFnIsCalled {
		t.Errorf("change sleep hours не был вызван")
	}
	if mockDailyNotesRepository.changeSleepHoursUserId != userId {
		t.Errorf("userId was expected - %v", userId)
	}
	if mockDailyNotesRepository.changeSleepHoursDate != date {
		t.Errorf("date was expected - %v", date)
	}
	if mockDailyNotesRepository.changeSleepHoursSleepHours != sleepHours {
		t.Errorf("sleep hours was expected - %v", sleepHours)
	}
}

// Тест ChangeSleepHours - Провал (ErrWrongSleepHourValue)
func TestChangeSleepHoursErrWrongSleepHourValue(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeSleepHoursFn: func(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
			return nil
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sleepHoursNegative := -1.1
	sleepHoursOver := 11.1
	expectedError := ErrWrongSleepHourValue
	tests := []struct {
		name          string
		ctx           context.Context
		userId        uuid.UUID
		date          time.Time
		sleepHours    float64
		expectedError error
	}{
		{name: "negative sleep hours", ctx: ctx, userId: userId, date: date, sleepHours: sleepHoursNegative, expectedError: expectedError},
		{name: "over sleep hours", ctx: ctx, userId: userId, date: date, sleepHours: sleepHoursOver, expectedError: expectedError},
	}
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)

	// test + assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dailyNotesService.ChangeSleepHours(test.ctx, test.userId, test.date, test.sleepHours)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error - %v", test.expectedError)
			}
			if result != "" {
				t.Errorf("result was expected empty")
			}
			if mockDailyNotesRepository.changeSleepHoursFnIsCalled {
				t.Errorf("change sleep hours был вызван")
			}
		})
	}
}

// Тест ChangeSleepHours - Провал (ErrNoteNotExists)
func TestChangeSleepHoursErrNoteNotExists(t *testing.T) {
	// preparing
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeSleepHoursFn: func(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
			return repository.ErrDailyEntryNotFound
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sleepHours := 5.5
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedError := ErrNoteNotExists

	// test
	response, err := dailyNotesService.ChangeSleepHours(ctx, userId, date, sleepHours)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error - %v", expectedError)
	}
	if response != "" {
		t.Errorf("response was expected empty")
	}
	if !mockDailyNotesRepository.changeSleepHoursFnIsCalled {
		t.Errorf("change sleep hours не был вызван")
	}
	if mockDailyNotesRepository.changeSleepHoursUserId != userId {
		t.Errorf("expected userId - %v", userId)
	}
	if mockDailyNotesRepository.changeSleepHoursDate != date {
		t.Errorf("expected date - %v", date)
	}
	if mockDailyNotesRepository.changeSleepHoursSleepHours != sleepHours {
		t.Errorf("expected sleep hours - %v", sleepHours)
	}
}

// Тест ChangeSleepHours - Провал (Err)
func TestChangeSleepHoursErr(t *testing.T) {
	// preparing
	needError := errors.New("need error")
	mockDailyNotesRepository := &MockDailyNotesRepository{
		ChangeSleepHoursFn: func(ctx context.Context, userId uuid.UUID, date time.Time, sleepHours float64) error {
			return needError
		},
	}
	ctx, userId := context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sleepHours := 5.5
	dailyNotesService := NewDailyNotesService(mockDailyNotesRepository, nil)
	expectedError := needError

	// test
	response, err := dailyNotesService.ChangeSleepHours(ctx, userId, date, sleepHours)

	// assert
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error - %v", expectedError)
	}
	if response != "" {
		t.Errorf("response was expected empty")
	}
	if !mockDailyNotesRepository.changeSleepHoursFnIsCalled {
		t.Errorf("change sleep hours не был вызван")
	}
	if mockDailyNotesRepository.changeSleepHoursUserId != userId {
		t.Errorf("userId was expected - %v", userId)
	}
	if mockDailyNotesRepository.changeSleepHoursDate != date {
		t.Errorf("date was expected - %v", date)
	}
	if mockDailyNotesRepository.changeSleepHoursSleepHours != sleepHours {
		t.Errorf("sleep hours was expected - %v", sleepHours)
	}
}
