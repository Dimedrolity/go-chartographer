package chart_test

import (
	"errors"
	"image"
	"image/color"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go-chartographer/internal/chart"
	"go-chartographer/pkg/kvstore"
)

//
// Создание изображения
//

// TestTileService - заглушка (stub)
type tileKey struct {
	x, y int
}
type TestTileService struct {
	images map[string]map[tileKey]image.Image
}

func (r *TestTileService) GetTile(id string, x int, y int) (image.Image, error) {
	return r.images[id][tileKey{x: x, y: y}], nil
}
func (r *TestTileService) SaveTile(id string, x int, y int, img image.Image) error {
	_, ok := r.images[id]
	if !ok {
		r.images[id] = make(map[tileKey]image.Image)
	}
	r.images[id][tileKey{x: x, y: y}] = img
	return nil
}
func (r *TestTileService) DeleteImage(id string) error {
	delete(r.images, id)
	return nil
}
func (r *TestTileService) Encode(image.Image) ([]byte, error) {
	panic("implement me")
}
func (r *TestTileService) Decode([]byte) (image.Image, error) {
	panic("implement me")
}

//TestImageRepo - заглушка (stub)
type TestImageRepo struct {
	images map[string]*chart.TiledImage
}

func (r *TestImageRepo) Add(key string, value interface{}) {
	r.images[key] = value.(*chart.TiledImage)
}

func (r *TestImageRepo) Get(key string) (interface{}, error) {
	return r.images[key], nil
}
func (r *TestImageRepo) Delete(id string) error {
	delete(r.images, id)
	return nil
}

type TestTileServiceEmpty struct{}

func (s TestTileServiceEmpty) SaveTile(string, int, int, image.Image) error {
	return nil
}
func (s TestTileServiceEmpty) GetTile(string, int, int) (image.Image, error) {
	return nil, nil
}
func (s TestTileServiceEmpty) DeleteImage(string) error {
	return nil
}
func (s *TestTileServiceEmpty) Encode(image.Image) ([]byte, error) {
	panic("implement me")
}
func (s *TestTileServiceEmpty) Decode([]byte) (image.Image, error) {
	panic("implement me")
}

type TestAdapterEmpty struct{}

func (t TestAdapterEmpty) ShiftRect(image.Image, int, int) {}

func TestNewRGBA(t *testing.T) {
	imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
	tileRepo := &TestTileServiceEmpty{}
	tileMaxSize := 1000
	chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, tileMaxSize)

	Convey("init", t, func() {
		const (
			minWidth  = 1
			minHeight = 1
			maxWidth  = 20_000
			maxHeight = 50_000
		)

		// Позитивные тесты

		testSize := func(width, height int) {
			img, err := chartService.AddImage(width, height)
			So(err, ShouldBeNil)

			So(img.Width, ShouldEqual, width)
			So(img.Height, ShouldEqual, height)
		}
		Convey("MinSize", func() {
			testSize(minWidth, maxHeight)
		})
		Convey("MaxSize", func() {
			testSize(maxWidth, minHeight)
		})

		// Негативные тесты

		var errSize *chart.SizeError

		Convey("test minWidth-1", func() {
			_, err := chartService.AddImage(minWidth-1, 1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test minHeight-1", func() {
			_, err := chartService.AddImage(1, minHeight-1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test maxWidth+1", func() {
			_, err := chartService.AddImage(maxWidth+1, 1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
		Convey("test maxHeight+1", func() {
			_, err := chartService.AddImage(1, maxHeight+1)
			So(errors.As(err, &errSize), ShouldBeTrue)
		})
	})
}

//
// Получение фрагмента изображения
//

func TestGetFragment_In(t *testing.T) {
	Convey("Fragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции Fragment красный пиксель изображения должен появиться в фрагменте", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 1000
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const imgSize = 2
		img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 1
			redY = 1
		)
		img.SetRGBA(redX, redY, red)

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = 1
			y              = 1
			fragmentWidth  = 1
			fragmentHeight = 1
		)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.In(img.Bounds()), ShouldBeTrue)

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Tiles:  []image.Rectangle{img.Bounds()},
		}
		fragment, err := chartService.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

		bounds := fragment.Bounds()
		So(bounds.Dx(), ShouldEqual, fragmentWidth)
		So(bounds.Dy(), ShouldEqual, fragmentHeight)

		So(fragment.At(redX, redY), ShouldResemble, red)
	})
}

func TestGetFragment_PartIntersect(t *testing.T) {
	Convey("Fragment когда прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения\n"+
		"После вызова функции Fragment во фрагменте должен появиться один красный пиксель", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 1000
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

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

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = redX
			y              = redY
			fragmentWidth  = 2
			fragmentHeight = 2
		)

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Tiles:  []image.Rectangle{img.Bounds()},
		}

		// Убеждаемся, что прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.Bounds().Overlaps(img.Bounds()) && !fragmentRect.Bounds().In(img.Bounds()), ShouldBeTrue)

		fragment, err := chartService.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
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

func TestGetFragment_NotOverlaps(t *testing.T) {
	Convey("Fragment когда прямоугольники не пересекаются\n"+
		"Результатом Fragment должно быть полностью черное изображение", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		//adapter := &TestAdapterEmpty{}
		tileMaxSize := 1000
		chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, tileMaxSize)

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

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			x              = imgWidth
			y              = imgHeight
			fragmentWidth  = 1
			fragmentHeight = 1
		)

		// Убеждаемся, что прямоугольники не пересекаются
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.Overlaps(img.Bounds()), ShouldBeFalse)

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),

			Tiles: []image.Rectangle{img.Bounds()},
		}
		_, err := chartService.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)

		So(errors.Is(err, chart.ErrNotOverlaps), ShouldBeTrue)
	})
}

func TestGetFragment_In_NotFirstTile(t *testing.T) {
	Convey("когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения "+
		"и параметры x,y,width,height относятся не к первому тайлу\n"+
		"После вызова функции GetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 10
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			tileX      = 10
			tileY      = 0
			tileWidth  = 5
			tileHeight = 10
		)

		img := image.NewRGBA(image.Rect(tileX, tileY, tileX+tileWidth, tileY+tileHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = tileX
			redY = tileY
		)
		img.SetRGBA(redX, redY, red)

		const id = "0"
		_ = tileRepo.SaveTile(id, tileX, tileY, img) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

		const (
			imgWidth  = 15
			imgHeight = 15
		)
		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  imgWidth,
			Height: imgHeight,
			Tiles:  []image.Rectangle{image.Rect(tileX, tileY, tileX+tileWidth, tileY+tileHeight)},
		}

		const (
			x              = tileX
			y              = tileY
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		imgRect := image.Rect(0, 0, imgWidth, imgHeight)
		So(fragmentRect.In(imgRect), ShouldBeTrue)

		fragment, err := chartService.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

		bounds := fragment.Bounds()
		So(bounds.Dx(), ShouldEqual, fragmentWidth)
		So(bounds.Dy(), ShouldEqual, fragmentHeight)

		const (
			imgRedX = tileX
			imgRedY = tileY
		)
		So(fragment.At(imgRedX, imgRedY), ShouldResemble, red)
	})
}

func TestGetFragment_In_TwoTiles(t *testing.T) {
	Convey("GetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения "+
		"и фрагмент затрагивает 2 тайла.\n"+
		"Результат GetFragment должен иметь пиксели изображения", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 10
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			tile1X      = 0
			tile1Y      = 0
			tile1Width  = 10
			tile1Height = 10
		)
		const (
			tile2X      = 10
			tile2Y      = 0
			tile2Width  = 5
			tile2Height = 10
		)

		t1 := image.NewRGBA(image.Rect(tile1X, tile1Y, tile1X+tile1Width, tile1Y+tile1Height))
		t2 := image.NewRGBA(image.Rect(tile2X, tile2Y, tile2X+tile2Width, tile2Y+tile2Height))

		const (
			redX = 9
			redY = 0
		)
		const (
			greenX = 10
			greenY = 0
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		t1.SetRGBA(redX, redY, red)
		t2.SetRGBA(greenX, greenY, green)

		const id = "0"
		_ = tileRepo.SaveTile(id, tile1X, tile1Y, t1) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб
		_ = tileRepo.SaveTile(id, tile2X, tile2Y, t2) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		const (
			imgWidth  = 15
			imgHeight = 15
		)

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  imgWidth,
			Height: imgHeight,
			Tiles: []image.Rectangle{
				image.Rect(tile1X, tile1Y, tile1X+tile1Width, tile1Y+tile1Height),
				image.Rect(tile2X, tile2Y, tile2X+tile2Width, tile2Y+tile2Height),
			},
		}
		imageRepo.Add(id, tiledImg)

		const (
			x              = 9
			y              = 0
			fragmentWidth  = 2
			fragmentHeight = 1
		)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		imgRect := image.Rect(0, 0, imgWidth, imgHeight)
		fragmentRect := image.Rect(x, y, x+fragmentWidth, y+fragmentHeight)
		So(fragmentRect.In(imgRect), ShouldBeTrue)

		fragment, err := chartService.GetFragment(tiledImg, x, y, fragmentWidth, fragmentHeight)
		So(err, ShouldBeNil)

		So(fragment.At(redX, redY), ShouldResemble, red)
		So(fragment.At(greenX, greenY), ShouldResemble, green)
	})
}

func TestGetFragment_Size(t *testing.T) {
	imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
	tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
	adapter := &TestAdapterEmpty{}
	tileMaxSize := 1000
	chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

	emptyImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	const id = "0"
	_ = tileRepo.SaveTile(id, 0, 0, emptyImg) // чтобы getTile, вызываемый в chart.GetFragment, возвращал стаб

	tiledEmptyImg := &chart.TiledImage{
		Id:     id,
		Width:  emptyImg.Bounds().Dx(),
		Height: emptyImg.Bounds().Dy(),
		Tiles:  []image.Rectangle{emptyImg.Bounds()},
	}

	const (
		fragmentMinWidth  = 1
		fragmentMinHeight = 1
		fragmentMaxWidth  = 5_000
		fragmentMaxHeight = 5_000
	)

	// Позитивные тесты

	testSize := func(width, height int) {
		img, err := chartService.GetFragment(tiledEmptyImg, 0, 0, width, height)
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
		_, err := chartService.GetFragment(tiledEmptyImg, 0, 0, fragmentMinWidth-1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test minHeight-1", t, func() {
		_, err := chartService.GetFragment(tiledEmptyImg, 0, 0, 1, fragmentMinHeight-1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxWidth+1", t, func() {
		_, err := chartService.GetFragment(tiledEmptyImg, 0, 0, fragmentMaxWidth+1, 1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
	Convey("test maxHeight+1", t, func() {
		_, err := chartService.GetFragment(tiledEmptyImg, 0, 0, 1, fragmentMaxHeight+1)
		So(errors.As(err, &errSize), ShouldBeTrue)
	})
}

// -----------
// SetFragment
// -----------
func TestSetFragment_In(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения.\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 1000
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Tiles:  []image.Rectangle{img.Bounds()},
		}

		const (
			x = 1
			y = 0
		)
		const (
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))
		const (
			fragmentRedX = x
			fragmentRedY = y
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(fragmentRedX, fragmentRedY, red)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		So(fragment.Bounds().In(img.Bounds()), ShouldBeTrue)

		err := chartService.SetFragment(tiledImg, 0, 0, fragment)
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
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 1000
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Tiles:  []image.Rectangle{img.Bounds()},
		}

		const (
			x = 3
			y = 3
		)
		const (
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))

		const (
			fragmentRedX = x
			fragmentRedY = y
		)
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(fragmentRedX, fragmentRedY, red)

		// Убеждаемся, что прямоугольники не пересекаются
		So(!fragment.Bounds().Overlaps(img.Bounds()), ShouldBeTrue)

		err := chartService.SetFragment(tiledImg, 0, 0, fragment)
		So(errors.Is(err, chart.ErrNotOverlaps), ShouldBeTrue)

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
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		tileMaxSize := 1000
		adapter := &TestAdapterEmpty{}
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			imgWidth  = 2
			imgHeight = 2
		)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

		const id = "0"
		_ = tileRepo.SaveTile(id, 0, 0, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Tiles:  []image.Rectangle{img.Bounds()},
		}

		const (
			x              = 1
			y              = 1
			fragmentWidth  = 2
			fragmentHeight = 2
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))

		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		for x := x; x < fragmentWidth; x++ {
			for y := y; y < fragmentHeight; y++ {
				fragment.SetRGBA(x, y, red)
			}
		}

		const (
			imgRedX = 1
			imgRedY = 1
		)
		// Убеждаемся, что прямоугольники пересекаются, но фрагмент частично вне прямоугольника изображения
		So(fragment.Bounds().Overlaps(img.Bounds()) && !fragment.Bounds().In(img.Bounds()), ShouldBeTrue)

		err := chartService.SetFragment(tiledImg, 0, 0, fragment)
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

func TestSetFragment_In_NotFirstTile(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения "+
		"и параметры x,y,width,height относятся не к первому тайлу\n"+
		"После вызова функции SetFragment красный пиксель фрагмента должен появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 10
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			tileX      = 10
			tileY      = 0
			tileWidth  = 5
			tileHeight = 10
		)
		img := image.NewRGBA(image.Rect(tileX, tileY, tileX+tileWidth, tileY+tileHeight))

		const id = "0"
		_ = tileRepo.SaveTile(id, tileX, tileY, img) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		const (
			imgWidth  = 15
			imgHeight = 15
		)
		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  imgWidth,
			Height: imgHeight,
			Tiles:  []image.Rectangle{img.Bounds()},
		}

		const (
			x = tileX
			y = tileY
		)
		const (
			fragmentWidth  = 1
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fragment.SetRGBA(x, y, red)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		imgRect := image.Rect(0, 0, imgWidth, imgHeight)
		So(fragment.Bounds().In(imgRect), ShouldBeTrue)

		err := chartService.SetFragment(tiledImg, 0, 0, fragment)
		So(err, ShouldBeNil)

		So(img.At(x, y), ShouldResemble, red)
	})
}

func TestSetFragment_In_TwoTiles(t *testing.T) {
	Convey("SetFragment когда прямоугольник фрагмента полностью лежит в прямоугольнике изображения "+
		"и фрагмент затрагивает 2 тайла.\n"+
		"После вызова функции SetFragment красные пиксели фрагмента должны появиться в изображении", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		adapter := &TestAdapterEmpty{}
		tileMaxSize := 10
		chartService := chart.NewChartographerService(imageRepo, tileRepo, adapter, tileMaxSize)

		const (
			tile1X = 0
			tile1Y = 0
			// TODO tile1_x0 _x1
			tile1Width  = 10
			tile1Height = 10
		)

		const (
			tile2X      = 10
			tile2Y      = 0
			tile2Width  = 5
			tile2Height = 10
		)

		t1 := image.NewRGBA(image.Rect(tile1X, tile1Y, tile1X+tile1Width, tile1Y+tile1Height))
		t2 := image.NewRGBA(image.Rect(tile2X, tile2Y, tile2X+tile2Width, tile2Y+tile2Height))

		const id = "0"
		_ = tileRepo.SaveTile(id, tile1X, tile1Y, t1) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб
		_ = tileRepo.SaveTile(id, tile2X, tile2Y, t2) // чтобы getTile, вызываемый в chart.SetFragment, возвращал стаб

		const (
			imgWidth  = 15
			imgHeight = 15
		)

		tiledImg := &chart.TiledImage{
			Id:     id,
			Width:  imgWidth,
			Height: imgHeight,
			Tiles: []image.Rectangle{
				image.Rect(tile1X, tile1Y, tile1X+tile1Width, tile1Y+tile1Height),
				image.Rect(tile2X, tile2Y, tile2X+tile2Width, tile2Y+tile2Height),
			},
		}

		const (
			x = 9
			y = 0
		)
		const (
			fragmentWidth  = 2
			fragmentHeight = 1
		)
		fragment := image.NewRGBA(image.Rect(x, y, x+fragmentWidth, y+fragmentHeight))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		fragment.SetRGBA(x, y, red)
		fragment.SetRGBA(x+1, y, green)

		// Убеждаемся, что прямоугольник фрагмента полностью лежит в прямоугольнике изображения
		imgRect := image.Rect(0, 0, imgWidth, imgHeight)
		So(fragment.Bounds().In(imgRect), ShouldBeTrue)

		err := chartService.SetFragment(tiledImg, 0, 0, fragment)
		So(err, ShouldBeNil)

		const (
			imgRedX = 9
			imgRedY = 0
		)
		const (
			imgGreenX = 10
			imgGreenY = 0
		)
		So(t1.At(imgRedX, imgRedY), ShouldResemble, red)
		So(t2.At(imgGreenX, imgGreenY), ShouldResemble, green)
	})
}

//
// Удаление изображения
//

func TestDeleteImage_Success(t *testing.T) {
	Convey("Должно удалить все данные изображения", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		//adapter := &TestAdapterEmpty{}
		chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, 0)

		const id = "0"
		tiledImg := &chart.TiledImage{
			Id: id,
		}
		imageRepo.Add(id, tiledImg)

		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		_ = tileRepo.SaveTile(id, 0, 0, img)

		getImg, _ := imageRepo.Get(id)
		So(getImg, ShouldNotBeNil)

		getTile, _ := tileRepo.GetTile(id, 0, 0)
		So(getTile, ShouldNotBeNil)

		err := chartService.DeleteImage(id)
		So(err, ShouldBeNil)

		getImg, _ = imageRepo.Get(id)
		So(getImg, ShouldBeNil)

		getTile, _ = tileRepo.GetTile(id, 0, 0)
		So(getTile, ShouldBeNil)
	})
}

type TestImageRepoDeleteNotExist struct{}

func (r *TestImageRepoDeleteNotExist) Add(string, interface{}) {
	panic("implement me")
}
func (r *TestImageRepoDeleteNotExist) Get(string) (interface{}, error) {
	panic("implement me")
}
func (r *TestImageRepoDeleteNotExist) Delete(string) error {
	return kvstore.ErrNotExist
}

func TestDeleteImage_NotExist(t *testing.T) {
	Convey("При запросе несуществующего изображения должна быть ошибка", t, func() {
		imageRepo := &TestImageRepoDeleteNotExist{}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		//adapter := &TestAdapterEmpty{}
		chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, 0)

		err := chartService.DeleteImage("0")
		So(errors.Is(err, chart.ErrNotExist), ShouldBeTrue)
	})
}

//
// Получение изображения
//

type TestImageRepoGetNotExist struct{}

func (r *TestImageRepoGetNotExist) Add(string, interface{}) {
	panic("implement me")
}
func (r *TestImageRepoGetNotExist) Get(string) (interface{}, error) {
	return nil, kvstore.ErrNotExist
}
func (r *TestImageRepoGetNotExist) Delete(string) error {
	panic("implement me")
}

func TestGetImage_NotExist(t *testing.T) {
	Convey("При запросе несуществующего изображения должна быть ошибка", t, func() {
		imageRepo := &TestImageRepoGetNotExist{}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		//adapter := &TestAdapterEmpty{}
		chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, 0)

		_, err := chartService.GetImage("0")
		So(errors.Is(err, chart.ErrNotExist), ShouldBeTrue)
	})
}

func TestGetImage_Success(t *testing.T) {
	Convey("Должно удалить все данные изображения", t, func() {
		imageRepo := &TestImageRepo{images: make(map[string]*chart.TiledImage)}
		tileRepo := &TestTileService{images: make(map[string]map[tileKey]image.Image)}
		//adapter := &TestAdapterEmpty{}
		chartService := chart.NewChartographerService(imageRepo, tileRepo, nil, 0)

		const id = "0"
		tiledImg := &chart.TiledImage{
			Id: id,
		}
		imageRepo.Add(id, tiledImg)

		img, err := chartService.GetImage(id)
		So(err, ShouldBeNil)
		So(img, ShouldResemble, tiledImg)
	})
}
