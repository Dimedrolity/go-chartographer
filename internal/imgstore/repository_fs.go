package imgstore

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystemTileRepository - хранилище изображений-тайлов в файлах на диске.
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

// SaveTile сохраняет тайл-изображение на диск.
// По id создается папка на диске для тайлов изображения, для каждого тайла создается файл и именуется по координатам "Y=<y>; X=<x>.bmp".
func (r *FileSystemTileRepository) SaveTile(id string, x int, y int, img []byte) error {
	dir := r.imgDirPath(id)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	path := r.tilePath(id, x, y)
	err = os.WriteFile(path, img, 0777)
	if err != nil {
		return err
	}

	return nil
}

// GetTile считывает с диска изображение-тайл с координатами (x; y) изображения id.
// Возможны ошибки типа *os.PathError, например os.ErrNotExist.
func (r *FileSystemTileRepository) GetTile(id string, x, y int) ([]byte, error) {
	path := r.tilePath(id, x, y)
	return os.ReadFile(path)
}

// DeleteImage удаляет изображение с диска.
func (r *FileSystemTileRepository) DeleteImage(id string) error {
	// If the path does not exist, RemoveAll returns nil (no error).
	return os.RemoveAll(filepath.Join(r.dirPath, id))
}
