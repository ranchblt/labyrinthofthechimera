package labyrinth

import "github.com/hajimehoshi/ebiten"

type wizard struct {
	image             *ebiten.Image
	TopLeft           *coord
	moveSpeed         int
	fireballs         []*fireball
	fireCreator       *fireballCreator
	upKeys            []ebiten.Key
	downKeys          []ebiten.Key
	minPlayAreaHeight int
}

// newWizard returns an initialized wizard
func newWizard(i *ebiten.Image, speed int, fbc *fireballCreator, upkeys, downkeys []ebiten.Key, minPlayAreaHeight int) *wizard {
	c := &coord{
		x: 50,
		y: minPlayAreaHeight,
	}

	return &wizard{
		image:             i,
		TopLeft:           c,
		moveSpeed:         speed,
		fireCreator:       fbc,
		upKeys:            upkeys,
		downKeys:          downkeys,
		minPlayAreaHeight: minPlayAreaHeight,
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
			// We cannot send a pointer to newFireball because then when we move
			// the fireball it will move the wizard! Comical, yes but not desired.
			f := w.fireCreator.newFireball(*center(w.TopLeft, w.image), normalFireball)
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
	w.TopLeft.y -= w.moveSpeed
	if w.TopLeft.y <= w.minPlayAreaHeight {
		w.TopLeft.y = w.minPlayAreaHeight
	}
}

func (w *wizard) moveDown() {
	w.TopLeft.y += w.moveSpeed
	_, h := w.image.Size()
	if w.TopLeft.y+h >= ScreenHeight {
		w.TopLeft.y = ScreenHeight - h
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
	return defaultStationaryDST(i, w.TopLeft, w.image)
}

func (w *wizard) Src(i int) (x0, y0, x1, y1 int) {
	width, height := w.image.Size()
	return 0, 0, width, height
}

func (w *wizard) activate(p *powerup) {
	// TODO implement powerup logic here, somehow.
	switch {
	case p.class == fastPowerup:
		w.moveSpeed += 5
	}
}
