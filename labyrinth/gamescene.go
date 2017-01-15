package labyrinth

import (
	"errors"

	"github.com/hajimehoshi/ebiten"
)

const gameStateID = "game"

type gameState struct {
	keyboardWrapper *KeyboardWrapper
}

func (s *gameState) OnEnter() error {
	return nil
}

func (s *gameState) OnExit() error {
	return nil
}

func (s *gameState) Draw(screen *ebiten.Image) error {
	return nil
}

func (s *gameState) Update() error {
	if s.keyboardWrapper.KeyPushed(ebiten.KeyEscape) {
		return errors.New("User wanted to quit") //Best way to do this?
	}

	return nil
}

func (s *gameState) ID() string {
	return gameStateID
}
