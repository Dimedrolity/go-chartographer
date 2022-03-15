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

// TileMaxSize определяет максимальный размер тайла по ширине и высоте.
// Необходимо проинициализировать перед использованием функций текущего pkg
// TODO выделить в структуру TileRepository, и фукнцию NewTileRepo(tileMaxSize). Тогда сделать все фукнции методами Repo
var TileMaxSize int

func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ImageRepo TODO не использовать глобальную переменную.
// Выделить сущность TileRepo, она будет принимать в конструкторе ImRepo
var ImageRepo TiledImageRepository

// GetTile считывает с диска изображение-тайл с координатами (x; y) изображения id и декодирует в формат BMP.
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

// SaveTile декодирует тайл-изображение в формат BMP и сохраняет на диск.
// По id создается папка на диске для тайлов изображения,
// каждый тайл хранится в папке начальной координаты Y, файл именуется координатой X.
//
// Пример структуры файлов для изображения с id="3a8cc52-8997-4adb-8a09-918c29aa10c4" и координатами тайла (0; 10):
//
// 23a8cc52-8997-4adb-8a09-918c29aa10c4
// 	+---- 10
//		+---- 0.bmp
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
	err := ImageRepo.DeleteImage(id)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(dirPath, id))
	if err != nil {
		return err
	}

	return nil
}

// CreateImage создает RGBA изображение формата BMP, возвращает модель изображения.
// Если ширина/высота изображения больше TileMaxSize, то изображение разделяется на тайлы.
func CreateImage(width, height int) (*TiledImage, error) {
	img := ImageRepo.CreateImage(width, height)

	for _, t := range img.Tiles {
		i := image.NewRGBA(t)

		// TODO использовать зависимость Tiler
		err := SaveTile(img.Id, i)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}
