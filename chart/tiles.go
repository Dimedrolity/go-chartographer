package chart

import (
	"github.com/google/uuid"
	"image"
	"os"
	"path/filepath"
	"strconv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CreateImage разбивает размеры изображения на тайлы, создает RGBA изображения в соответствии с тайлами,
// записывает изображения на диск в формате BMP
// TODO разбить функцию, сейчас она делает все подряд. Возможно, вынести обработку SizeError в CreateTiles
func CreateImage(width, height, maxTileSize int) (string, error) {
	if width < minWidth || width > maxWidth ||
		height < minHeight || height > maxHeight {
		return "", &SizeError{
			minWidth: minWidth, width: width, maxWidth: maxWidth,
			minHeight: minHeight, height: height, maxHeight: maxHeight,
		}
	}

	tiles := CreateTiles(width, height, maxTileSize)
	id := uuid.NewString()

	for _, tile := range tiles {
		img := image.NewRGBA(tile)
		err := SaveTiledImage(id, img)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

// CreateTiles делит прямоугольник указанного размера (width и height) на несколько тайлов (прямоугольников)
// с максимальным размером maxTileSize
func CreateTiles(width, height, maxTileSize int) []image.Rectangle {
	tiles := make([]image.Rectangle, 0, (width/maxTileSize+1)*(height/maxTileSize+1))

	for y := 0; y < height; y += maxTileSize {
		for x := 0; x < width; x += maxTileSize {
			w := min(width-x, maxTileSize)
			h := min(height-y, maxTileSize)
			tiles = append(tiles, image.Rect(x, y, x+w, y+h))
		}
	}

	return tiles
}

// SaveTiledImage сохраняет тайл-изображение на диск.
// Необходим id и img.Bounds для именования папок и файлов.
// Необходим сам img, чтобы сделать Encode и получить байты.
func SaveTiledImage(id string, img image.Image) error {
	// TODO можно было бы обойтись без буфера, Create файл и bmp.Encode(файл)
	// 	file, err := os.OpenFile(filepath.Join(dir, x+".bmp"), os.O_WRONLY|os.O_CREATE, 0777)
	encode, err := Encode(img)
	if err != nil {
		return err
	}

	y := strconv.Itoa(img.Bounds().Min.Y)
	dir := filepath.Join(dirPath, id, y)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	x := strconv.Itoa(img.Bounds().Min.X)
	err = os.WriteFile(filepath.Join(dir, x+".bmp"), encode, 0777)

	if err != nil {
		return err
	}

	return nil
}
