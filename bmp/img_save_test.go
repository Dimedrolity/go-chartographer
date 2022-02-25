package bmp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Выделить структуру ChartasService, содержит путь к каталогу данных и все методы
// Тогда нужно выделить интерфейс и методы сервиса? что даст?

// инициализировать сервис - путь к каталогу данных

func TestGuid(t *testing.T) {
	Convey("Guid", t, func() {
		id := Guid()
		So(id, ShouldHaveLength, 36)
	})
}

func TestAppendExt(t *testing.T) {
	Convey("AppendExtension", t, func() {
		filename := AppendExtension("qwerty")
		So(filename, ShouldEqual, "qwerty.bmp")
	})
}

// интеграционный тест для проверка, что изображение корректно сохраняется на диск
