package postgres

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/StepanDubrovin/shortener/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func New(dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(url string) (string, error) {
	query := `INSERT INTO urls (url) VALUES ($1) RETURNING id`

	var id int64

	err := s.db.QueryRow(query, url).Scan(&id)

	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func (s *Store) Get(code string) (string, error) {
	id, err := strconv.ParseInt(code, 10, 64)

	if err != nil || id <= 0 {
		return "", storage.ErrNotFound
	}

	query := `SELECT url FROM urls WHERE id = $1`

	var url string

	err = s.db.QueryRow(query, id).Scan(&url)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
