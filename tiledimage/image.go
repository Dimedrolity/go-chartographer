package tiledimage

import "image"

// Image - модель изображения, разделенного на тайлы.
type Image struct {
	Id            string
	Width, Height int
	TileMaxSize   int
	Tiles         []image.Rectangle
}
