package store

import (
	"fmt"

	"lxtend.com/m/adapter"
)

type InMemoryLogStore struct {
	logs []adapter.Entry
}

func (s *InMemoryLogStore) Append(entry []adapter.Entry) {
	s.logs = append(s.logs, entry...)
}

func (s *InMemoryLogStore) Get(index int) adapter.Entry {
	if index < 0 || index >= len(s.logs) {
		return adapter.Entry{}
	}
	return s.logs[index]
}

func (s *InMemoryLogStore) GetSince(index int) []adapter.Entry {
	if index < 0 || index >= len(s.logs) {
		return []adapter.Entry{}
	}
	return s.logs[index:]
}

func (s *InMemoryLogStore) LastIndex() int {
	return len(s.logs) - 1
}

func (s *InMemoryLogStore) LastTerm() int {
	if len(s.logs)-1 < 0 {
		return 0
	}
	return s.logs[len(s.logs)-1].Term
}

func (s *InMemoryLogStore) ReplaceFrom(index int, entries []adapter.Entry) error {
	if index < 0 || index > len(s.logs) {
		return fmt.Errorf("index out of bounds")
	}
	s.logs = append(s.logs[:index], entries...)
	return nil
}
