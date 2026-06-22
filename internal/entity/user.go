package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	SecondName string    `json:"second_name" db:"second_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Email      string    `json:"email" db:"email"`
	BirthDate  time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUser struct {
	FirstName  string    `json:"first_name" db:"first_name"`
	SecondName string    `json:"second_name" db:"second_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Email      string    `json:"email" db:"email"`
	BirthDate  time.Time `json:"birth_date" db:"birth_date"`
}
