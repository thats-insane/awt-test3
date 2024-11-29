package data

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("duplicate email")
var ErrEditConflict = errors.New("edit conflict")
