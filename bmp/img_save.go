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
