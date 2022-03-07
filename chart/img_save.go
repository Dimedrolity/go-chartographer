package chart

import (
	"bytes"
	"image"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"golang.org/x/image/bmp"
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

var dirPath string

func filename(id string) string {
	return filepath.Join(dirPath, id+".bmp")
}

// SetImagesDir - создает директорию, при необходимости, и устанавливает путь к директории
func SetImagesDir(path string) error {
	// If path is already a directory, MkdirAll does nothing and returns nil.
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	dirPath = path

	return nil
}

// GetImage считывает байты изображения и декодирует в image.Image.
// Один из вариантов ошибки - os.ErrNotExist.
func GetImage(id string) (image.Image, error) {
	file, err := os.Open(filename(id))
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

	err = os.WriteFile(filename(id), imgBytes, 0777)
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
	err := os.Remove(filename(id))
	if err != nil {
		return err
	}

	return nil
}
