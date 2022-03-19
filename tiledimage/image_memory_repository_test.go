package tiledimage_test

import (
	"chartographer-go/tiledimage"
	"errors"

	"image"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO Протестировать на многопоточность.

// Методы тестируются в парах. По-другому никак

func TestFileSystemTileRepo_SuccessGet(t *testing.T) {
	Convey("После создания получение изображения должно вернуть исходное", t, func() {
		imageRepo := tiledimage.NewInMemoryImageRepo()

		const (
			id     = "0"
			width  = 1
			height = 1
		)
		tiledImg := &tiledimage.Image{
			Id:     id,
			Width:  width,
			Height: height,
			Tiles:  []image.Rectangle{},
		}

		imageRepo.Add(tiledImg)
		got, err := imageRepo.Get(id)

		So(err, ShouldBeNil)
		So(got, ShouldResemble, tiledImg)
	})
}

func TestFileSystemTileRepo_SuccessDelete(t *testing.T) {
	Convey("После создания и удаления получение изображения должно вернуть исходное", t, func() {
		imageRepo := tiledimage.NewInMemoryImageRepo()

		const (
			id     = "0"
			width  = 1
			height = 1
		)
		tiledImg := &tiledimage.Image{
			Id:     id,
			Width:  width,
			Height: height,
			Tiles:  []image.Rectangle{},
		}

		imageRepo.Add(tiledImg)
		err := imageRepo.Delete(id)
		So(err, ShouldBeNil)
		_, err = imageRepo.Get(id)
		So(err, ShouldNotBeNil)
	})
}

func TestFileSystemTileRepo_Delete(t *testing.T) {
	Convey("Удаление несуществующего изображения должно вернуть ошибку", t, func() {
		imageRepo := tiledimage.NewInMemoryImageRepo()

		err := imageRepo.Delete("0")

		So(errors.Is(err, tiledimage.ErrNotExist), ShouldBeTrue)
	})
}

func TestFileSystemTileRepo_Get(t *testing.T) {
	Convey("Получение несуществующего изображения должно вернуть ошибку", t, func() {
		imageRepo := tiledimage.NewInMemoryImageRepo()

		_, err := imageRepo.Get("0")

		So(errors.Is(err, tiledimage.ErrNotExist), ShouldBeTrue)
	})
}
