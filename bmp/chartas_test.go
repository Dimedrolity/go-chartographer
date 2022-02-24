package bmp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateImg_ZeroWidth(t *testing.T) {
	Convey("ZeroWidth", t, func() {
		_, err := CreateImg(0, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_ZeroHeight(t *testing.T) {
	Convey("ZeroHeight", t, func() {
		_, err := CreateImg(1, 0)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_WidthExceeded(t *testing.T) {
	Convey("WidthExceeded", t, func() {
		_, err := CreateImg(maxWidth+1, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_HeightExceeded(t *testing.T) {
	Convey("HeightExceeded", t, func() {
		_, err := CreateImg(1, maxHeight+1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_MinMaxSize(t *testing.T) {
	testSize := func(width, height int) {
		img, err := CreateImg(width, height)
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
}

// Выделить структуру ChartasService, содержит путь к каталогу данных и все методы
// Тогда нужно выделить интерфейс и методы сервиса? что даст?

// инициализировать сервис - путь к каталогу данных

// TestSaveImg - интеграционный тест для проверка, что изображение корректно сохраняется на диск
//func TestSaveImg(t *testing.T) {
//	Convey("SaveImg", t, func() {
//		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
//		id, err := SaveImg(img)
//		So(err, ShouldNotBeNil)
//		So(id, ShouldHaveLength, 36)
//		// отдельная папка для тестовых изображений?
//		So(got, ShouldResemble, want)
//
//		// удалить это изображение.
//
//	})
//}

// сохранять нужно и по текущему id. не обязательно создание нового id.

// удалить файл и создать новый.
//ЛИБО переписать полностью текущий файл.
//ЛИБО переписать частично текущий файл

//func TestCreateFile(t *testing.T) {
//	_ = os.WriteFile("img/a.bmp", []byte{0, 1, 0}, 0777)
//}
