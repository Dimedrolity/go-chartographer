package imagetile

import (
	"bytes"
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"os"
	"path/filepath"
)

// FileSystemTileRepository - хранилище изображений-тайлов формата BMP в файлах на диске.
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

func (r *FileSystemTileRepository) imgDirPath(id string) string {
	return filepath.Join(r.dirPath, id)
}
func tileFilename(x, y int) string {
	return fmt.Sprintf("Y=%d; X=%d.bmp", y, x)
}
func (r *FileSystemTileRepository) tilePath(id string, x, y int) string {
	return filepath.Join(r.imgDirPath(id), tileFilename(x, y))
}

// GetTile считывает с диска изображение-тайл с координатами (x; y) изображения id и декодирует в формат BMP.
// У возвращаемого image.Image Bounds().Min равен (0; 0).
// Возможны ошибки типа *os.PathError, например os.ErrNotExist.
func (r *FileSystemTileRepository) GetTile(id string, x, y int) (image.Image, error) {
	path := r.tilePath(id, x, y)
	file, err := os.Open(path)
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

	ShiftRect(img, x, y)
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
func (r *FileSystemTileRepository) SaveTile(id string, x int, y int, img image.Image) error {
	// TODO можно было бы обойтись без буфера, Create файл и bmp.Encode(файл)
	// 	file, err := os.OpenFile(filepath.Join(dir, x+".bmp"), os.O_WRONLY|os.O_CREATE, 0777)
	encode, err := Encode(img)
	if err != nil {
		return err
	}

	dir := r.imgDirPath(id)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	path := r.tilePath(id, x, y)
	err = os.WriteFile(path, encode, 0777)
	if err != nil {
		return err
	}

	return nil
}

// DeleteImage удаляет изображение с диска.
func (r *FileSystemTileRepository) DeleteImage(id string) error {
	// If the path does not exist, RemoveAll returns nil (no error).
	err := os.RemoveAll(filepath.Join(r.dirPath, id))
	if err != nil {
		return err
	}
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
