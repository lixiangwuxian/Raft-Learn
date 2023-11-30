package store

import (
	"fmt"

	"lxtend.com/m/adapter"
)

type InMemoryLogStore struct {
	logs []adapter.Entry
}

func (s *InMemoryLogStore) Append(entry adapter.Entry) error {
	s.logs = append(s.logs, entry)
	return nil
}

func (s *InMemoryLogStore) Get(index int) (adapter.Entry, error) {
	if index < 0 || index >= len(s.logs) {
		return adapter.Entry{}, fmt.Errorf("index out of bounds")
	}
	return s.logs[index], nil
}

func (s *InMemoryLogStore) LastIndex() (int, error) {
	return len(s.logs) - 1, nil
}

func (s *InMemoryLogStore) ReplaceFrom(index int, entries []adapter.Entry) error {
	if index < 0 || index > len(s.logs) {
		return fmt.Errorf("index out of bounds")
	}
	s.logs = append(s.logs[:index], entries...)
	return nil
}
