package main

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"go-chartographer/internal/imagetile"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"testing"
)

// TestEncodeDecode - интеграционный тест библиотеки "golang.org/x/image/bmp" для того, чтобы повысить доверие.
// Библиотека должна кодировать/декодировать верно, должны совпадать все реальные байты с ожидаемыми.
func TestEncodeDecode(t *testing.T) {
	Convey("Байты изображения после Encode и Decode должны совпадать с исходными", t, func() {
		// Assign
		img := image.NewRGBA(image.Rect(0, 0, 2, 2))

		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		const (
			redX = 0
			redY = 0
		)
		img.Set(redX, redY, red)

		green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		const (
			greenX = 1
			greenY = 0
		)
		img.Set(greenX, greenY, green)

		blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		const (
			blueX = 0
			blueY = 1
		)
		img.Set(blueX, blueY, blue)

		white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		const (
			whiteX = 1
			whiteY = 1
		)
		img.Set(whiteX, whiteY, white)

		// Act
		buffer := bytes.Buffer{}
		err := bmp.Encode(&buffer, img)
		So(err, ShouldBeNil)

		decodedImg, err := bmp.Decode(&buffer)
		So(err, ShouldBeNil)

		// Assert
		So(img.Bounds(), ShouldResemble, decodedImg.Bounds())
		So(decodedImg.At(redX, redY), ShouldResemble, red)
		So(decodedImg.At(greenX, greenY), ShouldResemble, green)
		So(decodedImg.At(blueX, blueY), ShouldResemble, blue)
		So(decodedImg.At(whiteX, whiteY), ShouldResemble, white)
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
		encodeBytes, _ := imagetile.Encode(img)

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
