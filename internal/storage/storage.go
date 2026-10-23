package storage

import (
	"errors"
	"strconv"
)

var ErrNotFound = errors.New("url not found")

type Store struct {
	urls   map[string]string
	nextID int
}

func New() *Store {
	return &Store{
		urls:   make(map[string]string),
		nextID: 1,
	}
}

func (s *Store) Save(url string) (string, error) {
	id := strconv.Itoa(s.nextID)
	s.nextID++

	s.urls[id] = url

	return id, nil
}

func (s *Store) Get(code string) (string, error) {
	value, ok := s.urls[code]

	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}
