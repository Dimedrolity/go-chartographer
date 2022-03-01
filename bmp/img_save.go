package bmp

import (
	"bytes"
	"github.com/google/uuid"
	"golang.org/x/image/bmp"
	"image"
)

// Encode
//если еще сделать обертку над Decode, то интеграционный тест нужно будет переписать.
func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func AppendExtension(filename string) string {
	return filename + ".bmp"
}

func Guid() string {
	uuidWithHyphen := uuid.New()
	return uuidWithHyphen.String()
}

//const pathToFolder = "testdata/"

//
//func SetPath(path string) {
//	pathToFolder = path
//}

// по сути такой код должен быть в HTTP-API или в интеграционном тесте
//func SaveImage(name string, img *image.RGBA) error {
//	imgBytes, err := Encode(img)
//	if err != nil {
//		return err
//	}
//
//	err = os.WriteFile(pathToFolder+name, imgBytes, 0777)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// по сути такой код должен быть в HTTP-API или в интеграционном тесте
//func CreateAndSaveImage(img *image.RGBA) (string, error) {
//	id := Guid()
//	name := AppendExtension(id)
//
//	err := SaveImage(name, img)
//	if err != nil {
//		return "", err
//	}
//
//	return id, nil
//}
