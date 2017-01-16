package labyrinth

import (
	"errors"

	"github.com/hajimehoshi/ebiten"
)

const gameStateID = "game"

type gameState struct {
	keyboardWrapper *KeyboardWrapper
	wizard          *wizard
	heartImage      *ebiten.Image
	maxLives        int
	lives           int
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

	return nil
}

func (s *gameState) drawLives(r *ebiten.Image) error {
	w, h := s.heartImage.Size()
	//fmt.Println(ScreenWidth - (s.maxLives * w))
	heartStartX := ScreenWidth - (s.maxLives * w / 2)
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

	return nil
}

func (s *gameState) ID() string {
	return gameStateID
}
