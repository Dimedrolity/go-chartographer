package chart

import "image"

// Service определяет бизнес логику обработки изображений.
type Service interface {
	AddImage(width, height int) (*TiledImage, error)
	GetImage(id string) (*TiledImage, error)
	DeleteImage(id string) error

	SetFragment(img *TiledImage, x int, y int, fragment image.Image) error
	GetFragment(img *TiledImage, x, y, width, height int) (image.Image, error)

	Encode(img image.Image) ([]byte, error)
	Decode(b []byte) (image.Image, error)
}
