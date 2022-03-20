package kvstore

import (
	"sync"
)

// InMemoryStore - потокобезопасное in-memory хранилище.
type InMemoryStore struct {
	store map[string]interface{}
	mu    *sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]interface{}),
		mu:    &sync.Mutex{},
	}
}

func (r *InMemoryStore) Add(key string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.add(key, value)
}

func (r *InMemoryStore) add(key string, value interface{}) {
	r.store[key] = value
}

func (r *InMemoryStore) Get(key string) (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.get(key)
}

func (r *InMemoryStore) get(key string) (interface{}, error) {
	value, ok := r.store[key]
	if !ok {
		return nil, ErrNotExist
	}
	return value, nil
}

func (r *InMemoryStore) Delete(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.delete(key)
}

func (r *InMemoryStore) delete(key string) error {
	_, ok := r.store[key]
	if !ok {
		return ErrNotExist
	}

	delete(r.store, key)

	return nil
}
