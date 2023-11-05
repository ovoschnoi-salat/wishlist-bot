package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrNoSuchUser   = errors.New("no user found")
	ErrWrongArg     = errors.New("wrong argument passed")
	ErrReservedWish = errors.New("wish is already reserved")
	ErrNoReserver   = errors.New("wish is not reserved")
)
