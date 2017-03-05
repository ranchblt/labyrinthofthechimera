package labyrinth

import (
	"image"

	"github.com/ranchblt/labyrinthofthechimera/collision"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/draw"
)

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
func defaultDST(i int, topLeft collision.Coord, image *ebiten.Image) (x0, y0, x1, y1 int) {
	width, height := image.Size()
	return topLeft.X(),
		topLeft.Y(),
		topLeft.X() + width,
		topLeft.Y() + height
}

func scaledDST(i int, topLeft collision.Coord, image *ebiten.Image, scale float64) (x0, y0, x1, y1 int) {
	width, height := image.Size()
	return topLeft.X(),
		topLeft.Y(),
		topLeft.X() + int(float64(width)*scale),
		topLeft.Y() + int(float64(height)*scale)
}

// center finds the center coords based on the topLeft and the size of the image
func center(topLeft collision.Coord, image *ebiten.Image) *coord {
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

func containsPowerupClass(powerups []*powerup, class powerupClass) bool {
	// eat it, multiple returns.
	for _, p := range powerups {
		if p.class == class {
			return true
		}
	}
	return false
}
