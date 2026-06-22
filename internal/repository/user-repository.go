package repository

import (
	"database/sql"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateTx(tx *sql.Tx, data entity.CreateUser) (entity.User, error)
	GetByID(userID uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
}
