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

func NewImage(width, height int) (image.Image, error) {
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

// TODO w h д.б. целыми числами. В api делать преобразование в int. При ошибке возвращать ошибку, не вызывая SubImage

const (
	subImageMinWidth  = 1
	subImageMinHeight = 1
	subImageMaxWidth  = 5_000
	subImageMaxHeight = 5_000
)

func SubImage(img image.Image, x, y, width, height int) (image.Image, error) {
	if x < subImageMinWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Ширина изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if y < subImageMinHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Высота изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	if width < subImageMinWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Ширина изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height < subImageMinHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Высота изображения должна быть положительным числом.\n") +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	if width > subImageMaxWidth {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая ширина = %d.\n", subImageMaxWidth) +
			fmt.Sprintf("Полученная ширина имеет значение %d", width))
	}
	if height > subImageMaxHeight {
		return nil, errors.New(fmt.Sprintf("ошибка. Изображение превышает допустимый размер, максимально допустимая высота = %d.\n", subImageMaxHeight) +
			fmt.Sprintf("Полученная высота имеет значение %d", height))
	}

	// проверка что x, y содержится в img.Bounds (метод In)
	//if !(Point{x, y}.In(p.Rect)) {
	//	return color.RGBA64{}
	//}

	// проверка что w, h не превышают

	// rgba.SubImage не подходит, так как работает через Intersect. По требованию нужен??
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return nil, errors.New("ошибка. Тип изображения не поддерживается")
	}
	sub := rgba.SubImage(image.Rect(x, y, x+width, y+height))

	// будет ли черный по умолч?
	// БУДЕТ ЛИ ВЕРНО Encode?

	// Создание Rect с координат x y
	//r := image.Rect(x, y, x+width, y+height)
	//sub := image.NewRGBA(r)
	//
	//for h := 0; h < height; h++ {
	//	for w := 0; w < width; w++ {
	//		c := img.At(x+w, y+h)
	//		sub.Set(x+w, y+h, c)
	//	}
	//}

	// Создание Rect с координат 0 0
	//sub := image.NewRGBA(image.Rect(0, 0, width, height))
	//xSub := 0
	//ySub := 0
	//for yImg := 0; yImg < height; yImg++ {
	//	for xImg := 0; xImg < width; xImg++ {
	//		c := img.At(x+xImg, y+yImg)
	//		sub.Set(xSub, ySub, c)
	//		xSub++
	//	}
	//	ySub++
	//}

	return sub, nil
}
