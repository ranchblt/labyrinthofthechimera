package labyrinth

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Sprite is an animated slice of ebiten.Image
type Sprite interface {
	CurrentFrame() *ebiten.Image
	Animate()
	Len() int
}

type ebitenSprite struct {
	images      []*ebiten.Image
	frame       int
	frameTicker *time.Ticker
}

// NewSprite returns a new animatable sprite
func NewSprite(images []*ebiten.Image, miliSecondsBtwFrames int) Sprite {
	return &ebitenSprite{
		images:      images,
		frame:       0,
		frameTicker: time.NewTicker(time.Millisecond * time.Duration(miliSecondsBtwFrames)),
	}
}

// CurrentFrame returns the image to draw now
func (s *ebitenSprite) CurrentFrame() *ebiten.Image {
	return s.images[s.frame]
}

// Animate is required to run in the background and will run the Ticker
// that updates which frame this sprite is on.
func (s *ebitenSprite) Animate() {
	for range s.frameTicker.C {
		s.frame++
		if s.frame >= s.Len() {
			s.frame = 0
		}
	}
}

func (s *ebitenSprite) Len() int {
	return len(s.images)
}
