package chart

import (
	"errors"
	"fmt"
)

var ErrNotExist = errors.New("изображение не найдено")

var ErrNotOverlaps = errors.New("изображение и фрагмент не пересекаются по координатам")

// SizeError означает, что ширина/высота изображения за пределом минимального/максимального значения
type SizeError struct {
	minWidth, width, maxWidth,
	minHeight, height, maxHeight int
}

func (e *SizeError) Error() string {
	return fmt.Sprintf("width должно быть в диапазоне [%d; %d] и ", e.minWidth, e.maxWidth) +
		fmt.Sprintf("height в диапазоне [%d; %d].\n", e.minHeight, e.maxHeight) +
		fmt.Sprintf("Получено width=%d, height=%d", e.width, e.height)
}
