package labyrinth

import (
	"bytes"
	"image"

	_ "image/png" // needed for png images

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
