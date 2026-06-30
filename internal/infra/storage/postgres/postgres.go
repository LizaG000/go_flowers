package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
)

type Storage struct {
	DB *sql.DB
}

func New(cfg config.Database) (*Storage, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: open database: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("%s: ping database: %w", op, err)
	}

	return &Storage{
		DB: db,
	}, nil
}
