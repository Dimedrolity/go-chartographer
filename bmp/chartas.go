package bmp

import (
	"errors"
	"fmt"
	"image"
	"image/color"
)

//Реализация бизнес логики обработки изображений

// Image дополняет image.Image
type Image interface {
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

// SetFragment измененяет пиксели img пикселями fragment, начиная с (x;y) по ширине width и высоте height
// Меняется существующий массив байт img, это производительнее чем создавать абсолютно новое изображение.
// Примечание: если фрагмент перекрывает границы изображения, то часть фрагмента вне изображения игнорируется.
func SetFragment(img Image, x, y, width, height int, fragment image.Image) {
	subRect := image.Rect(x, y, x+width, y+height)
	intersect := subRect.Intersect(img.Bounds())

	for h := 0; h < intersect.Bounds().Dy(); h++ {
		for w := 0; w < intersect.Bounds().Dx(); w++ {
			c := fragment.At(w, h)
			img.Set(x+w, y+h, c)
		}
	}
}

// TODO упростить, сделать параметрами только изображения. У фрагмента должен быть прямоугольник с x,y,x+w,y+h
