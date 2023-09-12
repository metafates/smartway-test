package hashset

import (
	"encoding/json"

	"github.com/samber/lo"
)

type Set[T comparable] struct {
	m map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

func (s *Set[T]) Put(value T) {
	s.m[value] = struct{}{}
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Set[T]) Remove(value T) {
	delete(s.m, value)
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var items []T

	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	s.m = make(map[T]struct{})

	for _, item := range items {
		s.Put(item)
	}

	return nil
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(lo.Keys(s.m))
}
