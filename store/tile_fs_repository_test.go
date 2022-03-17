package store_test

import (
	"chartographer-go/store"
	"image"
	"image/color"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileSystemTileRepo(t *testing.T) {
	Convey("Проверка того, что центральные тайлы будут максимального размера, а крайние не максимального", t, func() {
		tileRepo, err := store.NewFileSystemTileRepo("testdata")
		So(err, ShouldBeNil)

		id := "0"
		const (
			x      = 1
			y      = 1
			width  = 1
			height = 1
		)
		img := image.NewRGBA(image.Rect(x, y, x+width, y+height))
		img.Set(x, y, color.RGBA{A: 255}) // чтобы в SaveTile Decode сохранил как 24-битное

		err = tileRepo.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		tile, err := tileRepo.GetTile(id, x, y)
		So(err, ShouldBeNil)

		So(tile.Bounds(), ShouldResemble, img.Bounds())
	})
}
