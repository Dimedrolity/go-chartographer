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

func TestSetFragment_In(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const (
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(0, 0, fragmentWidth, fragmentHeight))
		const (
			fragmentRedX = 0
			fragmentRedY = 0
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(fragmentRedX, fragmentRedY, red)

		const (
			x = 1
			y = 0
		)
		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		rect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(rect.In(img.Bounds()), ShouldBeTrue)

		SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)

		const (
			imgRedX = x
			imgRedY = y
		)
		for x := 0; x < imgWidth; x++ {
			for y := 0; y < imgHeight; y++ {
				c := color.RGBA{}
				if x == imgRedX && y == imgRedY {
					c = red
				}

				So(img.At(x, y), ShouldResemble, c)
			}
		}
	})
}

func TestSetFragment_In_FragmentWrongStart(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения\n"+
		"После вызова функции SetFragment красный пиксель фрагмента не должен появиться в изображении, "+
		"так как прямоугольник фрагмента имеет начальные координаты не соотв. функции", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const (
			x              = 1
			y              = 0
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))

		const (
			redX = 1
			redY = 0
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(redX, redY, red)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		So(fragment.Bounds().In(img.Bounds()), ShouldBeTrue)

		SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)

		for x := 0; x < imgWidth; x++ {
			for y := 0; y < imgHeight; y++ {
				So(img.At(x, y), ShouldResemble, color.RGBA{})
			}
		}
	})
}

func TestSetFragment_NotOverlaps(t *testing.T) {
	Convey("SetFragment прямоугольники не пересекаются\n"+
		"После вызова функции SetFragment красный пиксель фрагмента не должен появиться в изображении, "+
		"так как прямоугольники не пересекаются", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const (
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(0, 0, fragmentWidth, fragmentHeight))
		const (
			fragmentRedX = 0
			fragmentRedY = 0
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(fragmentRedX, fragmentRedY, red)

		const (
			x = 3
			y = 3
		)
		rect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		// Убеждаемся, что прямоугольники не пересекаются
		So(!rect.Overlaps(img.Bounds()), ShouldBeTrue)

		SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)

		for x := 0; x < imgWidth; x++ {
			for y := 0; y < imgHeight; y++ {
				c := color.RGBA{}
				So(img.At(x, y), ShouldResemble, c)
			}
		}
	})
}

func TestSetFragment_PartIntersect(t *testing.T) {
	Convey("SetFragment когда прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const (
			fragmentWidth  = 2
			fragmentHeight = 2
		)
		fragment := image.NewRGBA(image.Rect(0, 0, fragmentWidth, fragmentHeight))

		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		for x := 0; x < fragmentWidth; x++ {
			for y := 0; y < fragmentHeight; y++ {
				fragment.SetRGBA(x, y, red)
			}
		}

		const (
			x       = 1
			y       = 1
			imgRedX = 1
			imgRedY = 1
		)
		rect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		// Убеждаемся, что прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения
		So(rect.Bounds().Overlaps(img.Bounds()) && !rect.Bounds().In(img.Bounds()), ShouldBeTrue)

		SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)

		for x := 0; x < imgWidth; x++ {
			for y := 0; y < imgHeight; y++ {
				c := color.RGBA{}
				if x == imgRedX && y == imgRedY {
					c = red
				}

				So(img.At(x, y), ShouldResemble, c)
			}
		}
	})
}
