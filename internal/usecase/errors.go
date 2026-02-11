package usecase

import "errors"

// users
var ErrUserExists = errors.New("user already exists")
var ErrUserNotExist = errors.New("user not exist")
var ErrWrongPassword = errors.New("wrong password")

// notes
var ErrWrongSleepHourValue = errors.New("wrong sleep hours value")
var ErrNoteAlreadyExists = errors.New("note already exusts")
