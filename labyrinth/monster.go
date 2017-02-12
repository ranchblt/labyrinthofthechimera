package labyrinth

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type monster struct {
	health    int
	sprite    []*ebiten.Image
	frame     int
	rgba      *image.RGBA
	topLeft   *coord
	active    bool
	moveClass movementClass
	speed     int
}

type movementClass string

const (
	straightLine = movementClass("straightLine")
)

func (m *monster) Update() error {
	m.frame++
	if m.frame >= m.Len() {
		m.frame = 0
	}
	if m.offScreen() {
		m.active = false
		return nil
	}

	switch m.moveClass {
	case straightLine:
		m.straightLineMove()
	default:
		panic("unknown move class")
	}

	return nil
}

func (m *monster) Draw(r *ebiten.Image) error {
	if m.active {
		r.DrawImage(m.sprite[m.frame], &ebiten.DrawImageOptions{
			ImageParts: m,
		})
	}

	return nil
}

func (m *monster) straightLineMove() {
	m.topLeft.x -= m.speed
}

func (m *monster) Len() int {
	return len(m.sprite)
}

func (m *monster) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultStationaryDST(i, m.topLeft, m.sprite[i])
}

func (m *monster) Src(i int) (x0, y0, x1, y1 int) {
	width, height := m.sprite[i].Size()
	return 0, 0, width, height
}

func (m *monster) offScreen() bool {
	w, _ := m.sprite[0].Size()
	return m.topLeft.X()+w > ScreenWidth
}

func (m *monster) RGBAImage() *image.RGBA {
	if m.rgba == nil {
		m.rgba = toRGBA(m.sprite[0])
	}
	return m.rgba
}

func (m *monster) hit(fireball *fireball) {
	// TODO take the damage of the fireball into account somehow...
	m.health -= 1
	if m.health <= 0 {
		// TODO should have some kind of "I DIED" animation somehow...
		m.active = false
	}
}

func (m *monster) powerup(powerup *powerup) {
	switch {
	case powerup.class == fastPowerup:
		m.speed += 5
	}
}
