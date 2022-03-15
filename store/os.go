package store

import (
	"bytes"
	"golang.org/x/image/bmp"
	"image"
	"os"
	"path/filepath"
	"strconv"
)

// TileRepository
// наверно не должен работать с image.Image, чтобы в репозитории не было кодирования/декодирования.
type TileRepository interface {
	SaveTile(id string, img image.Image) error
	GetTile(id string, x, y int) (image.Image, error)
	DeleteImage(id string) error
}

// FileSystemTileRepository - хранит изображения-тайлы в файлах на диске.
type FileSystemTileRepository struct {
	dirPath string
}

func NewFileSystemTileRepo(dirPath string) (*FileSystemTileRepository, error) {
	// If path is already a directory, MkdirAll does nothing and returns nil.
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	repo := &FileSystemTileRepository{
		dirPath: dirPath,
	}
	return repo, nil
}

var TileRepo TileRepository

// GetTile считывает с диска изображение-тайл с координатами (x; y) изображения id и декодирует в формат BMP.
// Возможны ошибки типа *os.PathError, например os.ErrNotExist.
func (r *FileSystemTileRepository) GetTile(id string, x, y int) (image.Image, error) {
	filename := filepath.Join(r.dirPath, id, strconv.Itoa(y), strconv.Itoa(x)+".bmp")

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
func (r *FileSystemTileRepository) SaveTile(id string, img image.Image) error {
	// TODO можно было бы обойтись без буфера, Create файл и bmp.Encode(файл)
	// 	file, err := os.OpenFile(filepath.Join(dir, x+".bmp"), os.O_WRONLY|os.O_CREATE, 0777)
	encode, err := Encode(img)
	if err != nil {
		return err
	}

	y := strconv.Itoa(img.Bounds().Min.Y)
	dir := filepath.Join(r.dirPath, id, y)
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

func (r *FileSystemTileRepository) DeleteImage(id string) error {
	err := os.RemoveAll(filepath.Join(r.dirPath, id))
	if err != nil {
		return err
	}
	return nil
}

// DeleteImage
//TODO выделить в слой сервиса
func DeleteImage(id string) error {
	err := ImageRepo.DeleteImage(id)
	if err != nil {
		return err
	}

	err = TileRepo.DeleteImage(id)
	if err != nil {
		return err
	}

	return nil
}

////

// TileMaxSize определяет максимальный размер тайла по ширине и высоте.
// Необходимо проинициализировать перед использованием функций текущего pkg
// TODO выделить в структуру Tiler, и фукнцию NewTiler(tileMaxSize). Тогда сделать все фукнции методами Tiler
// Это не нужно хранить в TiledImgRepo, так как репо предназначен для хранения самой модели TiledImage.
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
// Выделить сущность FileSystemTileRepository, она будет принимать в конструкторе ImRepo. Нет, репозитории не должны зависеть друг от друга.
var ImageRepo TiledImageRepository
