package bmp

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/bmp"
	"os"
	"testing"
)

func TestEmptyImg(t *testing.T) {
	Convey("Testing `createBmp` function", t, func() {
		imgBytes := createBmp(2, 3)
		So(len(imgBytes), ShouldEqual, 78)
	})
}

// TestBMPLib - интеграционный тест библиотеки "golang.org/x/image/bmp" для того, чтобы повысить доверие.
// Библиотека должна кодировать/декодировать верно,
// должны совпадать реальные с ожидаемыми: размер файла, размер изображений, пиксели - все байты.
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
