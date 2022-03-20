package server_test

import (
	"bytes"
	"chartographer-go/chart"
	"chartographer-go/kvstore"
	"chartographer-go/server"
	"golang.org/x/image/bmp"
	"image"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	. "github.com/smartystreets/goconvey/convey"
)

// Может быть подключить библиотеку для создания стабов в рантайме? типа FakeItEasy на C#
// Сейчас для каждого теста руками создана стаб-структура

// region Создание изображения

type TestChartServiceAllPanic struct{}

func (t TestChartServiceAllPanic) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceAllPanic) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceAllPanic) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceAllPanic) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceAllPanic) GetImage(string) (*chart.TiledImage, error) {
	panic("implement me")
}

func TestCreate_WrongSize(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceAllPanic{})

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

}

type TestChartServiceCreateMethodSizeErr struct{}

func (t TestChartServiceCreateMethodSizeErr) AddImage(int, int) (*chart.TiledImage, error) {
	return nil, &chart.SizeError{}
}
func (t TestChartServiceCreateMethodSizeErr) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceCreateMethodSizeErr) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}

func (t TestChartServiceCreateMethodSizeErr) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceCreateMethodSizeErr) GetImage(string) (*chart.TiledImage, error) {
	panic("implement me")
}

func TestCreate_SizeErr(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceCreateMethodSizeErr{})

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

	Convey("negative size", t, func() {
		// сервис-stub вернет SizeError
		testWrongSize(tmpl, &Size{Width: 0, Height: 0}) // 0, чтобы пройти проверки на числа
	})
}

type TestChartServiceCreateMethodSuccess struct{}

const id = "new"

func (t TestChartServiceCreateMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	return &chart.TiledImage{
		Id: id,
	}, nil
}
func (t TestChartServiceCreateMethodSuccess) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceCreateMethodSuccess) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceCreateMethodSuccess) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceCreateMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	panic("implement me")
}

func TestCreate_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceCreateMethodSuccess{})

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

//endregion

// region Удаление изображения

type TestChartServiceDeleteMethodNotFound struct{}

func (t TestChartServiceDeleteMethodNotFound) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceDeleteMethodNotFound) DeleteImage(string) error {
	return kvstore.ErrNotExist
}
func (t TestChartServiceDeleteMethodNotFound) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}

func (t TestChartServiceDeleteMethodNotFound) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceDeleteMethodNotFound) GetImage(string) (*chart.TiledImage, error) {
	panic("implement me")
}

func TestDelete_NotFound(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceDeleteMethodNotFound{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

type TestChartServiceDeleteMethodSuccess struct{}

func (t TestChartServiceDeleteMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceDeleteMethodSuccess) DeleteImage(string) error {
	return nil
}
func (t TestChartServiceDeleteMethodSuccess) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceDeleteMethodSuccess) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceDeleteMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	panic("implement me")
}

func TestDelete_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceDeleteMethodSuccess{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

//endregion

//region Получение фрагмента изображения

func TestGet_WrongParams(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceAllPanic{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	testWrongSize := func(tmpl *template.Template, frag *Fragment) {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, frag)
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	}

	Convey("no X", t, func() {
		tmplNoX, _ := template.New("no X").Parse("/chartas/{{.Id}}/?y={{.Y}}&width={{.Width}}&height={{.Height}}")
		testWrongSize(tmplNoX, &Fragment{Y: 0, Width: 1, Height: 1})
	})
	Convey("no Y", t, func() {
		tmplNoY, _ := template.New("no Y").Parse("/chartas/{{.Id}}/?x={{.X}}&&width={{.Width}}")
		testWrongSize(tmplNoY, &Fragment{X: 0, Width: 1, Height: 1})
	})
	Convey("no width", t, func() {
		tmplNoWidth, _ := template.New("no width").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&&height={{.Height}}")
		testWrongSize(tmplNoWidth, &Fragment{X: 0, Y: 0, Height: 1})
	})
	Convey("no height", t, func() {
		tmplNoHeight, _ := template.New("no height").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}")
		testWrongSize(tmplNoHeight, &Fragment{X: 0, Y: 0, Width: 1})
	})

	Convey("empty X", t, func() {
		testWrongSize(tmpl, &Fragment{X: "", Y: 0, Width: 1, Height: 1})
	})
	Convey("empty Y", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: "", Width: 1, Height: 1})
	})
	Convey("empty width", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: "", Height: 1})
	})
	Convey("empty height", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: 1, Height: ""})
	})

	Convey("string X", t, func() {
		testWrongSize(tmpl, &Fragment{X: "a", Y: 0, Width: 1, Height: 1})
	})
	Convey("string Y", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: "a", Width: 1, Height: 1})
	})
	Convey("string width", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: "a", Height: 1})
	})
	Convey("string height", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: 1, Height: "a"})
	})
}

type TestChartServiceGetMethodNotFound struct{}

func (t TestChartServiceGetMethodNotFound) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceGetMethodNotFound) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodNotFound) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}

func (t TestChartServiceGetMethodNotFound) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return nil, nil
}
func (t TestChartServiceGetMethodNotFound) GetImage(string) (*chart.TiledImage, error) {
	return nil, kvstore.ErrNotExist
}

func TestGet_NotFound(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceGetMethodNotFound{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

type TestChartServiceGetMethodSizeError struct{}

func (t TestChartServiceGetMethodSizeError) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")

}
func (t TestChartServiceGetMethodSizeError) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodSizeError) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}

func (t TestChartServiceGetMethodSizeError) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return nil, &chart.SizeError{}
}
func (t TestChartServiceGetMethodSizeError) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestGet_SizeErr(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceGetMethodSizeError{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

type TestChartServiceGetMethodNotOverlaps struct{}

func (t TestChartServiceGetMethodNotOverlaps) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceGetMethodNotOverlaps) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodNotOverlaps) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodNotOverlaps) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return nil, chart.ErrNotOverlaps
}
func (t TestChartServiceGetMethodNotOverlaps) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestGet_NotOverlaps(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceGetMethodNotOverlaps{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

type TestChartServiceGetMethodSuccess struct{}

func (t TestChartServiceGetMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceGetMethodSuccess) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodSuccess) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceGetMethodSuccess) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return image.Image(image.Rectangle{}), nil
}
func (t TestChartServiceGetMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestGet_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceGetMethodSuccess{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

//endregion

//region Установка фрагмента изображения

type Fragment struct {
	Id                  string
	X, Y, Width, Height interface{}
}

type TestChartServiceSetMethodWrongSize struct{}

func (t TestChartServiceSetMethodWrongSize) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodWrongSize) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceSetMethodWrongSize) SetFragment(*chart.TiledImage, image.Image) error {
	return nil
}
func (t TestChartServiceSetMethodWrongSize) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodWrongSize) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestSet_WrongSize(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodWrongSize{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	testWrongSize := func(tmpl *template.Template, frag *Fragment) {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, frag)
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("POST", url, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	}

	Convey("no X", t, func() {
		tmplNoX, _ := template.New("no X").Parse("/chartas/{{.Id}}/?y={{.Y}}&width={{.Width}}&height={{.Height}}")
		testWrongSize(tmplNoX, &Fragment{Y: 0, Width: 1, Height: 1})
	})
	Convey("no Y", t, func() {
		tmplNoY, _ := template.New("no Y").Parse("/chartas/{{.Id}}/?x={{.X}}&&width={{.Width}}")
		testWrongSize(tmplNoY, &Fragment{X: 0, Width: 1, Height: 1})
	})
	Convey("no width", t, func() {
		tmplNoWidth, _ := template.New("no width").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&&height={{.Height}}")
		testWrongSize(tmplNoWidth, &Fragment{X: 0, Y: 0, Height: 1})
	})
	Convey("no height", t, func() {
		tmplNoHeight, _ := template.New("no height").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}")
		testWrongSize(tmplNoHeight, &Fragment{X: 0, Y: 0, Width: 1})
	})

	Convey("empty X", t, func() {
		testWrongSize(tmpl, &Fragment{X: "", Y: 0, Width: 1, Height: 1})
	})
	Convey("empty Y", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: "", Width: 1, Height: 1})
	})
	Convey("empty width", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: "", Height: 1})
	})
	Convey("empty height", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: 1, Height: ""})
	})

	Convey("string X", t, func() {
		testWrongSize(tmpl, &Fragment{X: "a", Y: 0, Width: 1, Height: 1})
	})
	Convey("string Y", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: "a", Width: 1, Height: 1})
	})
	Convey("string width", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: "a", Height: 1})
	})
	Convey("string height", t, func() {
		testWrongSize(tmpl, &Fragment{X: 0, Y: 0, Width: 1, Height: "a"})
	})
}

type TestChartServiceSetMethodNotFound struct{}

func (t TestChartServiceSetMethodNotFound) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotFound) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotFound) SetFragment(*chart.TiledImage, image.Image) error {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotFound) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotFound) GetImage(string) (*chart.TiledImage, error) {
	return nil, kvstore.ErrNotExist
}

func TestSet_NotFound(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodNotFound{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		img := image.Image(image.Rect(0, 0, 1, 1))
		imgBuffer := bytes.Buffer{}
		_ = bmp.Encode(&imgBuffer, img)
		req := httptest.NewRequest("POST", url, &imgBuffer)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

type TestChartServiceSetMethodNotOverlaps struct{}

func (t TestChartServiceSetMethodNotOverlaps) AddImage(int, int) (*chart.TiledImage, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotOverlaps) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotOverlaps) SetFragment(*chart.TiledImage, image.Image) error {
	return chart.ErrNotOverlaps
}
func (t TestChartServiceSetMethodNotOverlaps) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodNotOverlaps) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestSet_NotOverlaps(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodNotOverlaps{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		img := image.Image(image.Rect(0, 0, 1, 1))
		imgBuffer := bytes.Buffer{}
		_ = bmp.Encode(&imgBuffer, img)
		req := httptest.NewRequest("POST", url, &imgBuffer)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

type TestChartServiceSetMethodSuccess struct{}

func (t TestChartServiceSetMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	return nil, nil
}
func (t TestChartServiceSetMethodSuccess) DeleteImage(string) error {
	panic("implement me")
}
func (t TestChartServiceSetMethodSuccess) SetFragment(*chart.TiledImage, image.Image) error {
	return nil
}
func (t TestChartServiceSetMethodSuccess) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	panic("implement me")
}
func (t TestChartServiceSetMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}

func TestSet_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodSuccess{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		img := image.Image(image.Rect(0, 0, 1, 1))
		imgBuffer := bytes.Buffer{}
		_ = bmp.Encode(&imgBuffer, img)
		req := httptest.NewRequest("POST", url, &imgBuffer)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

//endregion
