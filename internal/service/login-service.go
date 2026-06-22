package service

import (
	"database/sql"
	"fmt"
	"time"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/security"
)

type LoginService interface {
	Login(email string, password string) (string, error)
	Registration(data entity.CreateUser, password string) (entity.User, error)
}

type loginService struct {
	db                 *sql.DB
	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository
	auth               config.Auth
}

func NewLoginService(
	db *sql.DB,
	userRepository repository.UserRepository,
	passwordRepository repository.PasswordRepository,
	auth config.Auth) LoginService {
	return &loginService{
		db:                 db,
		userRepository:     userRepository,
		passwordRepository: passwordRepository,
		auth:               auth,
	}
}

func (service *loginService) Login(
	email string,
	password string,
) (string, error) {
	user, err := service.userRepository.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	passwordData, err := service.passwordRepository.Get(user.ID)
	if err != nil {
		return "", fmt.Errorf("get password hash: %w", err)
	}

	if err := security.ComparePassword(
		passwordData.Password,
		password,
	); err != nil {
		return "", fmt.Errorf("неверный email или пароль")
	}

	token, err := security.GenerateToken(
		user.ID,
		service.auth.PrivateKeyPath,
		service.auth.TokenTTL,
	)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (service *loginService) Registration(
	data entity.CreateUser,
	password string,
) (entity.User, error) {
	now := time.Now()

	age := now.Year() - data.BirthDate.Year()

	if now.Month() < data.BirthDate.Month() ||
		(now.Month() == data.BirthDate.Month() &&
			now.Day() < data.BirthDate.Day()) {
		age--
	}

	if age < 14 {
		return entity.User{}, fmt.Errorf("возраст не может быть младше 14 лет")
	}

	tx, err := service.db.Begin()
	if err != nil {
		return entity.User{}, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	user, err := service.userRepository.CreateTx(tx, data)
	if err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}

	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return entity.User{}, fmt.Errorf("hash password: %w", err)
	}

	err = service.passwordRepository.CreateTx(tx, entity.CreatePassword{
		UserID:   user.ID,
		Password: passwordHash,
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("create password: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return entity.User{}, fmt.Errorf("commit transaction: %w", err)
	}

	return user, nil
}
