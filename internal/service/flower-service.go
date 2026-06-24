package service

import (
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"github.com/google/uuid"
)

type FlowerService interface {
	Create(data entity.CreateFlower) (entity.Flower, error)
	GetAll() ([]entity.Flower, error)
	Update(id uuid.UUID, data entity.UpdateFlower) (entity.Flower, error)
	Delete(id uuid.UUID) (entity.Flower, error)
}

type flowerService struct {
	repository repository.FlowerRepository
}

func NewFlowerService(repository repository.FlowerRepository) FlowerService {
	return &flowerService{
		repository: repository,
	}
}

func (service *flowerService) Create(data entity.CreateFlower) (entity.Flower, error) {
	if data.Title == "" {
		return entity.Flower{}, fmt.Errorf("название цветка не может быть пустым")
	}

	if data.Price < 0 {
		return entity.Flower{}, fmt.Errorf("цена не может быть отрицательной")
	}

	if data.Height < 0 {
		return entity.Flower{}, fmt.Errorf("высота не может быть отрицательной")
	}

	if data.Count < 0 {
		return entity.Flower{}, fmt.Errorf("количество не может быть отрицательным")
	}

	return service.repository.Create(data)
}

func (service *flowerService) GetAll() ([]entity.Flower, error) {

	return service.repository.GetAll()
}

func (service *flowerService) Update(
	id uuid.UUID,
	data entity.UpdateFlower,
) (entity.Flower, error) {
	if data.Price != nil && *data.Price < 0 {
		return entity.Flower{}, fmt.Errorf("цена не может быть отрицательной")
	}

	if data.Height != nil && *data.Height < 0 {
		return entity.Flower{}, fmt.Errorf("высота не может быть отрицательной")
	}

	if data.Count != nil && *data.Count < 0 {
		return entity.Flower{}, fmt.Errorf("количество не может быть отрицательным")
	}

	return service.repository.Update(id, data)
}

func (service *flowerService) Delete(id uuid.UUID) (entity.Flower, error) {
	return service.repository.Delete(id)
}
