// Package chart является уровнем сервиса, содержит бизнес логику обработки изображений.
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

// ImageRepo
//TODO не использовать глобальные переменную. Выделить структуру сервиса и принимать репо в конструкторе (NewInMemoryImageRepo)
var ImageRepo tiledimage.Repository
var TileRepo store.TileRepository

type MutableImage interface {
	image.Image
	Set(x, y int, c color.Color)
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
func NewRgbaBmp(width, height int) (*tiledimage.Image, error) {
	if width < minWidth || width > maxWidth ||
		height < minHeight || height > maxHeight {
		return nil, &SizeError{
			minWidth: minWidth, width: width, maxWidth: maxWidth,
			minHeight: minHeight, height: height, maxHeight: maxHeight,
		}
	}

	tiles := tile.CreateTiles(width, height, tile.MaxSize)

	img := &tiledimage.Image{
		Config: image.Config{
			ColorModel: color.RGBAModel,
			Width:      width,
			Height:     height,
		},
		Id:          uuid.NewString(),
		TileMaxSize: tile.MaxSize,
		Tiles:       tiles,
	}
	ImageRepo.Add(img)

	for _, t := range img.Tiles {
		i := NewOpaqueRGBA(t)

		err := TileRepo.SaveTile(img.Id, i)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// NewOpaqueRGBA создает image.RGBA и устанавливает alpha-канал максимальным значением.
// Таким образом, изображение в дальнейшем будет кодироваться без учета альфа канала (24-бит на пиксель).
func NewOpaqueRGBA(r image.Rectangle) image.Image {
	img := image.NewRGBA(r)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0xFF})
		}
	}

	return img
}

func DeleteImage(id string) error {
	err := ImageRepo.Delete(id)
	if err != nil {
		return err
	}

	err = TileRepo.DeleteImage(id)
	if err != nil {
		return err
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
// Примечание: часть фрагмента вне границ изображения будет иметь чёрный цвет (цвет по умолчанию).
// Возможны ошибки SizeError, ErrNotOverlaps и типа *os.PathError, например os.ErrNotExist.
func GetFragment(imgConfig *tiledimage.Image, x, y, width, height int) (image.Image, error) {
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

	img := image.NewRGBA(fragmentRect)

	for _, t := range overlapped {
		tileImg, err := TileRepo.GetTile(imgConfig.Id, t.Min.X, t.Min.Y)
		if err != nil {
			return nil, err
		}

		tile.DrawIntersection(img, tileImg, fragmentRect)
	}

	return img, nil
}

// SetFragment измененяет пиксели изображения id пикселями фрагмента fragment, начиная с координат изобржаения (x;y) по ширине width и высоте height.
// Меняется существующий массив байт изображения, это производительнее чем создавать абсолютно новое изображение.
// Примечание:
// если фрагмент частично выходит за границы изображения, то часть фрагмента вне изображения игнорируется.
func SetFragment(tiledImageId string, fragment image.Image, x, y, width, height int) error {
	img, err := ImageRepo.Get(tiledImageId)
	if err != nil {
		return err
	}

	imgRect := image.Rect(0, 0, img.Width, img.Height)
	fragmentRect := image.Rect(x, y, x+width, y+height)
	if !imgRect.Overlaps(fragmentRect) {
		return ErrNotOverlaps
	}

	tiles := img.Tiles

	overlapped := tile.FilterOverlappedTiles(tiles, fragmentRect)

	for _, t := range overlapped {
		tileImg, err := TileRepo.GetTile(img.Id, t.Min.X, t.Min.Y)
		if err != nil {
			return err
		}

		mutableImage, ok := tileImg.(draw.Image)
		if !ok {
			return errors.New("изображение должно реализовывать draw.Image")
		}

		intersect := mutableImage.Bounds().Intersect(fragmentRect)
		draw.Draw(mutableImage, intersect, fragment, fragment.Bounds().Min, draw.Src)

		err = TileRepo.SaveTile(tiledImageId, mutableImage)
		if err != nil {
			return err
		}
	}

	return nil
}
