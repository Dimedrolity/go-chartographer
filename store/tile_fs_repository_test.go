package store_test

import (
	"chartographer-go/store"
	"errors"
	"image"
	"image/color"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestFileSystemTileRepo - интеграционные тесты, так как происходит взаимодействие с файловой системой.

func TestFileSystemTileRepo_SaveGet(t *testing.T) {
	Convey("Проверка сохранения изображений на диске и получения изображения."+
		"Проверяются также координаты изображения, должны соответствовать исходынм.", t, func() {
		tileRepo, err := store.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		const (
			x      = 1
			y      = 1
			width  = 1
			height = 1
		)
		img := image.NewRGBA(image.Rect(x, y, x+width, y+height))
		img.Set(x, y, color.RGBA{A: 255}) // чтобы в SaveTile Decode сохранил как 24-битное

		id := "0"
		err = tileRepo.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		tile, err := tileRepo.GetTile(id, x, y)
		So(err, ShouldBeNil)

		So(tile.Bounds(), ShouldResemble, img.Bounds())
	})
}

func TestFileSystemTileRepo_SaveDeleteGet(t *testing.T) {
	Convey("После создания и удаления файла вызов фукнции получения должен вернуть ошибку.", t, func() {
		tileRepo, err := store.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		const (
			x      = 0
			y      = 0
			width  = 1
			height = 1
		)
		img := image.NewRGBA(image.Rect(x, y, x+width, y+height))

		id := "0"
		err = tileRepo.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		err = tileRepo.DeleteImage(id)
		So(err, ShouldBeNil)

		_, err = tileRepo.GetTile(id, x, y)
		So(err, ShouldNotBeNil)
	})
}

func TestFileSystemTileRepo_Get(t *testing.T) {
	Convey("При запросе не существующего файла должна быть ошибка.", t, func() {
		tileRepo, err := store.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		_, err = tileRepo.GetTile("0", 0, 0)
		So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)
	})
}

func TestFileSystemTileRepo_Delete(t *testing.T) {
	Convey("При удалении не сущствующего файла не должно быть ошибки", t, func() {
		tileRepo, err := store.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)
		id := "0"

		err = tileRepo.DeleteImage(id)
		So(err, ShouldBeNil)
	})
}
