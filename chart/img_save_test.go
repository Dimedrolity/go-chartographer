package chart

import (
	"bytes"
	"chartographer-go/store"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
	"testing"
)

// Выделить структуру ChartService, содержит путь к каталогу данных и все методы
// Тогда нужно выделить интерфейс и методы сервиса? что даст?

// TestDecodeEncode - интеграционный тест библиотеки "golang.org/x/image/bmp" для того, чтобы повысить доверие.
// Библиотека должна кодировать/декодировать верно, должны совпадать все реальные байты с ожидаемыми.
func TestDecodeEncode(t *testing.T) {
	Convey("Байты изображения после Decode и Encode должны совпадать с исходными", t, func() {
		const path = "testdata/rgb.bmp"

		file, err := os.Open(path)
		So(err, ShouldBeNil)

		img, err := bmp.Decode(file)
		So(err, ShouldBeNil)
		err = file.Close()
		So(err, ShouldBeNil)

		buffer := bytes.Buffer{}
		err = bmp.Encode(&buffer, img)
		So(err, ShouldBeNil)
		encodeBytes := buffer.Bytes()

		initialBytes, err := os.ReadFile(path)
		So(err, ShouldBeNil)
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
		encodeBytes, _ := store.Encode(img)

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
