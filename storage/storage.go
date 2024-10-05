package storage

import (
	"encoding/gob"
	"os"
	"sync"
)

type Storage[T any] struct {
	m      sync.Mutex
	states map[int64]T
}

func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		states: make(map[int64]T),
	}
}

func (s *Storage[T]) Set(id int64, state T) {
	s.m.Lock()
	defer s.m.Unlock()
	s.states[id] = state
}

func (s *Storage[T]) Get(id int64) (T, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	state, ok := s.states[id]
	return state, ok
}

func (s *Storage[T]) DeleteUserState(id int64) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.states, id)
}

const fileName = "storage.data"

func (s *Storage[T]) Load() error {
	s.m.Lock()
	defer s.m.Unlock()
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return decoder.Decode(&s.states)
}

func (s *Storage[T]) Save() error {
	s.m.Lock()
	defer s.m.Unlock()
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	decoder := gob.NewEncoder(file)
	return decoder.Encode(s.states)
}
