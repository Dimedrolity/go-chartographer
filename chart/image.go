package chart

import "image"

// TiledImage - изображение, разделенное на части тайлы.
// Деление на тайлы необходимо, чтобы приложение не помещать в оперативную память изображения больших размеров.
type TiledImage struct {
	Id            string
	Width, Height int
	TileMaxSize   int // Определяет максимальный размер тайла по ширине и высоте.
	Tiles         []image.Rectangle
}
