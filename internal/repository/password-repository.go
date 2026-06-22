package repository

import (
	"database/sql"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"github.com/google/uuid"
)

type PasswordRepository interface {
	CreateTx(tx *sql.Tx, data entity.CreatePassword) error
	Get(userID uuid.UUID) (entity.Password, error)
}

type passwordRepository struct {
	db *sql.DB
}

func NewPasswordRepository(db *sql.DB) PasswordRepository {
	return &passwordRepository{
		db: db,
	}
}
func (repository *passwordRepository) CreateTx(
	tx *sql.Tx,
	data entity.CreatePassword,
) error {
	const query = `
		INSERT INTO passwords (
			user_id,
			password_hash
		)
		VALUES ($1, $2)
	`

	_, err := tx.Exec(
		query,
		data.UserID,
		data.Password,
	)
	if err != nil {
		return fmt.Errorf("postgres create password in transaction: %w", err)
	}

	return nil
}

func (repository *passwordRepository) Get(userID uuid.UUID) (entity.Password, error) {
	const query = `
		SELECT
			id,
			user_id,
			password_hash,
			created_at,
			updated_at
		FROM passwords WHERE user_id = $1
	`

	var password entity.Password

	err := repository.db.QueryRow(
		query,
		userID,
	).Scan(
		&password.ID,
		&password.UserID,
		&password.Password,
		&password.CreatedAt,
		&password.UpdatedAt,
	)

	if err != nil {
		return entity.Password{}, fmt.Errorf("postgres get password: %w", err)
	}

	return password, nil
}
