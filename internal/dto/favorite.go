package dto

import "github.com/google/uuid"

type RequestCreateFavorite struct {
	FlowerID uuid.UUID `json:"flowerID"`
}
