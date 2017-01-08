package labyrinth

import "github.com/hajimehoshi/ebiten"

const gameStateId = "game"

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
	return nil
}

func (s *gameState) ID() string {
	return gameStateId
}
