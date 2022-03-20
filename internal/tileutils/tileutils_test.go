package tileutils_test

import (
	"image"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "go-chartographer/internal/tileutils"
)

func TestCreateTiles(t *testing.T) {
	Convey("Проверка того, что центральные тайлы будут максимального размера, а крайние не максимального", t, func() {
		const (
			width       = 25
			height      = 25
			maxTileSize = 10
		)

		tiles := CreateTiles(width, height, maxTileSize)

		So(tiles, ShouldHaveLength, 9)
		// 1-я строка
		So(tiles[0].Min, ShouldResemble, image.Pt(0, 0))
		So(tiles[0].Max, ShouldResemble, image.Pt(10, 10))

		So(tiles[1].Min, ShouldResemble, image.Pt(10, 0))
		So(tiles[1].Max, ShouldResemble, image.Pt(20, 10))

		So(tiles[2].Min, ShouldResemble, image.Pt(20, 0))
		So(tiles[2].Max, ShouldResemble, image.Pt(25, 10))

		// 2-я строка
		So(tiles[3].Min, ShouldResemble, image.Pt(0, 10))
		So(tiles[3].Max, ShouldResemble, image.Pt(10, 20))

		So(tiles[4].Min, ShouldResemble, image.Pt(10, 10))
		So(tiles[4].Max, ShouldResemble, image.Pt(20, 20))

		So(tiles[5].Min, ShouldResemble, image.Pt(20, 10))
		So(tiles[5].Max, ShouldResemble, image.Pt(25, 20))

		// 3-я строка
		So(tiles[6].Min, ShouldResemble, image.Pt(0, 20))
		So(tiles[6].Max, ShouldResemble, image.Pt(10, 25))

		So(tiles[7].Min, ShouldResemble, image.Pt(10, 20))
		So(tiles[7].Max, ShouldResemble, image.Pt(20, 25))

		So(tiles[8].Min, ShouldResemble, image.Pt(20, 20))
		So(tiles[8].Max, ShouldResemble, image.Pt(25, 25))
	})
}

func TestOverlappedTiles(t *testing.T) {
	Convey("", t, func() {
		fragment := image.Rect(5, 5, 10, 10)
		tiles := []image.Rectangle{
			image.Rect(0, 0, 10, 10),   // пересекается
			image.Rect(50, 50, 10, 10), // не пересекается
		}

		overlapped := OverlappedTiles(tiles, fragment)

		So(overlapped, ShouldHaveLength, 1)
		So(overlapped[0], ShouldResemble, tiles[0])
	})
}
