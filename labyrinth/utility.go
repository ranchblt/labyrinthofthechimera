package labyrinth

import (
	"bytes"
	"image"

	_ "image/png" // needed for png images

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/labyrinthofthechimera/resource"
)

// openImage gets the image out of go-bindata
func openImage(path string) (image.Image, error) {
	b, err := resource.Asset(path)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	return image, nil
}

// handleErr panics on any error, makes error handling cleaner
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// utility function to do DST for a mobile object with no unusual modifications to the return values
func defaultMobileDST(i int, center *coord, image *ebiten.Image) (x0, y0, x1, y1 int) {
	width, height := image.Size()
	height = height / 2
	width = width / 2
	return center.x - width,
		center.y - height,
		center.x + width,
		center.y + height
}

// utility function to do DST for a stationary object with no unusual modifications to the return values
func defaultStationaryDST(i int, center *coord, image *ebiten.Image) (x0, y0, x1, y1 int) {
	width, height := image.Size()
	return center.x,
		center.y,
		center.x + width,
		center.y + height
}
