package chart_test

import (
	"chartographer-go/store"
	"chartographer-go/tile"
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
// TODO дописать тест. Получается интеграционный тест, так как происходит запись на диск
func TestCreateImage(t *testing.T) {
	Convey("Создание тайлов и запись тайлов на диск", t, func() {
		err := store.SetImagesDir("testdata")
		So(err, ShouldBeNil)

		const tileSize = 10
		store.TileMaxSize = tileSize
		_, err = chart.NewRgbaBmp(25, 25)

		// проверка, что созданы изображения-тайлы нужных размеров

		So(err, ShouldBeNil)
	})
}

// TODO не нужно записывать на диск при выполнении теста.
// Необходимо сделать ОС зависимостью и передавать стаб.
func TestNewRGBA(t *testing.T) {
	Convey("init", t, func() {
		const (
			minWidth  = 1
			minHeight = 1
			maxWidth  = 20_000
			maxHeight = 50_000
		)

		store.TileMaxSize = 1000
		err := store.SetImagesDir("testdata")
		So(err, ShouldBeNil)

		// Позитивные тесты

		testSize := func(width, height int) {
			img, err := chart.NewRgbaBmp(width, height)
			So(err, ShouldBeNil)

			So(img.Width, ShouldEqual, width)
			So(img.Height, ShouldEqual, height)
		}
		Convey("MinSize", func() {
			testSize(minWidth, minHeight)
		})
		Convey("MaxSize", func() {
			testSize(maxWidth, maxHeight)
		})

		// Негативные тесты

		var errSize *chart.SizeError

		Convey("test minWidth-1", func() {
			_, err := chart.NewRgbaBmp(minWidth-1, 1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test minHeight-1", func() {
			_, err := chart.NewRgbaBmp(1, minHeight-1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test maxWidth+1", func() {
			_, err := chart.NewRgbaBmp(maxWidth+1, 1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test maxHeight+1", func() {
			_, err := chart.NewRgbaBmp(1, maxHeight+1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
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

func TestCreateTiles(t *testing.T) {
	Convey("Проверка того, что центральные тайлы будут максимального размера, а крайние не максимального", t, func() {
		const (
			width       = 25
			height      = 25
			maxTileSize = 10
		)

		tiles := tile.CreateTiles(width, height, maxTileSize)

		So(tiles, ShouldHaveLength, 9)
		// 1-я строка
		So(tiles[0].Min, ShouldResemble, image.Pt(0, 0))
		So(tiles[0].Max, ShouldResemble, image.Pt(10, 10))

		So(tiles[1].Min, ShouldResemble, image.Pt(10, 0))
		So(tiles[1].Max, ShouldResemble, image.Pt(20, 10))

		So(tiles[2].Min, ShouldResemble, image.Pt(20, 0))
		So(tiles[2].Max, ShouldResemble, image.Pt(25, 10))

		// 2-я строка
		So(tiles[3].Min, ShouldResemble, image.Pt(0, 10))
		So(tiles[3].Max, ShouldResemble, image.Pt(10, 20))

		So(tiles[4].Min, ShouldResemble, image.Pt(10, 10))
		So(tiles[4].Max, ShouldResemble, image.Pt(20, 20))

		So(tiles[5].Min, ShouldResemble, image.Pt(20, 10))
		So(tiles[5].Max, ShouldResemble, image.Pt(25, 20))

		// 3-я строка
		So(tiles[6].Min, ShouldResemble, image.Pt(0, 20))
		So(tiles[6].Max, ShouldResemble, image.Pt(10, 25))

		So(tiles[7].Min, ShouldResemble, image.Pt(10, 20))
		So(tiles[7].Max, ShouldResemble, image.Pt(20, 25))

		So(tiles[8].Min, ShouldResemble, image.Pt(20, 20))
		So(tiles[8].Max, ShouldResemble, image.Pt(25, 25))
	})
}

//func TestCreateImage2(t *testing.T) {
//	Convey("Создание тайлов и запись тайлов на диск", t, func() {
//		err := chart.SetImagesDir("testdata")
//		So(err, ShouldBeNil)
//
//		imgStore := store.New()
//
//		w, h, err := chart.GetImageSize("c59ae0ce-fdda-4442-88ab-6ddb8abad8a0")
//		So(err, ShouldBeNil)
//
//		So(w, ShouldEqual, 255)
//		So(h, ShouldEqual, 255)
//	})
//}
