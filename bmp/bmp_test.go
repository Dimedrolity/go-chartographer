package bmp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEmptyImg(t *testing.T) {
	Convey("Testing `createBmp` function", t, func() {
		imgBytes := createBmp(2, 3)
		So(len(imgBytes), ShouldEqual, 78)
	})

	//os.WriteFile("img.bmp", img, 0777)
}
