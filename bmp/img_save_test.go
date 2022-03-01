package bmp

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
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

// TestDecodeEncode - интеграционный тест библиотеки "golang.org/x/image/bmp" для того, чтобы повысить доверие.
// Библиотека должна кодировать/декодировать верно, должны совпадать все реальные байты с ожидаемыми.
// TODO не использовать библиотеку bmp напрямую, только через публичное апи моего сервиса
func TestDecodeEncode(t *testing.T) {
	Convey("Байты изображения после Decode и Encode должны совпадать с исходными", t, func() {
		initialBytes, _ := os.ReadFile("testdata/rgb.bmp")
		r := bytes.NewReader(initialBytes)

		img, _ := bmp.Decode(r)
		encodeBytes, _ := Encode(img)

		want := initialBytes
		got := encodeBytes
		So(got, ShouldResemble, want)
	})
}

func TestEncodeDecode_RectStartNotZero(t *testing.T) {
	Convey("Decode должен работать корректно с изображением, "+
		"у которого прямоугольник имеет начальные координаты, отличные от (0;0)", t, func() {
		const (
			x0 = 2
			y0 = 2
			x1 = 2
			y1 = 2
		)
		img := image.NewRGBA(image.Rect(x0, y0, x1, y1))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				img.Set(x, y, red)
			}
		}
		encodeBytes, _ := Encode(img)

		r := bytes.NewReader(encodeBytes)
		decodedImg, _ := bmp.Decode(r)

		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				So(decodedImg.At(x, y), ShouldResemble, img.At(x, y))
			}
		}
		So(decodedImg.Bounds().Min, ShouldResemble, image.Pt(0, 0))
	})
}

// нужна константная строка байтов, чтобы не зависеть от Ф.С.?
// TODO посмотреть тесты в библиотеках для image и bmp
