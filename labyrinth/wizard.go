package labyrinth

import "github.com/hajimehoshi/ebiten"

type wizard struct {
	image       *ebiten.Image
	Center      *coord
	moveSpeed   int
	fireballs   []*fireball
	fireCreator *fireballCreator
	upKeys      []ebiten.Key
	downKeys    []ebiten.Key
}

// newWizard returns an initialized wizard
func newWizard(i *ebiten.Image, speed int, fbc *fireballCreator, upkeys, downkeys []ebiten.Key) *wizard {
	c := &coord{
		x: 100,
		// 40 is half the image height for the wizard
		y: MinPlayAreaHeight + 40,
	}

	return &wizard{
		image:       i,
		Center:      c,
		moveSpeed:   speed,
		fireCreator: fbc,
		upKeys:      upkeys,
		downKeys:    downkeys,
	}
}

func (w *wizard) Update(keys *KeyboardWrapper) error {
	movedUp := false
	for _, k := range w.upKeys {
		if keys.IsKeyPressed(k) && !movedUp {
			w.moveUp()
			movedUp = true
		}
	}

	movedDown := false
	for _, k := range w.downKeys {
		if keys.IsKeyPressed(k) && !movedDown {
			w.moveDown()
			movedDown = true
		}
	}

	if keys.IsKeyPressed(ebiten.KeySpace) {
		if len(w.fireballs) == 0 {
			f := w.fireCreator.newFireball(*w.Center, normalFireball)
			w.fireballs = append(w.fireballs, f)
		}
	}

	activeFireballs := []*fireball{}
	for _, f := range w.fireballs {
		if err := f.Update(); err != nil {
			return err
		}
		if f.active {
			activeFireballs = append(activeFireballs, f)
		}
	}

	// Only keep active fireballs in the array
	w.fireballs = activeFireballs

	return nil
}

func (w *wizard) moveUp() {
	w.Center.y -= w.moveSpeed
	_, height := w.image.Size()
	if w.Center.y-(height/2) <= MinPlayAreaHeight {
		w.Center.y = MinPlayAreaHeight + (height / 2)
	}
}

func (w *wizard) moveDown() {
	w.Center.y += w.moveSpeed
	_, height := w.image.Size()
	if w.Center.y+(height/2) >= ScreenHeight {
		w.Center.y = ScreenHeight - (height / 2)
	}
}

func (w *wizard) Draw(r *ebiten.Image) error {
	r.DrawImage(w.image, &ebiten.DrawImageOptions{
		ImageParts: w,
	})

	for _, f := range w.fireballs {
		if err := f.Draw(r); err != nil {
			return err
		}
	}

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
