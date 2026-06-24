package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage"
	"github.com/google/uuid"
)

type IdempotencyRepository interface {
	Create(data entity.CreateIdempotency) (entity.Idempotency, error)
	Get(key uuid.UUID) (entity.Idempotency, error)
}

type idempotencyRepository struct {
	db *sql.DB
}

func NewIdempotencyRepository(db *sql.DB) IdempotencyRepository {
	return &idempotencyRepository{
		db: db,
	}
}

func (repository *idempotencyRepository) Create(data entity.CreateIdempotency) (entity.Idempotency, error) {
	const query = `
		INSERT INTO idempotency(
			key,
			status,
			response_code,
			response_body,
			payload_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING key, status, response_code, response_body, payload_hash, created_at, updated_at
	`

	var idempotency entity.Idempotency

	err := repository.db.QueryRow(
		query,
		data.Key,
		data.Status,
		data.ResponseCode,
		data.ResponseBody,
		data.PayloadHash,
	).Scan(
		&idempotency.Key,
		&idempotency.Status,
		&idempotency.ResponseCode,
		&idempotency.ResponseBody,
		&idempotency.PayloadHash,
		&idempotency.CreatedAt,
		&idempotency.UpdatedAt,
	)

	if err != nil {
		return entity.Idempotency{}, fmt.Errorf("postgres create idempotency in transaction: %w", err)
	}

	return idempotency, nil

}

func (repository *idempotencyRepository) Get(key uuid.UUID) (entity.Idempotency, error) {
	const query = `
		SELECT
			key,
			status,
			response_code,
			response_body,
			payload_hash,
			created_at,
			updated_at
		FROM idempotency WHERE key = $1
	`
	var idempotency entity.Idempotency

	err := repository.db.QueryRow(
		query,
		key,
	).Scan(
		&idempotency.Key,
		&idempotency.Status,
		&idempotency.ResponseCode,
		&idempotency.ResponseBody,
		&idempotency.PayloadHash,
		&idempotency.CreatedAt,
		&idempotency.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Idempotency{}, storage.ErrIdempotencyNotFound
	}
	if err != nil {
		return entity.Idempotency{}, fmt.Errorf("postgres get idempotency: %w", err)
	}

	return idempotency, nil
}
