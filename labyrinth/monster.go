package labyrinth

import "github.com/hajimehoshi/ebiten"

type monster struct {
	health int
	image  *ebiten.Image
	center *coord
	active bool
}

func (f *monster) Update() error {
	if f.offScreen() {
		f.active = false
		return nil
	}

	return nil
}

func (f *monster) Draw(r *ebiten.Image) error {
	if f.active {
		r.DrawImage(f.image, &ebiten.DrawImageOptions{
			ImageParts: f,
		})
	}

	return nil
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
