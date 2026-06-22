package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FavoriteRepository interface {
	Create(data entity.CreateFavorite) (entity.Favorite, error)
	GetByUserID(userID uuid.UUID) ([]entity.FavoriteFlower, error)
	Delete(favoriteID uuid.UUID) (entity.Favorite, error)
}

type favoriteRepository struct {
	db *sql.DB
}

func NewFavoriteREpository(db *sql.DB) FavoriteRepository {
	return &favoriteRepository{
		db: db,
	}
}

func (repository *favoriteRepository) Create(
	data entity.CreateFavorite,
) (entity.Favorite, error) {
	const query = `
	INSERT INTO favorites (
		user_id,
		flower_id
	)
	VALUES($1, $2)
	RETURNING
		id,
		user_id,
		flower_id,
		created_at,
		updated_at
	`

	var favorite entity.Favorite

	err := repository.db.QueryRow(
		query,
		data.UserID,
		data.FlowerID,
	).Scan(
		&favorite.ID,
		&favorite.UserID,
		&favorite.FlowerID,
		&favorite.CreatedAt,
		&favorite.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return entity.Favorite{}, storage.ErrFavoriteExists
		}

		return entity.Favorite{}, fmt.Errorf("postgres create favorite in transaction: %w", err)
	}

	return favorite, nil
}

func (repository *favoriteRepository) GetByUserID(
	userID uuid.UUID,
) ([]entity.FavoriteFlower, error) {
	const query = `
		SELECT
			favorites.id,
			favorites.user_id,
			favorites.flower_id,
			flowers.title,
			flowers.description,
			flowers.price,
			flowers.height,
			flowers.count,
			favorites.created_at,
			favorites.updated_at
		FROM favorites
		JOIN flowers ON favorites.flower_id = flowers.id
		WHERE favorites.user_id = $1
	`

	rows, err := repository.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("postgres get favorites: %w", err)
	}
	defer rows.Close()

	favorites := make([]entity.FavoriteFlower, 0)

	for rows.Next() {
		var favorite entity.FavoriteFlower

		err := rows.Scan(
			&favorite.ID,
			&favorite.UserID,
			&favorite.FlowerID,
			&favorite.Title,
			&favorite.Description,
			&favorite.Price,
			&favorite.Height,
			&favorite.Count,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres scan favorite: %w", err)
		}

		favorites = append(favorites, favorite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres iterate favorites: %w", err)
	}

	return favorites, nil
}

func (repository *favoriteRepository) Delete(
	favoriteID uuid.UUID,
) (entity.Favorite, error) {
	const query = `
	DELETE FROM favorites
	WHERE id = $1
	RETURNING
		id,
		user_id,
		flower_id,
		created_at,
		updated_at
	`

	var favorite entity.Favorite

	err := repository.db.QueryRow(
		query,
		favoriteID,
	).Scan(
		&favorite.ID,
		&favorite.UserID,
		&favorite.FlowerID,
		&favorite.CreatedAt,
		&favorite.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Favorite{}, storage.ErrFlowerNotFound
	}

	if err != nil {
		return entity.Favorite{}, fmt.Errorf("postgres delete favorites: %w", err)
	}

	return favorite, nil
}
