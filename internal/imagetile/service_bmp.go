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

	img, err := bmp.Decode(bytes.NewReader(tile))
	if err != nil {
		return nil, err
	}

	ShiftRect(img, x, y)

	return img, nil
}

// SaveTile декодирует тайл-изображение в формат BMP и сохраняет.
func (s *BmpService) SaveTile(id string, x int, y int, img image.Image) error {
	// TODO можно было бы обойтись без буфера, Create файл и bmp.Encode(файл)
	// 	file, err := os.OpenFile(filepath.Join(dir, x+".bmp"), os.O_WRONLY|os.O_CREATE, 0777)
	encode, err := Encode(img)
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

// Encode обертка над bmp. TODO нужна ли?
func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
