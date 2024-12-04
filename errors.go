package alchemy

import "errors"

var ErrInvalidCredentials error = errors.New("invalid credentials")
var ErrNotFound error = errors.New("record not found")
var ErrDuplicateRecord error = errors.New("record already exist")
