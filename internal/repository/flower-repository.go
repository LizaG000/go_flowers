package repository

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"github.com/google/uuid"
)

type FlowerRepository interface {
	Create(data entity.CreateFlower) (entity.Flower, error)
	GetAll() ([]entity.Flower, error)
	Update(id uuid.UUID, data entity.UpdateFlower) (entity.Flower, error)
	Delete(id uuid.UUID) (entity.Flower, error)
}
