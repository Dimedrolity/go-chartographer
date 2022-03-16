package chart_test

import (
	"chartographer-go/tile"
	"chartographer-go/tiledimage"
	"errors"
	"image"
	"image/color"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"chartographer-go/chart"
)

//
// Создание изображения
//

// TestTileRepository - заглушка (stub)
type TestTileRepository struct {
	images map[string]image.Image
}

func (r *TestTileRepository) GetTile(id string, _ int, _ int) (image.Image, error) {
	return r.images[id], nil
}
func (r *TestTileRepository) SaveTile(id string, img image.Image) error {
	r.images[id] = img
	return nil
}
func (r *TestTileRepository) DeleteImage(id string) error {
	delete(r.images, id)
	return nil
}

//TestImageRepo - заглушка (stub)
type TestImageRepo struct {
	images map[string]*tiledimage.Image
}

func (r *TestImageRepo) Add(img *tiledimage.Image) {
	r.images[img.Id] = img
}
func (r *TestImageRepo) Get(id string) (*tiledimage.Image, error) {
	return r.images[id], nil
}
func (r *TestImageRepo) Delete(id string) error {
	delete(r.images, id)
	return nil
}

func TestNewRGBA(t *testing.T) {
	chart.ImageRepo = &TestImageRepo{images: make(map[string]*tiledimage.Image)}

	Convey("init", t, func() {
		const (
			minWidth  = 1
			minHeight = 1
			maxWidth  = 20_000
			maxHeight = 50_000
		)

		tile.MaxSize = 1000
		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

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

//
// Получение фрагмента изображения
//

func TestFragment_In(t *testing.T) {
	Convey("Fragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции Fragment красный пиксель изображения должен появиться в фрагменте", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

		const imgSize = 2
		img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = 1
			y              = 1
			fragmentWidth  = 1
			fragmentHeight = 1
		)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.In(img.Bounds()), ShouldBeTrue)

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		fragment, err := chart.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

		bounds := fragment.Bounds()
		So(bounds.Dx(), ShouldEqual, fragmentWidth)
		So(bounds.Dy(), ShouldEqual, fragmentHeight)

		So(fragment.At(redX, redY), ShouldResemble, red)
	})
}

func TestFragment_PartIntersect(t *testing.T) {
	Convey("Fragment когда прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения\n"+
		"После вызова функции Fragment во фрагменте должен появиться один красный пиксель", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

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

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = redX
			y              = redY
			fragmentWidth  = 2
			fragmentHeight = 2
		)

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}

		// Убеждаемся, что прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.Bounds().Overlaps(img.Bounds()) && !fragmentRect.Bounds().In(img.Bounds()), ShouldBeTrue)

		fragment, err := chart.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

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
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

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

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = imgWidth
			y              = imgHeight
			fragmentWidth  = 1
			fragmentHeight = 1
		)

		// Убеждаемся, что прямоугольники не пересекаются
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.Overlaps(img.Bounds()), ShouldBeFalse)

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		_, err := chart.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)

		So(errors.Is(err, chart.ErrNotOverlaps), ShouldBeTrue)
	})
}

func TestFragment_Size(t *testing.T) {
	imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
	chart.ImageRepo = imageRepo

	tile.MaxSize = 1000

	tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
	chart.TileRepo = tileRepo

	emptyImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	id := "0"
	_ = tileRepo.SaveTile(id, emptyImg) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

	tiledEmptyImg := &tiledimage.Image{
		Id: id,
		Config: image.Config{
			Width:  emptyImg.Bounds().Dx(),
			Height: emptyImg.Bounds().Dy(),
		},
		TileMaxSize: tile.MaxSize,
		Tiles:       []image.Rectangle{emptyImg.Bounds()},
	}

	const (
		fragmentMinWidth  = 1
		fragmentMinHeight = 1
		fragmentMaxWidth  = 5_000
		fragmentMaxHeight = 5_000
	)

	// Позитивные тесты

	testSize := func(width, height int) {
		img, err := chart.GetFragment(tiledEmptyImg, 0, 0, width, height)
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
		_, err := chart.GetFragment(tiledEmptyImg, 0, 0, fragmentMinWidth-1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test minHeight-1", t, func() {
		_, err := chart.GetFragment(tiledEmptyImg, 0, 0, 1, fragmentMinHeight-1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxWidth+1", t, func() {
		_, err := chart.GetFragment(tiledEmptyImg, 0, 0, fragmentMaxWidth+1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxHeight+1", t, func() {
		_, err := chart.GetFragment(tiledEmptyImg, 0, 0, 1, fragmentMaxHeight+1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
}

// -----------
// SetFragment
// -----------
func TestSetFragment_In(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		imageRepo.Add(tiledImg)

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

		err := chart.SetFragment(id, fragment, x, y, fragmentWidth, fragmentHeight)
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

func TestSetFragment_In_FragmentWithOffset(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения"+
		"и прямоугольник фрагмента имеет начальные координаты не (0; 0).\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		imageRepo.Add(tiledImg)

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

		err := chart.SetFragment(id, fragment, x, y, fragmentWidth, fragmentHeight)
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

func TestSetFragment_NotOverlaps(t *testing.T) {
	Convey("SetFragment когда прямоугольники не пересекаются\n"+
		"После вызова функции SetFragment красный пиксель фрагмента не должен появиться в изображении, "+
		"так как прямоугольники не пересекаются", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		imageRepo.Add(tiledImg)

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

		err := chart.SetFragment(id, fragment, x, y, fragmentWidth, fragmentHeight)
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
		imageRepo := &TestImageRepo{images: make(map[string]*tiledimage.Image)}
		chart.ImageRepo = imageRepo

		tile.MaxSize = 1000

		tileRepo := &TestTileRepository{images: make(map[string]image.Image)}
		chart.TileRepo = tileRepo

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		id := "0"
		_ = tileRepo.SaveTile(id, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &tiledimage.Image{
			Id: id,
			Config: image.Config{
				Width:  img.Bounds().Dx(),
				Height: img.Bounds().Dy(),
			},
			TileMaxSize: tile.MaxSize,
			Tiles:       []image.Rectangle{img.Bounds()},
		}
		imageRepo.Add(tiledImg)

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

		err := chart.SetFragment(id, fragment, x, y, fragmentWidth, fragmentHeight)
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
