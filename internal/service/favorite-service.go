package service

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"github.com/google/uuid"
)

type FavoriteService interface {
	Create(data entity.CreateFavorite) (entity.Favorite, error)
	GetByUserID(userID uuid.UUID) (entity.FavoriteFlower, error)
	Delete(favoriteID uuid.UUID) (entity.Favorite, error)
}

type favoriteService struct {
	repository repository.FavoriteRepository
}

func NewFavoriteService(repository repository.FavoriteRepository) FavoriteService {
	return &favoriteService{
		repository: repository,
	}
}

func (service *favoriteService) Create(data entity.CreateFavorite) (entity.Favorite, error) {
	return service.repository.Create(data)
}

func (service *favoriteService) GetByUserID(userID uuid.UUID) (entity.FavoriteFlower, error) {
	return service.repository.GetByUserID(userID)
}

func (service *favoriteService) Delete(favoriteID uuid.UUID) (entity.Favorite, error) {
	return service.repository.Delete(favoriteID)
}
