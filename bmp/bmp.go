package bmp

import "encoding/binary"

func createBmp(width, height int) []byte {
	// можно использовать массив, а не слайс.
	// Тогда нужно будет хранить offset и увеличивать его

	image := make([]byte, 0)

	// BMP Header
	const bmpHeaderSize = 14

	format := []byte("BM")
	image = append(image, format...)

	sizeBytes := make([]byte, 4)
	pixelsDataSize := width * height * 4
	const DIBHeaderSize = 40
	size := bmpHeaderSize + DIBHeaderSize + pixelsDataSize
	binary.LittleEndian.PutUint32(sizeBytes, uint32(size))
	image = append(image, sizeBytes...)

	unused := make([]byte, 2)
	image = append(image, unused...)
	image = append(image, unused...)

	offsetBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(offsetBytes, uint32(bmpHeaderSize+DIBHeaderSize))
	image = append(image, offsetBytes...)

	// DIB Header
	DIBSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(DIBSizeBytes, uint32(DIBHeaderSize))
	image = append(image, DIBSizeBytes...)

	widthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(widthBytes, uint32(width))
	image = append(image, widthBytes...)

	heightBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(heightBytes, uint32(height))
	image = append(image, heightBytes...)

	planeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(planeBytes, uint16(1))
	image = append(image, planeBytes...)

	const bitsInPixel = 24
	bitsPerPixelBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bitsPerPixelBytes, uint16(bitsInPixel))
	image = append(image, bitsPerPixelBytes...)

	compressionDummy := make([]byte, 4)
	image = append(image, compressionDummy...)

	imgSizeDummy := make([]byte, 4)
	image = append(image, imgSizeDummy...)

	horResBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(horResBytes, uint32(2835))
	image = append(image, horResBytes...)

	vertResBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(vertResBytes, uint32(2835))
	image = append(image, vertResBytes...)

	colorsBytes := make([]byte, 4)
	image = append(image, colorsBytes...)

	impColorsDummy := make([]byte, 4)
	image = append(image, impColorsDummy...)

	// Pixel array

	// Цветное
	//red := []byte{0, 0, 255}
	//green := []byte{0, 255, 0}
	//blue := []byte{255, 0, 0}
	//white := []byte{255, 255, 255}
	//padding := 4 - (width * 3 % 4)
	//paddingBytes := make([]byte, padding)
	//
	//image = append(image, red...)
	//image = append(image, white...)
	//image = append(image, paddingBytes...)
	//
	//image = append(image, blue...)
	//image = append(image, green...)
	//image = append(image, paddingBytes...)

	//Черное
	black := []byte{0, 0, 0}

	paddingBytes := make([]byte, padding(width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			image = append(image, black...)
		}
		image = append(image, paddingBytes...)
	}

	return image
}

// padding расчитывает количество байт, которые необходимо добавлять в конец каждой строки.
// В соотв. с форматом BMP, каждая строка пикселей должна содержать количество байт, кратное 4.
// Пример: ширина = 2; байтов пикселей будет 2*3=6; необходим padding, равный 2-ум байтам.
func padding(width int) int {
	const bytesPerPixel = 3
	return 4 - (width * bytesPerPixel % 4)
}
