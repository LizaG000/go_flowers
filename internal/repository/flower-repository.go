package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/storage"
	"github.com/google/uuid"
)

type FlowerRepository interface {
	Create(data entity.CreateFlower) (entity.Flower, error)
	GetAll() ([]entity.Flower, error)
	Update(id uuid.UUID, data entity.UpdateFlower) (entity.Flower, error)
	Delete(id uuid.UUID) (entity.Flower, error)
}

type flowerRepository struct {
	db *sql.DB
}

func NewFlowerRepository(db *sql.DB) FlowerRepository {
	return &flowerRepository{
		db: db,
	}
}

func (repository *flowerRepository) Create(
	data entity.CreateFlower,
) (entity.Flower, error) {
	const query = `
		INSERT INTO flowers (
			title,
			description,
			price,
			height,
			count
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			title,
			description,
			price,
			height,
			count,
			created_at,
			updated_at
	`

	var flower entity.Flower

	err := repository.db.QueryRow(
		query,
		data.Title,
		data.Description,
		data.Price,
		data.Height,
		data.Count,
	).Scan(
		&flower.ID,
		&flower.Title,
		&flower.Description,
		&flower.Price,
		&flower.Height,
		&flower.Count,
		&flower.CreatedAt,
		&flower.UpdatedAt,
	)
	if err != nil {
		return entity.Flower{}, fmt.Errorf("postgres create flower: %w", err)
	}

	return flower, nil
}

func (repository *flowerRepository) GetAll() ([]entity.Flower, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			price,
			height,
			count,
			created_at,
			updated_at
		FROM flowers
		ORDER BY created_at DESC
	`

	rows, err := repository.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("postgres get flowers: %w", err)
	}
	defer rows.Close()

	flowers := make([]entity.Flower, 0)

	for rows.Next() {
		var flower entity.Flower

		if err := rows.Scan(
			&flower.ID,
			&flower.Title,
			&flower.Description,
			&flower.Price,
			&flower.Height,
			&flower.Count,
			&flower.CreatedAt,
			&flower.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres scan flower: %w", err)
		}

		flowers = append(flowers, flower)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres iterate flowers: %w", err)
	}

	return flowers, nil
}

func (repository *flowerRepository) Update(
	id uuid.UUID,
	data entity.UpdateFlower,
) (entity.Flower, error) {
	const query = `
		UPDATE flowers
		SET
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			price = COALESCE($3, price),
			height = COALESCE($4, height),
			count = COALESCE($5, count),
			updated_at = NOW()
		WHERE id = $6
		RETURNING
			id,
			title,
			description,
			price,
			height,
			count,
			created_at,
			updated_at
	`

	var flower entity.Flower

	err := repository.db.QueryRow(
		query,
		data.Title,
		data.Description,
		data.Price,
		data.Height,
		data.Count,
		id,
	).Scan(
		&flower.ID,
		&flower.Title,
		&flower.Description,
		&flower.Price,
		&flower.Height,
		&flower.Count,
		&flower.CreatedAt,
		&flower.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Flower{}, storage.ErrFlowerNotFound
	}

	if err != nil {
		return entity.Flower{}, fmt.Errorf("postgres update flower: %w", err)
	}

	return flower, nil
}

func (repository *flowerRepository) Delete(
	id uuid.UUID,
) (entity.Flower, error) {
	const query = `
		DELETE FROM flowers
		WHERE id = $1
		RETURNING
			id,
			title,
			description,
			price,
			height,
			count,
			created_at,
			updated_at
	`

	var flower entity.Flower

	err := repository.db.QueryRow(query, id).Scan(
		&flower.ID,
		&flower.Title,
		&flower.Description,
		&flower.Price,
		&flower.Height,
		&flower.Count,
		&flower.CreatedAt,
		&flower.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Flower{}, storage.ErrFlowerNotFound
	}

	if err != nil {
		return entity.Flower{}, fmt.Errorf("postgres delete flower: %w", err)
	}

	return flower, nil
}
