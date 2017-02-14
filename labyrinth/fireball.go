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
	sprite    Sprite
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
		r.DrawImage(f.sprite.CurrentFrame(), &ebiten.DrawImageOptions{
			ImageParts: f,
		})
	}

	return nil
}

func (f *fireball) Len() int {
	return f.sprite.Len()
}

func (f *fireball) updateNormalFireball() error {
	f.topLeft.x += f.moveSpeed
	return nil
}

func (f *fireball) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultDST(i, f.topLeft, f.sprite.CurrentFrame())
}

func (f *fireball) Src(i int) (x0, y0, x1, y1 int) {
	width, height := f.sprite.CurrentFrame().Size()
	return 0, 0, width, height
}

// offscreen checks if the left most part of the image is past the ScreenWidth
func (f *fireball) offScreen() bool {
	return f.topLeft.X() > ScreenWidth
}

func (f *fireball) RGBAImage() *image.RGBA {
	if f.rgba == nil {
		f.rgba = toRGBA(f.sprite.CurrentFrame())
	}
	return f.rgba
}

func (f *fireball) hit() {
	// this should probably trigger some kind of animation. could also
	// potentially have powerups that spawn more shots or something. clusterbombs!
	f.active = false
}

type fireballCreator struct {
	images    []*ebiten.Image
	moveSpeed int
}

func (f *fireballCreator) newFireball(c coord, class fireballClass, calculatedMoveSpeed int) *fireball {
	s := NewSprite(f.images, 70)
	go s.Animate()
	return &fireball{
		sprite:    s,
		topLeft:   &c,
		moveSpeed: calculatedMoveSpeed,
		class:     class,
		active:    true,
	}
}
