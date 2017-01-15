package labyrinth

import "github.com/hajimehoshi/ebiten"

type wizard struct {
	image  *ebiten.Image
	Center *coord
}

// newWizard returns an initialized wizard
func newWizard(i *ebiten.Image) *wizard {
	c := &coord{
		x: 100,
		y: 100,
	}

	return &wizard{
		image:  i,
		Center: c,
	}
}

func (w *wizard) Update() error {
	return nil
}

func (w *wizard) Draw(r *ebiten.Image) error {

	r.DrawImage(w.image, &ebiten.DrawImageOptions{
		ImageParts: w,
	})

	return nil
}

func (w *wizard) Len() int {
	return 1
}

func (w *wizard) Dst(i int) (x0, y0, x1, y1 int) {
	width, height := w.image.Size()
	height = height / 2
	width = width / 2
	return w.Center.x - height,
		w.Center.y - width,
		w.Center.x + height,
		w.Center.y + width
}

func (w *wizard) Src(i int) (x0, y0, x1, y1 int) {
	width, height := w.image.Size()
	return 0, 0, width, height
}
