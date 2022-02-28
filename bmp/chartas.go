package bmp

import (
	"errors"
	"fmt"
	"image"
	"image/color"
)

//Реализация бизнес логики обработки изображений

// MutableImage дополняет image.Image
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
