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
	// whether this powerup is active on the thing using it
	active         bool
	durationMillis int
	// this is kinda vestigial at this point, it only applies to fastPlayerPowerup
	// which we're not enamored with. can probably be removed.
	boost int
}

type powerupClass string

const (
	fastPlayerPowerup     = powerupClass("fastPlayer")
	fastFireballPowerup   = powerupClass("fastFireball")
	powerFireballPowerup  = powerupClass("powerFireball")
	branchFireballPowerup = powerupClass("branchFireball")
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
	return defaultDST(i, p.topLeft, p.image)
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

func (p *powerup) Despawn() {
	<-p.timer.C
	p.expired = true
}

func (p *powerup) Activate() {
	p.active = true
	go p.Deactivate()
}

func (p *powerup) Deactivate() {
	timer := time.NewTimer(time.Millisecond * time.Duration(p.durationMillis))
	<-timer.C
	p.active = false
}
