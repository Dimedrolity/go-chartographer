package imgstore_test

import (
	"errors"
	"go-chartographer/internal/imgstore"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestFileSystemTileRepo - интеграционные тесты, так как происходит взаимодействие с файловой системой.

func TestFileSystemTileRepo_SuccessGet(t *testing.T) {
	Convey("Проверка сохранения изображений на диске и получения изображения.", t, func() {
		tileRepo, err := imgstore.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		const (
			x = 1
			y = 1
		)

		const id = "0"
		img := []byte{1, 2, 3, 4, 5}
		err = tileRepo.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		tile, err := tileRepo.GetTile(id, x, y)
		So(err, ShouldBeNil)

		So(tile, ShouldResemble, img)
	})
}

func TestFileSystemTileRepo_SuccessDelete(t *testing.T) {
	Convey("После создания и удаления файла вызов фукнции получения должен вернуть ошибку.", t, func() {
		tileRepo, err := imgstore.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		const (
			x = 0
			y = 0
		)

		const id = "0"
		img := []byte{1, 2, 3, 4, 5}

		err = tileRepo.SaveTile(id, x, y, img)
		So(err, ShouldBeNil)

		_, err = tileRepo.GetTile(id, x, y)
		So(err, ShouldBeNil)

		err = tileRepo.DeleteImage(id)
		So(err, ShouldBeNil)

		_, err = tileRepo.GetTile(id, x, y)
		So(err, ShouldNotBeNil)
	})
}

func TestFileSystemTileRepo_ErrorGet(t *testing.T) {
	Convey("При запросе не существующего файла должна быть ошибка.", t, func() {
		tileRepo, err := imgstore.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		_, err = tileRepo.GetTile("0", 0, 0)
		So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)
	})
}

func TestFileSystemTileRepo_ErrorDelete(t *testing.T) {
	Convey("При удалении не сущствующего файла не должно быть ошибки", t, func() {
		tileRepo, err := imgstore.NewFileSystemTileRepo(t.TempDir())
		So(err, ShouldBeNil)

		err = tileRepo.DeleteImage("0")
		So(err, ShouldBeNil)
	})
}
