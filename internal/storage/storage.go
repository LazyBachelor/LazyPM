package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Storage[T any] struct {
	mu   sync.RWMutex
	path string
	Data *T
}

func NewJsonStorage[T any](path string, defaultData *T) *Storage[T] {
	return &Storage[T]{
		path: path,
		Data: defaultData,
	}
}

func (s *Storage[T]) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *Storage[T]) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bytes, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(bytes, s.Data)
}

func (s *Storage[T]) Init() error {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return s.Save()
	}
	return s.Load()
}

