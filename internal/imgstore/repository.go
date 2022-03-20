package imgstore

type Repository interface {
	SaveTile(id string, x int, y int, img []byte) error
	GetTile(id string, x, y int) ([]byte, error)
	DeleteImage(id string) error
}
