package repository

import "errors"

var ErrUniqueViolation = errors.New("unique violation")
var ErrNoRow = errors.New("no rows found")
