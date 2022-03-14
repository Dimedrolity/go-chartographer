package store

import (
	"sync"
)

// наверно надо разделять ImageStore и Repository по pkg.
//И только Repo pkg будет использовать Store pkg

// ImageStore является потокобезопасным in-memory хранилищем конфигов изображений.
type ImageStore struct {
	mu    sync.RWMutex
	store map[string]*Image
}

func New() *ImageStore {
	return &ImageStore{
		store: make(map[string]*Image),
	}
}

func (s *ImageStore) set(id string, config *Image) {
	s.store[id] = config
}

func (s *ImageStore) get(id string) (config *Image, ok bool) {
	config, ok = s.store[id]
	return
}

func (s *ImageStore) Set(id string, config *Image) {
	s.mu.Lock()
	s.set(id, config)
	s.mu.Unlock()
}

func (s *ImageStore) Get(id string) (config *Image, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.get(id)
}
