package imagetile

import (
	"bytes"
	"golang.org/x/image/bmp"
	"image"
)

// BmpService - хранилище изображений-тайлов формата BMP.
type BmpService struct {
	repo Repository
}

func NewBmpService(r Repository) *BmpService {
	repo := &BmpService{
		repo: r,
	}
	return repo
}

// GetTile возвращает изображение-тайл с координатами (x; y) изображения id в формате BMP.
// У возвращаемого image.Image Bounds().Min равен (x; y).
func (s *BmpService) GetTile(id string, x, y int) (image.Image, error) {
	tile, err := s.repo.GetTile(id, x, y)
	if err != nil {
		return nil, err
	}

	img, err := s.Decode(tile)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveTile декодирует тайл-изображение в формат BMP и сохраняет.
func (s *BmpService) SaveTile(id string, x int, y int, img image.Image) error {
	encode, err := s.Encode(img)
	if err != nil {
		return err
	}

	err = s.repo.SaveTile(id, x, y, encode)
	if err != nil {
		return err
	}

	return nil
}

// DeleteImage удаляет изображение.
func (s *BmpService) DeleteImage(id string) error {
	return s.repo.DeleteImage(id)
}

// Encode декодирует image.Image в формат BMP.
func (s *BmpService) Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode декодирует байты BMP изображения в image.Image.
func (s *BmpService) Decode(b []byte) (image.Image, error) {
	img, err := bmp.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return img, nil
}
