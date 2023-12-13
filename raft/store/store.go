package store

import (
	"fmt"

	"lxtend.com/m/packages"
)

type InMemoryLogStore struct {
	logs []packages.Entry
}

func (s *InMemoryLogStore) Appends(entry []packages.Entry) {
	s.logs = append(s.logs, entry...)
}

func (s *InMemoryLogStore) Append(entry packages.Entry) {
	s.logs = append(s.logs, entry)
}

func (s *InMemoryLogStore) Get(index int) packages.Entry {
	if index < 0 || index >= len(s.logs) {
		return packages.Entry{Term: 0, Command: ""}
	}
	return s.logs[index]
}

func (s *InMemoryLogStore) GetSince(index int) []packages.Entry {
	if index < 0 {
		return s.logs
	}
	if index >= len(s.logs) {
		return []packages.Entry{}
	}
	return s.logs[index:]
}

func (s *InMemoryLogStore) LastIndex() int {
	// if len(s.logs) == 0 {
	// 	return 0
	// }
	return len(s.logs) - 1
}

func (s *InMemoryLogStore) LastTerm() int {
	if len(s.logs)-1 < 0 {
		return 0
	}
	return s.logs[len(s.logs)-1].Term
}

func (s *InMemoryLogStore) ReplaceFrom(index int, entries []packages.Entry) error {
	if index < 0 || index > len(s.logs) {
		return fmt.Errorf("index out of bounds")
	}
	s.logs = append(s.logs[:index], entries...)
	return nil
}
