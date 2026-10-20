package storage

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

func (s *Store) Save(url string) string {
	return ""
}

func (s *Store) Get(code string) (string, bool) {
	return "", false
}
