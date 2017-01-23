package labyrinth

import (
	"errors"

	"github.com/hajimehoshi/ebiten"
)

type fireballClass string

const (
	normalFireball = fireballClass("normal")
)

type fireball struct {
	image     *ebiten.Image
	moveSpeed int
	center    *coord
	class     fireballClass
}

func (f *fireball) Update() error {
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
	r.DrawImage(f.image, &ebiten.DrawImageOptions{
		ImageParts: f,
	})

	return nil
}

func (f *fireball) Len() int {
	return 1
}

func (f *fireball) updateNormalFireball() error {
	f.center.y += f.moveSpeed
	return nil
}

func (f *fireball) Dst(i int) (x0, y0, x1, y1 int) {
	width, height := f.image.Size()
	height = height / 2
	width = width / 2
	return f.center.x - height,
		f.center.y - width,
		f.center.x + height,
		f.center.y + width
}

func (f *fireball) Src(i int) (x0, y0, x1, y1 int) {
	width, height := f.image.Size()
	return 0, 0, width, height
}

type fireballCreator struct {
	image     *ebiten.Image
	moveSpeed int
}

func (f *fireballCreator) newFireball(c *coord, class fireballClass) *fireball {
	return &fireball{
		image:     f.image,
		center:    c,
		moveSpeed: f.moveSpeed,
		class:     class,
	}
}
