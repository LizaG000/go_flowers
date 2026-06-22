package entity

import (
	"time"

	"github.com/google/uuid"
)

type Password struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Password  string    `json:"password" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePassword struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Password string    `json:"password" db:"password_hash"`
}
