package chart

import (
	"image"
)

// RectShifter - адаптер (паттерн) для смещения прямоугольника изображения.
type RectShifter interface {
	ShiftRect(img image.Image, x, y int)
}

type ImageAdapter struct{}

// ShiftRect смещает прямоугольник изображения img на координаты (x; y).
// В итоге img.Bounds() будет возвращать смещенный прямоугольник.
func (a *ImageAdapter) ShiftRect(img image.Image, x, y int) {
	switch i := img.(type) {
	case *image.RGBA:
		i.Rect = i.Rect.Add(image.Pt(x, y))
	default:
		panic("color model is not implemented")
	}
}
