package service

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"github.com/google/uuid"
)

type IdempotencyService interface {
	Create(data entity.CreateIdempotency) (entity.Idempotency, error)
	Get(key uuid.UUID) (entity.Idempotency, error)
}

type idempotencyService struct {
	repository repository.IdempotencyRepository
}

func NewIdempotencyService(repository repository.IdempotencyRepository) IdempotencyService {
	return &idempotencyService{
		repository: repository,
	}
}

func (service *idempotencyService) Create(data entity.CreateIdempotency) (entity.Idempotency, error) {
	return service.repository.Create(data)
}

func (service *idempotencyService) Get(key uuid.UUID) (entity.Idempotency, error) {
	return service.repository.Get(key)
}
