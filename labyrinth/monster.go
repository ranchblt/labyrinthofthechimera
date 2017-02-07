package labyrinth

import "github.com/hajimehoshi/ebiten"

type monster struct {
	health    int
	image     *ebiten.Image
	center    *coord
	active    bool
	moveClass movementClass
	speed     int
}

type movementClass string

const (
	straightLine = movementClass("straightLine")
)

func (m *monster) Update() error {
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
		r.DrawImage(m.image, &ebiten.DrawImageOptions{
			ImageParts: m,
		})
	}

	return nil
}

func (m *monster) straightLineMove() {
	m.center.x -= m.speed
}

func (m *monster) Len() int {
	return 1
}

func (m *monster) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultMobileDST(i, m.center, m.image)
}

func (m *monster) Src(i int) (x0, y0, x1, y1 int) {
	width, height := m.image.Size()
	return 0, 0, width, height
}

func (m *monster) offScreen() bool {
	w, _ := m.image.Size()
	return m.center.X()-w > ScreenWidth
}
