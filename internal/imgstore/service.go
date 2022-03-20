package imgstore

import "image"

// Service определяет операции над тайлами в виде image.Image.
type Service interface {
	// SaveTile
	// Координатами являются (x; y), а не img.Bounds().Min.
	SaveTile(id string, x int, y int, img image.Image) error
	// GetTile
	// У возвращаемого image.Image Bounds().Min равен (0; 0).
	// Для смещения на (x; y) использовать RectShifter.
	GetTile(id string, x, y int) (image.Image, error)
	DeleteImage(id string) error

	Encode(img image.Image) ([]byte, error)
	Decode(b []byte) (image.Image, error)
}
