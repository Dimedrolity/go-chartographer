package imagetile_test

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"go-chartographer/internal/imagetile"
	"image"
	"image/color"
	"testing"
)

type TestTileRepo struct {
	images map[string][]byte
}

func (r *TestTileRepo) SaveTile(id string, _ int, _ int, img []byte) error {
	r.images[id] = img
	return nil
}

func (r *TestTileRepo) GetTile(id string, _, _ int) ([]byte, error) {
	_, ok := r.images[id]
	if !ok {
		return nil, errors.New("")
	}
	return r.images[id], nil
}

func (r *TestTileRepo) DeleteImage(id string) error {
	delete(r.images, id)
	return nil
}

func TestBmpService_SuccessGet(t *testing.T) {
	Convey("Проверка сохранения и получения изображения."+
		"Координаты полученного изображения должны соответствовать исходынм.", t, func() {
		tileRepo := &TestTileRepo{images: make(map[string][]byte)}
		bmpService := imagetile.NewBmpService(tileRepo)

		const (
			x      = 1
			y      = 1
			width  = 1
			height = 1
		)
		img := image.NewRGBA(image.Rect(x, y, x+width, y+height))
		img.Set(x, y, color.RGBA{A: 255}) // чтобы в SaveTile Decode распознал как 24-битное

		id := "0"

		err := bmpService.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		got, err := bmpService.GetTile(id, x, y)
		So(err, ShouldBeNil)

		So(got.Bounds(), ShouldResemble, img.Bounds())
	})
}

func TestBmpService_SuccessDelete(t *testing.T) {
	Convey("После создания и удаления вызов фукнции получения должен вернуть ошибку.", t, func() {
		tileRepo := &TestTileRepo{images: make(map[string][]byte)}
		bmpService := imagetile.NewBmpService(tileRepo)

		const (
			x      = 0
			y      = 0
			width  = 1
			height = 1
		)
		img := image.NewRGBA(image.Rect(x, y, x+width, y+height))
		img.Set(x, y, color.RGBA{A: 0xFF}) // так как в GetTile вызывается ShiftRect, который работает только с RGBA

		id := "0"

		err := bmpService.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		_, err = bmpService.GetTile(id, x, y)
		So(err, ShouldBeNil)

		err = bmpService.DeleteImage(id)
		So(err, ShouldBeNil)

		_, err = bmpService.GetTile(id, x, y)
		So(err, ShouldNotBeNil)
	})
}

func TestBmpService_ErrorGet(t *testing.T) {
	Convey("При запросе не существующего тайла должна быть ошибка.", t, func() {
		tileRepo := &TestTileRepo{images: make(map[string][]byte)}
		bmpService := imagetile.NewBmpService(tileRepo)

		_, err := bmpService.GetTile("0", 0, 0)
		So(err, ShouldNotBeNil)
	})
}
