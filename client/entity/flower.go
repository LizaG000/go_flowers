package entity

import (
	"time"

	"github.com/google/uuid"
)

type Flower struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Height      float64   `json:"height" db:"height"`
	Count       int64     `json:"count" db:"count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateFlower struct {
	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	Height      float64 `json:"height" db:"height"`
	Count       int64   `json:"count" db:"count"`
}
