package hashset

import (
	"encoding/json"

	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type Set[T comparable] struct {
	m map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) PutAll(values []T) {
	for _, value := range values {
		s.Put(value)
	}
}

func (s *Set[T]) Values() []T {
	return lo.Keys(s.m)
}

func (s *Set[T]) Clear() {
	if s.m == nil {
		return
	}

	maps.Clear(s.m)
}

func (s *Set[T]) Put(value T) {
	if s.m == nil {
		return
	}

	s.m[value] = struct{}{}
}

func (s *Set[T]) Has(value T) bool {
	if s.m == nil {
		return false
	}

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
