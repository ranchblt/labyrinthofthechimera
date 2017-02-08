package labyrinth

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type powerup struct {
	image   *ebiten.Image
	rgba    *image.RGBA
	class   powerupClass
	expired bool
	topLeft *coord
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
	return defaultStationaryDST(i, p.topLeft, p.image)
}

func (p *powerup) Src(i int) (x0, y0, x1, y1 int) {
	width, height := p.image.Size()
	return 0, 0, width, height
}

func (p *powerup) RGBAImage() *image.RGBA {
	if p.rgba == nil {
		p.rgba = toRGBA(p.image)
	}
	return p.rgba
}

func (p *powerup) Center() *coord {
	width, height := p.image.Size()
	width = width / 2
	height = height / 2
	return &coord{
		x: p.topLeft.x + width,
		y: p.topLeft.y + height,
	}
}
