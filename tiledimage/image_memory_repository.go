package tiledimage

import (
	"chartographer-go/tile"
	"errors"
	"github.com/google/uuid"
	"image"
	"image/color"
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

func (r *InMemoryImageRepo) CreateImage(width, height int) *Image {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.createImage(width, height)
}

func (r *InMemoryImageRepo) createImage(width, height int) *Image {
	id := uuid.NewString()
	// TODO использовать зависимость Tiler
	// TileMaxSize будет доступен через Tiler.TileMaxSize. Передавать его в CreateTiles не нужно.
	tiles := tile.CreateTiles(width, height, tile.MaxSize)

	// TODO принимать структуру Image единственным параметром

	img := &Image{
		Config: image.Config{
			ColorModel: color.RGBAModel,
			Width:      width,
			Height:     height,
		},
		Id:          id,
		TileMaxSize: tile.MaxSize, // TileMaxSize будет доступен через Tiler.TileMaxSize
		Tiles:       tiles,
	}
	r.store[id] = img

	return r.store[id]
}

func (r *InMemoryImageRepo) GetImage(id string) (*Image, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.getImage(id)
}

func (r *InMemoryImageRepo) getImage(id string) (*Image, error) {
	img, ok := r.store[id]
	if !ok {
		return nil, ErrNotExist
	}
	return img, nil
}

func (r *InMemoryImageRepo) DeleteImage(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.deleteImage(id)
}

func (r *InMemoryImageRepo) deleteImage(id string) error {
	_, ok := r.store[id]
	if !ok {
		return ErrNotExist
	}

	delete(r.store, id)

	return nil
}
