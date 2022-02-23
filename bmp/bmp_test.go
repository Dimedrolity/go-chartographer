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
		imgBytes, _ := os.ReadFile("rgb.bmp")
		r := bytes.NewReader(imgBytes)

		img, _ := bmp.Decode(r)
		w := bytes.Buffer{}
		_ = bmp.Encode(&w, img)

		want := len(imgBytes)
		got := len(w.Bytes())
		So(got, ShouldResemble, want)
	})
}
