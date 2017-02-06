package labyrinth

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

type powerup struct {
	image   *ebiten.Image
	class   powerupClass
	expired bool
	x       int
	y       int
	timer   *time.Timer
}

type powerupClass string

const (
	fastPowerup = powerupClass("fast")
)

func (p *powerup) Update() error {
	// lasts x seconds then expires.
	// hit by wizard fireball expired.
	// hit by monster expired.
	return nil
}

func (p *powerup) Draw(r *ebiten.Image) error {
	if p.expired {
		return nil
	}

	return r.DrawImage(p.image, &ebiten.DrawImageOptions{
		ImageParts: p,
	})
}

func (p *powerup) Len() int {
	return 1
}

func (p *powerup) Dst(i int) (x0, y0, x1, y1 int) {
	width, height := p.image.Size()
	height = height / 2
	width = width / 2
	return p.x,
		p.y,
		p.x + width,
		p.y + height
}

func (p *powerup) Src(i int) (x0, y0, x1, y1 int) {
	width, height := p.image.Size()
	return 0, 0, width, height
}
