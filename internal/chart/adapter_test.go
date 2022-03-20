package chart_test

import (
	"go-chartographer/internal/chart"
	"image"
	"image/color"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdapter(t *testing.T) {
	Convey("Проверка смещения прямоугольника изображения."+
		"После смещения должны поменятся начальные и конечные координаты, "+
		"но размер изображения и пиксели не должны измениться", t, func() {
		const (
			x      = 0
			y      = 0
			width  = 2
			height = 2
		)
		img := image.NewRGBA(image.Rect(x, y, width, height))
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		for x := x; x < width; x++ {
			for y := y; y < height; y++ {
				img.SetRGBA(x, y, red)
			}
		}

		const (
			offsetX = 10
			offsetY = 20
		)
		a := chart.ImageAdapter{}
		a.ShiftRect(img, offsetX, offsetY)

		So(img.Bounds().Dx(), ShouldEqual, width)
		So(img.Bounds().Dy(), ShouldEqual, height)

		So(img.Bounds().Min, ShouldResemble, image.Pt(x+offsetX, y+offsetY))
		So(img.Bounds().Max, ShouldResemble, image.Pt(x+width+offsetX, y+height+offsetY))

		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				So(img.At(x, y), ShouldResemble, red)
			}
		}
	})
}
