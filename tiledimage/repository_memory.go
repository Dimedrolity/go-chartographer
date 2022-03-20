package tiledimage

import (
	"errors"
	"sync"
)

// InMemoryRepo - потокобезопасное in-memory хранилище.
type InMemoryRepo struct {
	store map[string]interface{}
	mu    *sync.Mutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		store: make(map[string]interface{}),
		mu:    &sync.Mutex{},
	}
}

var ErrNotExist = errors.New("не найдено")

func (r *InMemoryRepo) Add(key string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.add(key, value)
}

func (r *InMemoryRepo) add(key string, value interface{}) {
	r.store[key] = value
}

func (r *InMemoryRepo) Get(key string) (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.get(key)
}

func (r *InMemoryRepo) get(key string) (interface{}, error) {
	value, ok := r.store[key]
	if !ok {
		return nil, ErrNotExist
	}
	return value, nil
}

func (r *InMemoryRepo) Delete(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.delete(key)
}

func (r *InMemoryRepo) delete(key string) error {
	_, ok := r.store[key]
	if !ok {
		return ErrNotExist
	}

	delete(r.store, key)

	return nil
}
