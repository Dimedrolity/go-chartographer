package tiledimage

import "image"

// Image - модель изображения, разделенного на тайлы.
type Image struct {
	image.Config
	Id          string
	TileMaxSize int
	Tiles       []image.Rectangle
}
