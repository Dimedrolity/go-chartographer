package kvstore_test

import (
	"errors"
	"go-chartographer/pkg/kvstore"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Методы тестируются в парах, так как нет доступа к внутренней структуре данных

func TestFileSystemTileRepo_SuccessGet(t *testing.T) {
	Convey("После создания получение изображения должно вернуть исходное", t, func() {
		imageRepo := kvstore.NewInMemoryStore()

		const (
			key   = "0"
			value = 42
		)

		imageRepo.Add(key, value)
		got, err := imageRepo.Get(key)

		So(err, ShouldBeNil)
		So(got, ShouldResemble, value)
	})
}

func TestFileSystemTileRepo_SuccessDelete(t *testing.T) {
	Convey("После создания и удаления получение изображения должно вернуть исходное", t, func() {
		imageRepo := kvstore.NewInMemoryStore()

		const (
			key   = "0"
			value = 42
		)

		imageRepo.Add(key, value)
		err := imageRepo.Delete(key)
		So(err, ShouldBeNil)
		_, err = imageRepo.Get(key)
		So(err, ShouldNotBeNil)
	})
}

func TestFileSystemTileRepo_Delete(t *testing.T) {
	Convey("Удаление несуществующего изображения должно вернуть ошибку", t, func() {
		imageRepo := kvstore.NewInMemoryStore()

		err := imageRepo.Delete("0")

		So(errors.Is(err, kvstore.ErrNotExist), ShouldBeTrue)
	})
}

func TestFileSystemTileRepo_Get(t *testing.T) {
	Convey("Получение несуществующего изображения должно вернуть ошибку", t, func() {
		imageRepo := kvstore.NewInMemoryStore()

		_, err := imageRepo.Get("0")

		So(errors.Is(err, kvstore.ErrNotExist), ShouldBeTrue)
	})
}
