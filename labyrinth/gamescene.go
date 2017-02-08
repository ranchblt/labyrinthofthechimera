package labyrinth

import (
	"errors"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/geometry-jumper/collision"
)

const gameStateID = "game"

type gameState struct {
	keyboardWrapper     *KeyboardWrapper
	wizard              *wizard
	heartImage          *ebiten.Image
	fastPowerupImage    *ebiten.Image
	powerups            []*powerup
	powerupTimer        *time.Timer
	powerupTimerStarted bool
	maxLives            int
	lives               int
	rand                *rand.Rand
	minPlayAreaHeight   int
	// TFE this is just for testing, should not stay this way
	monsterImage *ebiten.Image
	monster      *monster
	monsters     []*monster
}

func (s *gameState) OnEnter() error {
	// TFE this is just for testing, should not stay this way
	s.monster = &monster{
		health: 3,
		image:  s.monsterImage,
		center: &coord{
			x: ScreenWidth - 50,
			y: 360,
		},
		active:    true,
		moveClass: straightLine,
		speed:     1,
	}
	s.monsters = append(s.monsters, s.monster)
	return nil
}

func (s *gameState) OnExit() error {
	return nil
}

func (s *gameState) Draw(r *ebiten.Image) error {

	if err := s.drawLives(r); err != nil {
		return err
	}

	if err := s.wizard.Draw(r); err != nil {
		return err
	}

	for _, p := range s.powerups {
		if err := p.Draw(r); err != nil {
			return err
		}
	}

	for _, m := range s.monsters {
		if err := m.Draw(r); err != nil {
			return err
		}
	}

	return nil
}

func (s *gameState) drawLives(r *ebiten.Image) error {
	w, h := s.heartImage.Size()
	//fmt.Println(ScreenWidth - (s.maxLives * w))
	heartStartX := ScreenWidth
	for i := s.lives; i != 0; i-- {
		h := &Stationary{
			Image: s.heartImage,
			topLeft: &coord{
				x: heartStartX - i*w - (5 * i),
				y: h / 2,
			},
		}
		r.DrawImage(h.Image, &ebiten.DrawImageOptions{
			ImageParts: h,
		})
	}

	return nil
}

func (s *gameState) Update() error {
	for _, m := range s.monsters {
		if err := m.Update(); err != nil {
			return err
		}
	}

	if err := s.wizard.Update(s.keyboardWrapper); err != nil {
		return err
	}

	if s.keyboardWrapper.KeyPushed(ebiten.KeyEscape) {
		return errors.New("User wanted to quit") //Best way to do this?
	}

	nonExpiredPowups := []*powerup{}
	for _, p := range s.powerups {
		if err := p.Update(); err != nil {
			return err
		}
		if !p.expired {
			nonExpiredPowups = append(nonExpiredPowups, p)
		}
	}
	s.powerups = nonExpiredPowups

	if !s.powerupTimerStarted {
		// support random? number in future
		s.powerupTimer = time.NewTimer(time.Second)
		go s.spawnPowerup()
		s.powerupTimerStarted = true
	}

	s.collisions()

	return nil
}

func (s *gameState) collisions() {
	for _, fireball := range s.wizard.fireballs {
		fireballHitbox := collision.Hitbox{
			Image:  fireball.RGBAImage(),
			Center: fireball.center,
		}
		for _, monster := range s.monsters {
			monsterHitbox := collision.Hitbox{
				Image:  monster.RGBAImage(),
				Center: monster.center,
			}
			if collision.IsColliding(&fireballHitbox, &monsterHitbox) {
				fireball.hit()
				monster.hit(fireball)
			}
		}

		for _, powerup := range s.powerups {
			width, height := powerup.image.Size()
			width = width / 2
			height = height / 2
			powerupHitbox := collision.Hitbox{
				Image: powerup.RGBAImage(),
				// boy I hope this works.
				Center: &coord{
					x: powerup.topLeft.x + width,
					y: powerup.topLeft.y + height,
				},
			}
			if collision.IsColliding(&fireballHitbox, &powerupHitbox) {
				fireball.hit()
				s.wizard.activate(powerup)
				powerup.expired = true
			}
		}
	}
}

func (s *gameState) ID() string {
	return gameStateID
}

// spawnPowerup puts a powup randomly on the map so the wizard can shoot it.
// This should be inside the playable area. Will have to manage in the future
// so it doesn't spawn on top of a monster.
func (s *gameState) spawnPowerup() {
	<-s.powerupTimer.C
	width, _ := s.wizard.image.Size()
	padding := 50
	s.powerups = append(s.powerups, &powerup{
		image: s.fastPowerupImage,
		class: fastPowerup,
		topLeft: &coord{
			x: s.rand.Intn(ScreenWidth-width/2-s.wizard.Center.X()-padding*2) + width/2 + s.wizard.Center.X() + padding,
			y: s.rand.Intn(ScreenHeight-s.minPlayAreaHeight-padding*2) + s.minPlayAreaHeight + padding,
		},
	})
	s.powerupTimerStarted = false
}
