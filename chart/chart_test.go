package chart_test

import (
	"errors"
	"image"
	"image/color"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"chartographer-go/chart"
)

// -----------
// NewRGBA
// -----------
func TestNewRGBA(t *testing.T) {
	const (
		minWidth  = 1
		minHeight = 1
		maxWidth  = 20_000
		maxHeight = 50_000
	)

	// Позитивные тесты

	testSize := func(width, height int) {
		img, err := chart.NewRGBA(width, height)
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

	// Негативные тесты

	var errSize *chart.SizeError

	Convey("test minWidth-1", t, func() {
		_, err := chart.NewRGBA(minWidth-1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test minHeight-1", t, func() {
		_, err := chart.NewRGBA(1, minHeight-1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxWidth+1", t, func() {
		_, err := chart.NewRGBA(maxWidth+1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxHeight+1", t, func() {
		_, err := chart.NewRGBA(1, maxHeight+1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
}

// -----------
// Fragment
// -----------
func TestFragment_In(t *testing.T) {
	Convey("Fragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции Fragment красный пиксель изображения должен появиться в фрагменте", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		const (
			x              = 1
			y              = 1
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment, err := chart.Fragment(img, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)
		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		So(fragment.Bounds().In(img.Bounds()), ShouldBeTrue)

		bounds := fragment.Bounds()
		So(bounds.Dx(), ShouldEqual, fragmentWidth)
		So(bounds.Dy(), ShouldEqual, fragmentHeight)

		So(fragment.At(redX, redY), ShouldResemble, red)
	})
}

func TestFragment_PartIntersect(t *testing.T) {
	Convey("Fragment когда прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения\n"+
		"После вызова функции Fragment во фрагменте должен появиться один красный пиксель", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		const (
			x              = redX
			y              = redY
			fragmentWidth  = 2
			fragmentHeight = 2
		)
		fragment, err := chart.Fragment(img, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)
		// Убеждаемся, что прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения
		So(fragment.Bounds().Overlaps(img.Bounds()) && !fragment.Bounds().In(img.Bounds()), ShouldBeTrue)

		bounds := fragment.Bounds()
		So(bounds.Dx(), ShouldEqual, fragmentWidth)
		So(bounds.Dy(), ShouldEqual, fragmentHeight)

		for y := 0; y < fragmentHeight; y++ {
			for x := 0; x < fragmentWidth; x++ {
				c := color.RGBA{}
				if x == redX && y == redY {
					c = red
				}

				So(fragment.At(x, y), ShouldResemble, c)
			}
		}
	})
}

func TestFragment_NotOverlaps(t *testing.T) {
	Convey("Fragment когда прямоугольники не пересекаются\n"+
		"Результатом Fragment должно быть полностью черное изображение", t, func() {
		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		const (
			x              = imgWidth
			y              = imgHeight
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		So(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight).Bounds().Overlaps(img.Bounds()), ShouldBeFalse)

		_, err := chart.Fragment(img, x, y, fragmentWidth, fragmentHeight)
		So(errors.Is(err, chart.ErrNotOverlaps), ShouldBeTrue)
	})
}

func TestFragment_Size(t *testing.T) {
	const (
		fragmentMinWidth  = 1
		fragmentMinHeight = 1
		fragmentMaxWidth  = 5_000
		fragmentMaxHeight = 5_000
	)
	emptyImg := image.NewRGBA(image.Rect(0, 0, 1, 1))

	// Позитивные тесты

	testSize := func(width, height int) {
		img, err := chart.Fragment(emptyImg, 0, 0, width, height)
		So(err, ShouldBeNil)

		rect := img.Bounds()
		So(rect.Dx(), ShouldEqual, width)
		So(rect.Dy(), ShouldEqual, height)
	}
	Convey("MinSize", t, func() {
		testSize(fragmentMinWidth, fragmentMinHeight)
	})
	Convey("MaxSize", t, func() {
		testSize(fragmentMaxWidth, fragmentMaxHeight)
	})

	// Негативные тесты

	var errSize *chart.SizeError
	Convey("test minWidth-1", t, func() {
		_, err := chart.Fragment(emptyImg, 0, 0, fragmentMinWidth-1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test minHeight-1", t, func() {
		_, err := chart.Fragment(emptyImg, 0, 0, 1, fragmentMinHeight-1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxWidth+1", t, func() {
		_, err := chart.Fragment(emptyImg, 0, 0, fragmentMaxWidth+1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxHeight+1", t, func() {
		_, err := chart.Fragment(emptyImg, 0, 0, 1, fragmentMaxHeight+1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
}

// -----------
// SetFragment
// -----------
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

		err := chart.SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

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
		"После вызова функции SetFragment должна быть ошибка, "+
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

		err := chart.SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldNotBeNil)
	})
}

func TestSetFragment_NotOverlaps(t *testing.T) {
	Convey("SetFragment когда прямоугольники не пересекаются\n"+
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

		err := chart.SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldNotBeNil)

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

		err := chart.SetFragment(img, fragment, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

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
