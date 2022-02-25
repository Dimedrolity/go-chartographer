package bmp

import (
	"errors"
	"fmt"
	"image"
)

//Реализация бизнес логики обработки изображений

// TODO w h д.б. целыми числами. В api делать преобразование в int. При ошибке возвращать ошибку, не вызывая NewImage

const (
	minWidth  = 1
	minHeight = 1
	maxWidth  = 20_000
	maxHeight = 50_000
)

func NewImage(width int, height int) (image.Image, error) {
	if width < minWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Ширина изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height < minHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Высота изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	if width > maxWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая ширина = %d.\n", maxWidth) +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height > maxHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая высота = %d.\n", maxHeight) +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return img, nil
}
