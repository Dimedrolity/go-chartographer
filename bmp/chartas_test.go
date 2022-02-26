package bmp

import (
	. "github.com/smartystreets/goconvey/convey"
	"image"
	"image/color"
	"testing"
)

func TestNewImage_ZeroWidth(t *testing.T) {
	Convey("ZeroWidth", t, func() {
		_, err := NewImage(0, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestNewImage_ZeroHeight(t *testing.T) {
	Convey("ZeroHeight", t, func() {
		_, err := NewImage(1, 0)
		So(err, ShouldNotBeNil)
	})
}

func TestNewImage_WidthExceeded(t *testing.T) {
	Convey("WidthExceeded", t, func() {
		_, err := NewImage(maxWidth+1, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestNewImage_HeightExceeded(t *testing.T) {
	Convey("HeightExceeded", t, func() {
		_, err := NewImage(1, maxHeight+1)
		So(err, ShouldNotBeNil)
	})
}

func TestNewImage_MinMaxSize(t *testing.T) {
	testSize := func(width, height int) {
		img, err := NewImage(width, height)
		So(err, ShouldBeNil)

		rect := img.Bounds()
		So(rect.Dx(), ShouldEqual, width)
		So(rect.Dy(), ShouldEqual, height)
	}
	Convey("MinSize", t, func() {
		testSize(minWidth, minHeight)
	})
	Convey("MaxSize", t, func() {
		testSize(maxWidth, maxHeight)
	})
}

func TestSubImage(t *testing.T) {
	Convey("SubImage", t, func() {
		img := image.NewRGBA(image.Rect(0, 0, 2, 2))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		const (
			subWidth  = 1
			subHeight = 1
		)
		sub, err := SubImage(img, redX, redY, subWidth, subHeight)
		So(err, ShouldBeNil)

		bounds := sub.Bounds()
		So(bounds.Dx(), ShouldEqual, subWidth)
		So(bounds.Dy(), ShouldEqual, subHeight)

		So(sub.At(redX, redY), ShouldResemble, red)
	})
}

// TODO тест, когда SubImage выходит за рамки исходного
// чёрным цветом закрашивается та часть фрагмента, которая оказывается вне границ изображения (см. пример ниже).
// Получается нужно будет вернуть sub-image запрошенного размера, но часть будет черным цветом.
// TODO уточнить у авторов.
