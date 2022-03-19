package imagetile

import "image"

// TileRepository
// наверно не должен работать с image.Image, чтобы в репозитории не было кодирования/декодирования.
type TileRepository interface {
	// SaveTile
	//по новой идее, x y не нужно передавать, будет вычисляться по imb.b.min
	SaveTile(id string, x int, y int, img image.Image) error
	// GetTile
	// У возвращаемого image.Image Bounds().Min равен (x; y).
	// TODO возвращать байты, Репо не должно кодировать и декодировать
	GetTile(id string, x, y int) (image.Image, error)
	DeleteImage(id string) error
}
