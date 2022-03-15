package tiledimage

type Repository interface {
	CreateImage(width, height int) *Image
	GetImage(id string) (*Image, error)
	DeleteImage(id string) error
}
