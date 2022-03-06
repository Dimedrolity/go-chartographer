package chart

import (
	"bytes"
	"github.com/google/uuid"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// TODO может в отдельный pkg?

// TODO директория должна указываться при инициализации приложения
// TODO должна создаваться, если она не существует.
var pathToFolder = "data/"

// GetImage считывает байты изображения и декодирует в image.Image.
// Один из вариантов ошибки - os.ErrNotExist.
func GetImage(id string) (image.Image, error) {
	filename := pathToFolder + id + ".bmp"

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img, err := bmp.Decode(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveImage кодирует изображение в байты и сохраняет на диск.
func SaveImage(id string, img image.Image) error {
	imgBytes, err := Encode(img)
	if err != nil {
		return err
	}

	filename := pathToFolder + id + ".bmp"

	err = os.WriteFile(filename, imgBytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

// SaveNewImage создает уникальный id, и сохраняет изображение, возвращает id.
func SaveNewImage(img image.Image) (string, error) {
	id := uuid.NewString()

	err := SaveImage(id, img)
	if err != nil {
		return "", err
	}

	return id, nil
}

// DeleteImage удаляет файл изображения.
func DeleteImage(id string) error {
	filename := pathToFolder + id + ".bmp"

	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}
