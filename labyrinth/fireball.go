package labyrinth

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten"
)

type fireballClass string

const (
	normalFireball = fireballClass("normal")
)

type fireball struct {
	image     *ebiten.Image
	rgba      *image.RGBA
	moveSpeed int
	topLeft   *coord
	class     fireballClass
	active    bool
}

func (f *fireball) Update() error {
	if f.offScreen() {
		f.active = false
		return nil
	}

	switch {
	case f.class == normalFireball:
		if err := f.updateNormalFireball(); err != nil {
			return err
		}
	default:
		return errors.New("Invalid class of fireball")
	}
	return nil
}

func (f *fireball) Draw(r *ebiten.Image) error {
	if f.active {
		r.DrawImage(f.image, &ebiten.DrawImageOptions{
			ImageParts: f,
		})
	}

	return nil
}

func (f *fireball) Len() int {
	return 1
}

func (f *fireball) updateNormalFireball() error {
	f.topLeft.x += f.moveSpeed
	return nil
}

func (f *fireball) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultStationaryDST(i, f.topLeft, f.image)
}

func (f *fireball) Src(i int) (x0, y0, x1, y1 int) {
	width, height := f.image.Size()
	return 0, 0, width, height
}

// offscreen checks if the left most part of the image is past the ScreenWidth
func (f *fireball) offScreen() bool {
	return f.topLeft.X() > ScreenWidth
}

func (f *fireball) RGBAImage() *image.RGBA {
	if f.rgba == nil {
		f.rgba = toRGBA(f.image)
	}
	return f.rgba
}

func (f *fireball) hit() {
	// this should probably trigger some kind of animation. could also
	// potentially have powerups that spawn more shots or something. clusterbombs!
	f.active = false
}

type fireballCreator struct {
	image     *ebiten.Image
	moveSpeed int
}

func (f *fireballCreator) newFireball(c coord, class fireballClass) *fireball {
	return &fireball{
		image:     f.image,
		topLeft:   &c,
		moveSpeed: f.moveSpeed,
		class:     class,
		active:    true,
	}
}
