// Package tile отвечает за разделение изображение на части (тайлы).
package tile

import (
	"image"
	"image/draw"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxSize определяет максимальный размер тайла по ширине и высоте.
// Необходимо проинициализировать перед использованием функций текущего pkg
// TODO выделить в структуру Tiler, и фукнцию NewTiler(tileMaxSize). Тогда сделать все фукнции методами Tiler
var MaxSize int

// CreateTiles делит прямоугольник указанного размера (width и height) на несколько тайлов (прямоугольников)
// с максимальным размером tileMaxSize
func CreateTiles(width, height, tileMaxSize int) []image.Rectangle {
	tiles := make([]image.Rectangle, 0, (width/tileMaxSize+1)*(height/tileMaxSize+1))

	for y := 0; y < height; y += tileMaxSize {
		for x := 0; x < width; x += tileMaxSize {
			w := min(width-x, tileMaxSize)
			h := min(height-y, tileMaxSize)
			tiles = append(tiles, image.Rect(x, y, x+w, y+h))
		}
	}

	return tiles
}

// FilterOverlappedTiles возвращает только те тайлы, которые пересекаются с фрагментом.
// TODO тест
func FilterOverlappedTiles(imgTiles []image.Rectangle, fragment image.Rectangle) []image.Rectangle {
	overlapped := make([]image.Rectangle, 0, len(imgTiles))
	for _, tile := range imgTiles {
		if tile.Overlaps(fragment) {
			overlapped = append(overlapped, tile)
		}
	}
	return overlapped
}

// DrawIntersection закрашивает dst пикселями tile, которые пересекаются с fragment
func DrawIntersection(dst draw.Image, tile image.Image, fragment image.Rectangle) {
	intersect := tile.Bounds().Intersect(fragment)
	draw.Draw(dst, intersect, tile, intersect.Min, draw.Src)
}
