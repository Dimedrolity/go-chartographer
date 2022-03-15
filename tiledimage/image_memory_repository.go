package tiledimage

import (
	"errors"
	"sync"
)

// InMemoryImageRepo - потокобезопасное in-memory хранилище конфигов изображений.
type InMemoryImageRepo struct {
	store map[string]*Image
	mu    *sync.Mutex
}

func NewInMemoryImageRepo() *InMemoryImageRepo {
	return &InMemoryImageRepo{
		store: make(map[string]*Image),
		mu:    &sync.Mutex{},
	}
}

var ErrNotExist = errors.New("изображение не найдено")

func (r *InMemoryImageRepo) Add(img *Image) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.add(img)
}

func (r *InMemoryImageRepo) add(img *Image) {
	r.store[img.Id] = img
}

func (r *InMemoryImageRepo) Get(id string) (*Image, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.get(id)
}

func (r *InMemoryImageRepo) get(id string) (*Image, error) {
	img, ok := r.store[id]
	if !ok {
		return nil, ErrNotExist
	}
	return img, nil
}

func (r *InMemoryImageRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.delete(id)
}

func (r *InMemoryImageRepo) delete(id string) error {
	_, ok := r.store[id]
	if !ok {
		return ErrNotExist
	}

	delete(r.store, id)

	return nil
}
