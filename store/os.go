// Package store содержит фукнции записи изображения на диск, получения тайлов
package store

import (
	"bytes"
	"golang.org/x/image/bmp"
	"image"
	"os"
	"path/filepath"
	"strconv"
)

var dirPath string

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

func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GetTile считывает с диска изображение-тайл с координатами (x; y) изображения id и декодирует в image.Image.
// Возможны ошибки типа *os.PathError, например os.ErrNotExist.
func GetTile(id string, x, y int) (image.Image, error) {
	filename := filepath.Join(dirPath, id, strconv.Itoa(y), strconv.Itoa(x)+".bmp")

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

// SaveTile сохраняет тайл-изображение на диск.
// Необходим id и img.Bounds для именования папок и файлов.
// Необходим сам img, чтобы сделать Encode и получить байты.
func SaveTile(id string, img image.Image) error {
	// TODO можно было бы обойтись без буфера, Create файл и bmp.Encode(файл)
	// 	file, err := os.OpenFile(filepath.Join(dir, x+".bmp"), os.O_WRONLY|os.O_CREATE, 0777)
	encode, err := Encode(img)
	if err != nil {
		return err
	}

	y := strconv.Itoa(img.Bounds().Min.Y)
	dir := filepath.Join(dirPath, id, y)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	x := strconv.Itoa(img.Bounds().Min.X)
	err = os.WriteFile(filepath.Join(dir, x+".bmp"), encode, 0777)

	if err != nil {
		return err
	}

	return nil
}

// DeleteImage удаляет изображение с диска.
func DeleteImage(id string) error {
	err := os.RemoveAll(filepath.Join(dirPath, id))
	if err != nil {
		return err
	}

	return nil
}

// SaveImage кодирует изображение в байты и сохраняет на диск.
//func SaveImage(id string, img image.Image) error {
//	imgBytes, err := Encode(img)
//	if err != nil {
//		return err
//	}
//
//	err = os.WriteFile(filename(id), imgBytes, 0777)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// SaveNewImage создает уникальный id, и сохраняет изображение, возвращает id.
//func SaveNewImage(img image.Image) (string, error) {
//	id := uuid.NewString()
//
//	err := SaveImage(id, img)
//	if err != nil {
//		return "", err
//	}
//
//	return id, nil
//}
