package server_test

import (
	"bytes"
	"image"
	"image/color"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/Dimedrolity/go-chartographer/internal/chart"
	"github.com/Dimedrolity/go-chartographer/internal/server"
)

// Может быть подключить библиотеку для создания стабов в рантайме? типа FakeItEasy на C#
// Сейчас для каждого теста руками создана стаб-структура

// region Создание изображения

type TestChartService struct {
	chart.Service
}

func TestCreate_WrongSize(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartService{})

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

type TestChartServiceCreateMethodSizeErr struct {
	chart.Service
}

func (t TestChartServiceCreateMethodSizeErr) AddImage(int, int) (*chart.TiledImage, error) {
	return nil, &chart.SizeError{}
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

type TestChartServiceCreateMethodSuccess struct {
	chart.Service
}

const id = "new"

func (t TestChartServiceCreateMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	return &chart.TiledImage{
		Id: id,
	}, nil
}

func TestCreate_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceCreateMethodSuccess{})

	type Size struct {
		Width, Height interface{}
	}
	tmpl, _ := template.New("test").
		Parse("/chartas/?width={{.Width}}&height={{.Height}}")

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

// endregion

// region Удаление изображения

type TestChartServiceDeleteMethodNotFound struct {
	chart.Service
}

func (t TestChartServiceDeleteMethodNotFound) DeleteImage(string) error {
	return chart.ErrNotExist
}

func TestDelete_NotFound(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceDeleteMethodNotFound{})

	tmpl, _ := template.New("right request").
		Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

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

type TestChartServiceDeleteMethodSuccess struct {
	chart.Service
}

func (t TestChartServiceDeleteMethodSuccess) DeleteImage(string) error {
	return nil
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

// endregion

// region Получение фрагмента изображения

func TestGet_WrongParams(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartService{})

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

type TestChartServiceGetMethodNotFound struct {
	chart.Service
}

func (t TestChartServiceGetMethodNotFound) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return nil, nil
}
func (t TestChartServiceGetMethodNotFound) GetImage(string) (*chart.TiledImage, error) {
	return nil, chart.ErrNotExist
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

type TestChartServiceGetMethodSizeError struct {
	chart.Service
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

type TestChartServiceGetMethodNotOverlaps struct {
	chart.Service
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

type TestChartServiceGetMethodSuccess struct {
	chart.Service
}

func (t TestChartServiceGetMethodSuccess) GetFragment(*chart.TiledImage, int, int, int, int) (image.Image, error) {
	return image.Image(image.Rectangle{}), nil
}
func (t TestChartServiceGetMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}
func (t TestChartServiceGetMethodSuccess) Encode(image.Image) ([]byte, error) {
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

// endregion

// region Установка фрагмента изображения

type Fragment struct {
	Id                  string
	X, Y, Width, Height interface{}
}

type TestChartServiceSetMethodWrongSize struct {
	chart.Service
}

func (t TestChartServiceSetMethodWrongSize) SetFragment(*chart.TiledImage, int, int, image.Image) error {
	return nil
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

type TestChartServiceSetMethodNotFound struct {
	chart.Service
}

func (t TestChartServiceSetMethodNotFound) GetImage(string) (*chart.TiledImage, error) {
	return nil, chart.ErrNotExist
}

func TestSet_NotFound(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodNotFound{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		req := httptest.NewRequest("POST", url, &bytes.Buffer{})
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

type TestChartServiceSetMethodNotOverlaps struct {
	chart.Service
}

func (t TestChartServiceSetMethodNotOverlaps) SetFragment(*chart.TiledImage, int, int, image.Image) error {
	return chart.ErrNotOverlaps
}
func (t TestChartServiceSetMethodNotOverlaps) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}
func (t TestChartServiceSetMethodNotOverlaps) Decode([]byte) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{A: 0xFF})
	return img, nil
}

func TestSet_NotOverlaps(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodNotOverlaps{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()
		// пустой буфер, так как Decode вернет стаб
		req := httptest.NewRequest("POST", url, &bytes.Buffer{})
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

type TestChartServiceSetMethodSuccess struct {
	chart.Service
}

func (t TestChartServiceSetMethodSuccess) AddImage(int, int) (*chart.TiledImage, error) {
	return nil, nil
}
func (t TestChartServiceSetMethodSuccess) SetFragment(*chart.TiledImage, int, int, image.Image) error {
	return nil
}
func (t TestChartServiceSetMethodSuccess) GetImage(string) (*chart.TiledImage, error) {
	return nil, nil
}
func (t TestChartServiceSetMethodSuccess) Decode([]byte) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{A: 0xFF})
	return img, nil
}

func TestSet_Success(t *testing.T) {
	srv := server.NewServer(&server.Config{}, &TestChartServiceSetMethodSuccess{})

	tmpl, _ := template.New("right request").Parse("/chartas/{{.Id}}/?x={{.X}}&y={{.Y}}&width={{.Width}}&height={{.Height}}")

	Convey("", t, func() {
		b := bytes.Buffer{}
		err := tmpl.Execute(&b, &Fragment{Id: "0", X: 0, Y: 0, Width: 1, Height: 1})
		So(err, ShouldBeNil)
		url := b.String()

		// пустой буфер, так как Decode вернет стаб
		req := httptest.NewRequest("POST", url, &bytes.Buffer{})
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

// endregion
