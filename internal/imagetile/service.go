package imagetile

import "image"

// Service определяет операции над тайлами в виде image.Image.
type Service interface {
	// SaveTile
	// Координатами являются (x; y), а не img.Bounds().Min.
	SaveTile(id string, x int, y int, img image.Image) error
	// GetTile
	// У возвращаемого image.Image Bounds().Min равен (x; y).
	GetTile(id string, x, y int) (image.Image, error)
	DeleteImage(id string) error
}
