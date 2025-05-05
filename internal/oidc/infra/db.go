package infra

import (
	"errors"
	"sync"
)

type DB struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

func NewDB() *DB {
	return &DB{
		store: make(map[string]interface{}),
	}
}

func (db *DB) Set(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.store[key] = value
}

func (db *DB) Get(key string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	value, exists := db.store[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (db *DB) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.store, key)
}
