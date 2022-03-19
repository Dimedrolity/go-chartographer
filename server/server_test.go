package server_test

import (
	"bytes"
	"chartographer-go/chart"
	"chartographer-go/server"
	"chartographer-go/tiledimage"
	"image"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	. "github.com/smartystreets/goconvey/convey"
)

//
// Создание изображения
//

type TestChartServiceSizeErr struct{}

func (t TestChartServiceSizeErr) NewRgbaBmp(int, int) (*tiledimage.Image, error) {
	return nil, &chart.SizeError{}
}
func (t TestChartServiceSizeErr) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceSizeErr) SetFragment(string, image.Image) error {
	panic("implement me")
}

func (t TestChartServiceSizeErr) GetFragment(*tiledimage.Image, int, int, int, int) (image.Image, error) {
	return nil, &chart.SizeError{}
}
func (t TestChartServiceSizeErr) GetTiledImage(string) (*tiledimage.Image, error) {
	panic("implement me")
}

func TestCreate_WrongSize(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSizeErr{})

	type Size struct {
		Width, Height interface{}
	}
	tmpl, _ := template.New("right request").Parse("/chartas/?width={{.Width}}&height={{.Height}}")

	testWrongSize := func(tmpl *template.Template, size *Size) {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, size)
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("POST", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	}

	Convey("no width", t, func() {
		tmplOnlyHeight, _ := template.New("no width").Parse("/chartas/?height={{.Height}}")
		testWrongSize(tmplOnlyHeight, &Size{Height: 1})
	})
	Convey("no height", t, func() {
		tmplOnlyWidth, _ := template.New("no height").Parse("/chartas/?width={{.Width}}")
		testWrongSize(tmplOnlyWidth, &Size{Width: 1})
	})

	Convey("empty width", t, func() {
		testWrongSize(tmpl, &Size{Width: "", Height: 1})
	})
	Convey("empty height", t, func() {
		testWrongSize(tmpl, &Size{Width: 1, Height: ""})
	})

	Convey("string width", t, func() {
		testWrongSize(tmpl, &Size{Width: "a", Height: 1})
	})
	Convey("string height", t, func() {
		testWrongSize(tmpl, &Size{Width: 1, Height: "a"})
	})

	Convey("negative width", t, func() {
		// сервис вернет SizeError
		testWrongSize(tmpl, &Size{Width: 0, Height: 0}) // 0, чтобы пройти проверки на числа
	})
}

type TestChartService struct{}

const id = "new"

func (t TestChartService) NewRgbaBmp(int, int) (*tiledimage.Image, error) {
	return &tiledimage.Image{
		Id: id,
	}, nil
}
func (t TestChartService) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartService) SetFragment(string, image.Image) error {
	panic("implement me")
}

func (t TestChartService) GetFragment(*tiledimage.Image, int, int, int, int) (image.Image, error) {
	return nil, nil
}
func (t TestChartService) GetTiledImage(string) (*tiledimage.Image, error) {
	panic("implement me")
}

func TestCreate_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartService{})

	type Size struct {
		Width, Height interface{}
	}
	tmpl, _ := template.New("test").Parse("/chartas/?width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Size{Width: 0, Height: 0}) // 0, чтобы пройти проверки strconv.Atoi
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("POST", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusCreated)
		So(w.Body.String(), ShouldEqual, id)
	})
}

//
// TODO Удаление изображения
//

//
// TODO Получение фрагмента изображения
//

//
// TODO Установка фрагмента изображения
//
