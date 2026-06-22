package repository

import (
	"database/sql"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"github.com/google/uuid"
)

type PasswordRepository interface {
	CreateTx(tx *sql.Tx, data entity.CreatePassword) error
	Get(userID uuid.UUID) (entity.Password, error)
}
