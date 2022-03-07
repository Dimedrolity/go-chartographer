package chart

import (
	"errors"
	"fmt"
	"image"
	"image/color"
)

//Реализация бизнес логики обработки изображений

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

// TODO создать свой тип ошибки?

func NewRGBA(width, height int) (image.Image, error) {
	if width < minWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Ширина изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная ширина=%d", width))
	}
	if height < minHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Высота изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная высота=%d", height))
	}

	if width > maxWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая ширина = %d.\n", maxWidth) +
			fmt.Sprintf("Полученная ширина=%d", width))
	}
	if height > maxHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая высота = %d.\n", maxHeight) +
			fmt.Sprintf("Полученная высота=%d", height))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return img, nil
}

const (
	fragmentMinWidth  = 1
	fragmentMinHeight = 1
	fragmentMaxWidth  = 5_000
	fragmentMaxHeight = 5_000
)

// Fragment возвращает фрагмент изображения img, начиная с координат изобржаения (x;y) по ширине width и высоте height.
// Примечание: часть фрагмента вне границ изображения будет иметь чёрный цвет (цвет по умолчанию).
func Fragment(img image.Image, x, y, width, height int) (image.Image, error) {
	if width < fragmentMinWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Ширина изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height < fragmentMinHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Высота изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	if width > fragmentMaxWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая ширина = %d.\n", fragmentMaxWidth) +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height > fragmentMaxHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая высота = %d.\n", fragmentMaxHeight) +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	fragment := image.NewRGBA(image.Rect(x, y, x+width, y+height))
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
// 1. если фрагмент перекрывает границы изображения, то часть фрагмента вне изображения игнорируется.
// 2. изображение и фрагмент должны иметь начальные координаты (0;0).
func SetFragment(img image.Image, fragment image.Image, x, y, width, height int) error {
	mutableImage, ok := img.(MutableImage)
	if !ok {
		return errors.New("ошибка. Изображение должно реализовывать MutableImage")
	}

	start := image.Pt(0, 0)
	if img.Bounds().Min != start || fragment.Bounds().Min != start {
		return errors.New("ошибка. Изображение и фрагмент должны иметь начальные координаты (0;0). " +
			fmt.Sprintf("изображение имеет %v, ", img.Bounds().Min) +
			fmt.Sprintf("фрагмент имеет %v.", fragment.Bounds().Min))
	}

	fragmentRect := image.Rect(x, y, x+width, y+height)
	intersect := img.Bounds().Intersect(fragmentRect)

	for h := 0; h < intersect.Bounds().Dy(); h++ {
		for w := 0; w < intersect.Bounds().Dx(); w++ {
			c := fragment.At(w, h)
			mutableImage.Set(x+w, y+h, c)
		}
	}

	return nil
}
