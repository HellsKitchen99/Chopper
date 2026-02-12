package usecase

import "errors"

// users
var ErrUserExists = errors.New("user already exists")
var ErrUserNotExist = errors.New("user not exist")
var ErrWrongPassword = errors.New("wrong password")

// notes
var ErrWrongMoodValue = errors.New("wrong mood value")
var ErrWrongSleepHourValue = errors.New("wrong sleep hours value")
var ErrWrongLoadValue = errors.New("wrong load value")
var ErrNoteAlreadyExists = errors.New("note already exusts")
var ErrNoteNotExists = errors.New("note not exists")
