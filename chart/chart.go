// Package chart - слой сервиса приложения.
package chart

import (
	"chartographer-go/store"
	"chartographer-go/tile"
	"chartographer-go/tiledimage"
	"errors"
	"github.com/google/uuid"
	"image"
	"image/color"
	"image/draw"
)

type Service interface {
	NewRgbaBmp(width, height int) (*tiledimage.Image, error)
	DeleteImage(id string) error
	SetFragment(tiledImageId string, fragment image.Image) error
	GetFragment(imgConfig *tiledimage.Image, x, y, width, height int) (image.Image, error)
	GetTiledImage(id string) (*tiledimage.Image, error)
}

// ChartographerService - содержит бизнес логику обработки изображений.
type ChartographerService struct {
	imageRepo   tiledimage.Repository
	tileRepo    store.TileRepository
	tileMaxSize int // Определяет максимальный размер тайла по ширине и высоте.
}

func NewChartographerService(imageRepo tiledimage.Repository, tileRepo store.TileRepository, tileMaxSize int) *ChartographerService {
	return &ChartographerService{imageRepo: imageRepo, tileRepo: tileRepo, tileMaxSize: tileMaxSize}
}

const (
	minWidth  = 1
	minHeight = 1
	maxWidth  = 20_000
	maxHeight = 50_000
)

// NewRgbaBmp разделяет размеры изображения на тайлы, создает image.RGBA изображения в соответствии с тайлами,
// записывает изображения на диск в формате BMP.
// Возможна ошибка типа *SizeError
func (cs *ChartographerService) NewRgbaBmp(width, height int) (*tiledimage.Image, error) {
	if width < minWidth || width > maxWidth ||
		height < minHeight || height > maxHeight {
		return nil, &SizeError{
			minWidth: minWidth, width: width, maxWidth: maxWidth,
			minHeight: minHeight, height: height, maxHeight: maxHeight,
		}
	}

	tiles := tile.CreateTiles(width, height, cs.tileMaxSize)

	img := &tiledimage.Image{
		Id:          uuid.NewString(),
		Width:       width,
		Height:      height,
		TileMaxSize: cs.tileMaxSize,
		Tiles:       tiles,
	}
	cs.imageRepo.Add(img)

	for _, t := range img.Tiles {
		i := newOpaqueRGBA(t)

		err := cs.tileRepo.SaveTile(img.Id, t.Min.X, t.Min.Y, i)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// newOpaqueRGBA создает image.RGBA и устанавливает alpha-канал максимальным значением.
// Таким образом, изображение в дальнейшем будет кодироваться без учета альфа канала (24-бит на пиксель).
func newOpaqueRGBA(r image.Rectangle) image.Image {
	img := image.NewRGBA(r)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0xFF})
		}
	}

	return img
}

func (cs *ChartographerService) DeleteImage(id string) error {
	err := cs.imageRepo.Delete(id)
	if err != nil {
		return err
	}

	err = cs.tileRepo.DeleteImage(id)
	if err != nil {
		return err
	}

	return nil
}

// SetFragment измененяет пиксели изображения id пикселями фрагмента fragment,
// накладывая прямогольник фрагмента - fragment.Bounds() - на изображение.
// Изображение имеет начальные координаты (0;0), фрагмент может иметь начальные координаты отличные от (0;0).
//
// Меняется существующий массив байт изображения, это производительнее чем создавать абсолютно новое изображение.
//
// Примечание:
// если фрагмент частично выходит за границы изображения, то часть фрагмента вне изображения игнорируется.
func (cs *ChartographerService) SetFragment(tiledImageId string, fragment image.Image) error {
	img, err := cs.imageRepo.Get(tiledImageId)
	if err != nil {
		return err
	}

	imgRect := image.Rect(0, 0, img.Width, img.Height)

	if !imgRect.Overlaps(fragment.Bounds()) {
		return ErrNotOverlaps
	}

	intersect := fragment.Bounds().Intersect(imgRect)

	overlapped := tile.FilterOverlappedTiles(img.Tiles, intersect)

	for _, t := range overlapped {
		tileImg, err := cs.tileRepo.GetTile(img.Id, t.Min.X, t.Min.Y)
		if err != nil {
			return err
		}

		mutableImage, ok := tileImg.(draw.Image)
		if !ok {
			return errors.New("изображение должно реализовывать draw.Image")
		}

		// TODO refactor&comment,
		tileIntersect := t.Bounds().Intersect(intersect)

		// sp не фрагмент, а пересечение фрагмента с t.
		// Нельзя передавать всегда один и тот же sp.
		fragIntersect := t.Bounds().Intersect(fragment.Bounds())
		draw.Draw(mutableImage, tileIntersect, fragment, fragIntersect.Bounds().Min, draw.Src)

		err = cs.tileRepo.SaveTile(img.Id, t.Min.X, t.Min.Y, mutableImage)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	fragmentMinWidth  = 1
	fragmentMinHeight = 1
	fragmentMaxWidth  = 5_000
	fragmentMaxHeight = 5_000
)

// GetFragment возвращает фрагмент изображения id, начиная с координат изобржаения (x; y) по ширине width и высоте height.
// Возвращаемое изображение будет иметь начальные координаты (x; y).
// Примечание: часть фрагмента вне границ изображения будет иметь чёрный цвет (цвет по умолчанию).
// Возможны ошибки SizeError, ErrNotOverlaps и типа *os.PathError, например os.ErrNotExist.
func (cs *ChartographerService) GetFragment(imgConfig *tiledimage.Image, x, y, width, height int) (image.Image, error) {
	if width < fragmentMinWidth || width > fragmentMaxWidth ||
		height < fragmentMinHeight || height > fragmentMaxHeight {
		return nil, &SizeError{
			minWidth: fragmentMinWidth, width: width, maxWidth: fragmentMaxWidth,
			minHeight: fragmentMinHeight, height: height, maxHeight: fragmentMaxHeight,
		}
	}

	imgRect := image.Rect(0, 0, imgConfig.Width, imgConfig.Height)
	fragmentRect := image.Rect(x, y, x+width, y+height)
	if !imgRect.Overlaps(fragmentRect) {
		return nil, ErrNotOverlaps
	}

	overlapped := tile.FilterOverlappedTiles(imgConfig.Tiles, fragmentRect)

	fragment := image.NewRGBA(fragmentRect)

	for _, t := range overlapped {
		tileImg, err := cs.tileRepo.GetTile(imgConfig.Id, t.Min.X, t.Min.Y)
		if err != nil {
			return nil, err
		}

		intersect := t.Intersect(fragmentRect)

		draw.Draw(fragment, intersect, tileImg, intersect.Min, draw.Src)
	}

	return fragment, nil
}

func (cs *ChartographerService) GetTiledImage(id string) (*tiledimage.Image, error) {
	return cs.imageRepo.Get(id)
}
