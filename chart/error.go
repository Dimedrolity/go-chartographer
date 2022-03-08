package chart

import (
	"errors"
	"fmt"
)

var ErrNotOverlaps = errors.New("изображение и фрагмент не пересекаются по координатам")

// SizeError означает, что ширина/высота изображения за пределом минимального/максимального значения
type SizeError struct {
	minWidth, width, maxWidth,
	minHeight, height, maxHeight int
}

func (err *SizeError) Error() string {
	return fmt.Sprintf("width должно быть в диапазоне [%d; %d] и ", err.minWidth, err.maxWidth) +
		fmt.Sprintf("height в диапазоне [%d; %d].\n", err.minHeight, err.maxHeight) +
		fmt.Sprintf("Получено width=%d, height=%d", err.width, err.height)
}
