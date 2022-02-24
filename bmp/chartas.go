package bmp

import (
	"errors"
	"fmt"
	"image"
)

//Реализация бизнес логики обработки изображений

// TODO w h д.б. целыми числами. В api делать преобразование в int. При ошибке возвращать ошибку, не вызывая CreateImg

const (
	minWidth  = 1
	minHeight = 1
	maxWidth  = 20_000
	maxHeight = 50_000
)

func CreateImg(width int, height int) (*image.RGBA, error) {
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

//var pathToFolder string
//
//func SetPath(path string) {
//	pathToFolder = path
//}

// или отделить логику именования от записи на диск?
// должен принимать путь до каталога, в котором сервис может хранить данные.
//func SaveImg(img *image.RGBA) (string, error) {
//	w := bytes.Buffer{}
//	err := bmp.Encode(&w, img)
//	if err != nil {
//		return "", err
//	}
//
//	id := guid()
//	name := id + ".bmp"
//
//	err = os.WriteFile(name, w.Bytes(), 0777)
//	if err != nil {
//		return "", err
//	}
//
//	return id, nil
//}
//
//func guid() string {
//	uuidWithHyphen := uuid.New()
//	return uuidWithHyphen.String()
//}
