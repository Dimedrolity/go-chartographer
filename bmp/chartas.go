package bmp

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

// TODO w h д.б. целыми числами. В api делать преобразование в int. При ошибке возвращать ошибку, не вызывая NewImage

const (
	minWidth  = 1
	minHeight = 1
	maxWidth  = 20_000
	maxHeight = 50_000
)

// TODO создать свой тип ошибки

func NewImage(width, height int) (image.Image, error) {
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

// TODO w h д.б. целыми числами. В api делать преобразование в int. При ошибке возвращать ошибку, не вызывая Fragment

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
// Примечание: если фрагмент перекрывает границы изображения, то часть фрагмента вне изображения игнорируется.
func SetFragment(img MutableImage, fragment image.Image, x, y, width, height int) {
	// Изображение и фрагмент должны иметь начальные координаты (0;0), иначе функция отработает некорректно.
	start := image.Pt(0, 0)
	if img.Bounds().Min != start || fragment.Bounds().Min != start {
		return
	}

	fragmentRect := image.Rect(x, y, x+width, y+height)
	intersect := img.Bounds().Intersect(fragmentRect)

	for h := 0; h < intersect.Bounds().Dy(); h++ {
		for w := 0; w < intersect.Bounds().Dx(); w++ {
			c := fragment.At(w, h)
			img.Set(x+w, y+h, c)
		}
	}
}
