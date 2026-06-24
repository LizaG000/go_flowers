package storage

import "errors"

var (
	ErrFlowerNotFound = errors.New("flower not found")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrIdempotencyNotFound      = errors.New("idempotency not found")
	ErrIdempotencyAlreadyExists = errors.New("user already exists")

	ErrPasswordNotFound = errors.New("password not found")

	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrFavoriteExists   = errors.New("flower already in favorites")
)
