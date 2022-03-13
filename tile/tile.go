// Package tile отвечает за разбиение изображение на части (тайлы).
package tile

import (
	"image"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
