package bmp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateImg_ZeroWidth(t *testing.T) {
	Convey("ZeroWidth", t, func() {
		_, err := NewImage(0, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_ZeroHeight(t *testing.T) {
	Convey("ZeroHeight", t, func() {
		_, err := NewImage(1, 0)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_WidthExceeded(t *testing.T) {
	Convey("WidthExceeded", t, func() {
		_, err := NewImage(maxWidth+1, 1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_HeightExceeded(t *testing.T) {
	Convey("HeightExceeded", t, func() {
		_, err := NewImage(1, maxHeight+1)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateImg_MinMaxSize(t *testing.T) {
	testSize := func(width, height int) {
		img, err := NewImage(width, height)
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
