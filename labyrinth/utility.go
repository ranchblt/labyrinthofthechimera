package labyrinth

import (
	"bytes"
	"image"

	"image/gif"
	_ "image/png" // needed for png images

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/labyrinthofthechimera/resource"
	"golang.org/x/image/draw"
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

// openGif gets the gif out of go-bindata
func openGif(path string) (*gif.GIF, error) {
	b, err := resource.Asset(path)
	if err != nil {
		return nil, err
	}

	image, err := gif.DecodeAll(bytes.NewReader(b))

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

// helper to convert image to RGBA
func toRGBA(img image.Image) *image.RGBA {
	switch img.(type) {
	case *image.RGBA:
		return img.(*image.RGBA)
	}
	out := image.NewRGBA(img.Bounds())
	draw.Copy(out, image.Pt(0, 0), img, img.Bounds(), draw.Src, nil)
	return out
}

// utility function to do DST for a stationary object with no unusual modifications to the return values
func defaultDST(i int, topLeft *coord, image *ebiten.Image) (x0, y0, x1, y1 int) {
	width, height := image.Size()
	return topLeft.X(),
		topLeft.Y(),
		topLeft.X() + width,
		topLeft.Y() + height
}

// center finds the center coords based on the topLeft and the size of the image
func center(topLeft *coord, image *ebiten.Image) *coord {
	w, h := image.Size()
	return &coord{
		x: topLeft.X() + w/2,
		y: topLeft.Y() + h/2,
	}
}

func getPoweredUpValue(startValue int, powerups []*powerup, class powerupClass) int {
	poweredUpValue := startValue

	for _, p := range powerups {
		if p.class == class {
			poweredUpValue += p.boost
		}
	}

	return poweredUpValue
}
