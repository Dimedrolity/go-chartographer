// Package chart содержит бизнес логику обработки изображений
package chart

import (
	"chartographer-go/store"
	"chartographer-go/tile"
	"errors"
	"fmt"
	"image"
	"image/color"
)

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
func NewRgbaBmp(width, height int) (*store.TiledImage, error) {
	if width < minWidth || width > maxWidth ||
		height < minHeight || height > maxHeight {
		return nil, &SizeError{
			minWidth: minWidth, width: width, maxWidth: maxWidth,
			minHeight: minHeight, height: height, maxHeight: maxHeight,
		}
	}

	img, err := store.CreateImage(width, height)
	if err != nil {
		return nil, err
	}

	return img, nil
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
func GetFragment(imgConfig *store.TiledImage, x, y, width, height int) (image.Image, error) {
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
		tileImg, err := store.TileRepo.GetTile(imgConfig.Id, t.Min.X, t.Min.Y)
		if err != nil {
			return nil, err
		}

		tile.DrawIntersection(img, tileImg, fragmentRect)
	}

	return img, nil
}

// Fragment возвращает фрагмент изображения img, начиная с координат изобржаения (x;y) по ширине width и высоте height.
// Примечание: часть фрагмента вне границ изображения будет иметь чёрный цвет (цвет по умолчанию).
// Возможны ошибки ErrNotOverlaps или типа *SizeError
func Fragment(img image.Image, x, y, width, height int) (image.Image, error) {
	if width < fragmentMinWidth || width > fragmentMaxWidth ||
		height < fragmentMinHeight || height > fragmentMaxHeight {
		return nil, &SizeError{
			minWidth: fragmentMinWidth, width: width, maxWidth: fragmentMaxWidth,
			minHeight: fragmentMinHeight, height: height, maxHeight: fragmentMaxHeight,
		}
	}

	fragment := image.NewRGBA(image.Rect(x, y, x+width, y+height))
	if !img.Bounds().Overlaps(fragment.Bounds()) {
		return nil, ErrNotOverlaps
	}
	intersect := img.Bounds().Intersect(fragment.Bounds())

	for h := intersect.Min.Y; h < intersect.Max.Y; h++ {
		for w := intersect.Min.X; w < intersect.Max.X; w++ {
			c := img.At(w, h)
			fragment.Set(w, h, c)
		}
	}

	return fragment, nil
}

// SetFragment измененяет пиксели изображения img пикселями фрагмента fragment, начиная с координат изобржаения (x;y) по ширине width и высоте height.
// Меняется существующий массив байт изображения, это производительнее чем создавать абсолютно новое изображение.
// Примечания:
// 1. если фрагмент частично выходит за границы изображения, то часть фрагмента вне изображения игнорируется.
// 2. изображение и фрагмент должны иметь начальные координаты (0;0).
func SetFragment(img image.Image, fragment image.Image, x, y, width, height int) error {
	mutableImage, ok := img.(MutableImage)
	if !ok {
		return errors.New("изображение должно реализовывать MutableImage")
	}

	start := image.Pt(0, 0)
	if img.Bounds().Min != start || fragment.Bounds().Min != start {
		return errors.New("изображение и фрагмент должны иметь начальные координаты (0;0). " +
			fmt.Sprintf("изображение имеет %v, ", img.Bounds().Min) +
			fmt.Sprintf("фрагмент имеет %v.", fragment.Bounds().Min))
	}

	fragmentRect := image.Rect(x, y, x+width, y+height)
	if !img.Bounds().Overlaps(fragmentRect) {
		return ErrNotOverlaps
	}

	intersect := img.Bounds().Intersect(fragmentRect)

	for h := 0; h < intersect.Bounds().Dy(); h++ {
		for w := 0; w < intersect.Bounds().Dx(); w++ {
			c := fragment.At(w, h)
			mutableImage.Set(x+w, y+h, c)
		}
	}

	return nil
}
