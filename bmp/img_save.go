package bmp

import (
	"bytes"
	"github.com/google/uuid"
	"golang.org/x/image/bmp"
	"image"
)

func Encode(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := bmp.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// TODO принимать r io.Reader, чтобы не создавать свой
func Decode(imgBytes []byte) (image.Image, error) {
	r := bytes.NewReader(imgBytes)
	img, err := bmp.Decode(r)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func AppendExtension(filename string) string {
	return filename + ".bmp"
}

func Guid() string {
	uuidWithHyphen := uuid.New()
	return uuidWithHyphen.String()
}
