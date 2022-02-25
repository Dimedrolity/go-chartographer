package bmp

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/bmp"
	"os"
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

// TestBMPLib - интеграционный тест библиотеки "golang.org/x/image/bmp" для того, чтобы повысить доверие.
// Библиотека должна кодировать/декодировать верно,
// должны совпадать реальные с ожидаемыми: размер файла, размер изображений, пиксели - все байты.
// TODO не использовать библиотеку bmp напрямую, только через публичное апи моего сервиса
func TestBMPLib(t *testing.T) {
	Convey("Testing `bmp lib` function", t, func() {
		initialBytes, _ := os.ReadFile("rgb.bmp")
		r := bytes.NewReader(initialBytes)

		img, _ := bmp.Decode(r)
		encodeBytes, _ := Encode(&img)

		want := initialBytes
		got := encodeBytes
		So(got, ShouldResemble, want)
	})
}

// нужна константная строка байтов, чтобы не зависеть от Ф.С.?
