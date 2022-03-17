package store

import (
	"image"
)

// ShiftRect смещает прямоугольник изображения img на координаты (x; y).
// В итоге img.Bounds() будет возвращать смещенный прямоугольник.
func ShiftRect(img image.Image, x, y int) {
	switch i := img.(type) {
	case *image.RGBA:
		i.Rect = i.Rect.Add(image.Pt(x, y))
	default:
		panic("not implemented")
	}
}
