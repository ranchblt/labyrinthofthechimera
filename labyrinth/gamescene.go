package labyrinth

import (
	"errors"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten"
)

const gameStateID = "game"

type gameState struct {
	keyboardWrapper   *KeyboardWrapper
	wizard            *wizard
	heartImage        *ebiten.Image
	fastPowerupImage  *ebiten.Image
	powerups          []*powerup
	powerupTimer      *time.Timer
	powerupSpawned    bool
	maxLives          int
	lives             int
	rand              *rand.Rand
	minPlayAreaHeight int
}

func (s *gameState) OnEnter() error {
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

	return nil
}

func (s *gameState) drawLives(r *ebiten.Image) error {
	w, h := s.heartImage.Size()
	//fmt.Println(ScreenWidth - (s.maxLives * w))
	heartStartX := ScreenWidth
	for i := s.lives; i != 0; i-- {
		h := &Stationary{
			Image: s.heartImage,
			X:     heartStartX - i*w - (5 * i),
			Y:     h / 2,
		}
		r.DrawImage(h.Image, &ebiten.DrawImageOptions{
			ImageParts: h,
		})
	}

	return nil
}

func (s *gameState) Update() error {
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

	if !s.powerupSpawned {
		// support random? number in future
		s.powerupTimer = time.NewTimer(time.Second * 2)
		go s.spawnPowerup()
	}

	return nil
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
	s.powerups = append(s.powerups, &powerup{
		image: s.fastPowerupImage,
		class: fastPowerup,
		x:     s.rand.Intn(ScreenHeight-s.minPlayAreaHeight) + s.minPlayAreaHeight,
		y:     s.rand.Intn(ScreenWidth-width) + width,
	})
	s.powerupSpawned = true
}
