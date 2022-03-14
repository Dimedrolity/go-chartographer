// Package store содержит CRUD-функции для разделенных на тайлы изображений.
package store

import (
	"chartographer-go/tile"
	"errors"
	"github.com/google/uuid"
	"image"
	"image/color"
	"sync"
)

// Возможно надо разделить ImageRepo и Tile по pkg.

type ImageRepository interface {
	CreateImage(width, height int) *TiledImage
	GetImage(id string) (*TiledImage, error)
	DeleteImage(id string)
}

// InMemoryImageRepo является потокобезопасным in-memory хранилищем конфигов изображений.
type InMemoryImageRepo struct {
	store map[string]*TiledImage
	mu    *sync.Mutex
}

func New() *InMemoryImageRepo {
	return &InMemoryImageRepo{
		store: make(map[string]*TiledImage),
		mu:    &sync.Mutex{},
	}
}

var ErrNotExist = errors.New("изображение не найдено")

func (r *InMemoryImageRepo) CreateImage(width, height int) *TiledImage {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.createImage(width, height)
}

func (r *InMemoryImageRepo) createImage(width, height int) *TiledImage {
	id := uuid.NewString()
	// TODO использовать зависимость Tiler
	tiles := tile.CreateTiles(width, height, TileMaxSize)

	img := &TiledImage{
		Config: image.Config{
			ColorModel: color.RGBAModel,
			Width:      width,
			Height:     height,
		},
		Id:          id,
		TileMaxSize: TileMaxSize,
		Tiles:       tiles,
	}
	r.store[id] = img

	return r.store[id]
}

func (r *InMemoryImageRepo) GetImage(id string) (*TiledImage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.getImage(id)
}

func (r *InMemoryImageRepo) getImage(id string) (*TiledImage, error) {
	img, ok := r.store[id]
	if !ok {
		return nil, ErrNotExist
	}
	return img, nil
}

func (r *InMemoryImageRepo) DeleteImage(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deleteImage(id)
}

func (r *InMemoryImageRepo) deleteImage(id string) {
	delete(r.store, id)
	return
}
