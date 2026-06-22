package service

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	GetByID(userID uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
}

type userService struct {
	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) GetByID(userID uuid.UUID) (entity.User, error) {
	user, _ := service.userRepository.GetByID(userID)
	return user, nil
}
func (service *userService) GetByEmail(email string) (entity.User, error) {
	user, _ := service.userRepository.GetByEmail(email)
	return user, nil
}
