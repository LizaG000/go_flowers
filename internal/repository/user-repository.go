package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateTx(tx *sql.Tx, data entity.CreateUser) (entity.User, error)
	GetByID(userID uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repository *userRepository) CreateTx(
	tx *sql.Tx,
	data entity.CreateUser,
) (entity.User, error) {
	const query = `
		INSERT INTO users (
			first_name,
			second_name,
			last_name,
			email,
			birth_date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			first_name,
			second_name,
			last_name,
			email,
			birth_date,
			created_at,
			updated_at
	`

	var user entity.User

	err := tx.QueryRow(
		query,
		data.FirstName,
		data.SecondName,
		data.LastName,
		data.Email,
		data.BirthDate,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.BirthDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, storage.ErrUserAlreadyExists
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("postgres create user in transaction: %w", err)
	}

	return user, nil
}

func (repository *userRepository) GetByID(userID uuid.UUID) (entity.User, error) {
	const query = `
		SELECT
			id,
			first_name,
			second_name,
			last_name,
			email,
			birth_date,
			created_at,
			updated_at
		FROM users WHERE id = $1
	`
	var user entity.User

	err := repository.db.QueryRow(
		query,
		userID,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.BirthDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
}

func (repository *userRepository) GetByEmail(email string) (entity.User, error) {
	const query = `
		SELECT
			id,
			first_name,
			second_name,
			last_name,
			email,
			birth_date,
			created_at,
			updated_at
		FROM users WHERE email = $1
	`
	var user entity.User

	err := repository.db.QueryRow(
		query,
		email,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.BirthDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
}
