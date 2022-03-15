package tiledimage

type Repository interface {
	Add(img *Image)
	Get(id string) (*Image, error)
	Delete(id string) error
}
