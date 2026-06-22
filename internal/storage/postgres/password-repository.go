package postgres

import (
	"database/sql"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"github.com/google/uuid"
)

type PasswordRepository struct {
	db *sql.DB
}

var _ repository.PasswordRepository = (*PasswordRepository)(nil)

func NewPasswordRepository(db *sql.DB) *PasswordRepository {
	return &PasswordRepository{
		db: db,
	}
}
func (repository *PasswordRepository) CreateTx(
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

func (repository *PasswordRepository) Get(userID uuid.UUID) (entity.Password, error) {
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
