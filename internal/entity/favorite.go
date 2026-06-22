package entity

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FlowerID  uuid.UUID `json:"flower_id" db:"flower_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateFavorite struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	FlowerID uuid.UUID `json:"flower_id" db:"flower_id"`
}

type FavoriteFlower struct {
	ID          uuid.UUID `json:"favorite_id" db:"favorite_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	FlowerID    uuid.UUID `json:"flower_id" db:"flower_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Height      float64   `json:"height" db:"height"`
	Count       int64     `json:"count" db:"count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
