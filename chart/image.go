package chart

import "image"

// TiledImage - изображение, разделенное на тайлы.
type TiledImage struct {
	Id            string
	Width, Height int
	TileMaxSize   int
	Tiles         []image.Rectangle
}
