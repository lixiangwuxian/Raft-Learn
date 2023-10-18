package main

import (
	"errors"
	"sync"
)

type KVStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (k *KVStore) Get(key string) (string, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if val, ok := k.store[key]; ok {
		return val, nil
	}

	return "", errors.New("key not found")
}

func (k *KVStore) Set(key string, value string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.store[key] = value
}

func (k *KVStore) Delete(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	delete(k.store, key)
}

const (
	Set = iota
	Get
	Delete
	None
)

type Action struct {
	Type  int
	Key   string
	Value string
	Term  int
	Index int
}

type LogStore struct {
	store []Action
}

func NewLogStore() *LogStore {
	return &LogStore{
		store: make([]Action, 0),
	}
}

func (l *LogStore) Append(log Action) {
	l.store = append(l.store, log)
}

func (l *LogStore) Get(index int) Action {
	return l.store[index]
}

func (l *LogStore) Len() int {
	return len(l.store)
}

func (l *LogStore) PopEnd() {
	if len(l.store) == 0 {
		return
	}
	l.store = l.store[:len(l.store)-1]
}
func (l *LogStore) PeekEnd() Action {
	if len(l.store) == 0 {
		return Action{}
	}
	return l.store[len(l.store)-1]
}
func (l *LogStore) PeekLastTerm() int {
	return l.PeekEnd().Term
}
func (l *LogStore) PeekLastIndex() int {
	return l.PeekEnd().Index
}
