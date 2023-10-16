package raft

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
