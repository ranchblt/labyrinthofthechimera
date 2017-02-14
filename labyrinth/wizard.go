package labyrinth

import "github.com/hajimehoshi/ebiten"

type wizard struct {
	image             *ebiten.Image
	TopLeft           *coord
	moveSpeed         int
	fireballs         []*fireball
	powerups          []*powerup
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
		powerups:          []*powerup{},
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
		w.generateFireball()
	}

	w.updateFireballs()

	w.removePowerups()

	return nil
}

func (w *wizard) generateFireball() {
	if len(w.fireballs) == 0 {
		// We cannot send a pointer to newFireball because then when we move
		// the fireball it will move the wizard! Comical, yes but not desired.

		if containsPowerupClass(w.powerups, branchFireballPowerup) {
			// branching fireballs are a bit different. we need three new fireballs,
			// one for up, center, and down.

			// center is just a standard one
			centerClass := w.getFireballClass(normalFireball, fastFireball, powerFireball, fastPowerFireball)
			centerF := w.fireCreator.newFireball(*center(w.TopLeft, w.image), centerClass)
			w.fireballs = append(w.fireballs, centerF)

			// then up & down need their own special classes, so we can
			// use the right movement behavior
			upClass := w.getFireballClass(upFireball, fastUpFireball, powerUpFireball, fastPowerUpFireball)
			upF := w.fireCreator.newFireball(*center(w.TopLeft, w.image), upClass)
			w.fireballs = append(w.fireballs, upF)

			downClass := w.getFireballClass(downFireball, fastDownFireball, powerDownFireball, fastPowerDownFireball)
			downF := w.fireCreator.newFireball(*center(w.TopLeft, w.image), downClass)
			w.fireballs = append(w.fireballs, downF)

		} else {
			class := w.getFireballClass(normalFireball, fastFireball, powerFireball, fastPowerFireball)
			f := w.fireCreator.newFireball(*center(w.TopLeft, w.image), class)
			w.fireballs = append(w.fireballs, f)
		}
	}
}

// each parameter determines what should be returned when powerups of that type are found.
// this isn't necessarily expandible easily, but I don't see us having too many fireball types.
func (w *wizard) getFireballClass(normalClass fireballClass, fastClass fireballClass, powerClass fireballClass, fastPowerClass fireballClass) fireballClass {
	fast := containsPowerupClass(w.powerups, fastFireballPowerup)
	power := containsPowerupClass(w.powerups, powerFireballPowerup)

	if fast && power {
		return fastPowerClass
	} else if power {
		return powerClass
	} else if fast {
		return fastClass
	}
	return normalClass
}

func (w *wizard) updateFireballs() error {
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

func (w *wizard) removePowerups() error {
	activePowerups := []*powerup{}

	for _, p := range w.powerups {
		if p.active {
			activePowerups = append(activePowerups, p)
		}
	}

	w.powerups = activePowerups

	return nil
}

func (w *wizard) moveUp() {
	w.TopLeft.y -= getPoweredUpValue(w.moveSpeed, w.powerups, fastPlayerPowerup)
	if w.TopLeft.y <= w.minPlayAreaHeight {
		w.TopLeft.y = w.minPlayAreaHeight
	}
}

func (w *wizard) moveDown() {
	w.TopLeft.y += getPoweredUpValue(w.moveSpeed, w.powerups, fastPlayerPowerup)
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
	return defaultDST(i, w.TopLeft, w.image)
}

func (w *wizard) Src(i int) (x0, y0, x1, y1 int) {
	width, height := w.image.Size()
	return 0, 0, width, height
}

func (w *wizard) activate(p *powerup) {
	w.powerups = append(w.powerups, p)
	p.Activate()
}
