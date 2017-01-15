package labyrinth

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/statemanager"
	"github.com/uber-go/zap"
)

const ScreenWidth = 1280
const ScreenHeight = 720

// Game is a labyrinth game.
type Game struct {
	logger          zap.Logger
	stateManager    statemanager.StateManager
	keyboardWrapper *KeyboardWrapper
}

// NewGame returns a new labyrinth game.
func NewGame(debug *bool) *Game {
	lvl := zap.ErrorLevel
	if *debug {
		lvl = zap.DebugLevel
	}

	logger := zap.New(zap.NewTextEncoder(zap.TextNoTime()), lvl)
	keyboardWrapper := NewKeyboardWrapper()

	stateManager := statemanager.New()
	stateManager.Add(&gameState{
		keyboardWrapper: keyboardWrapper,
	})
	stateManager.SetActive(gameStateID)

	return &Game{
		logger:          logger,
		stateManager:    stateManager,
		keyboardWrapper: keyboardWrapper,
	}

}

// Update is the Game loop
func (g *Game) Update(r *ebiten.Image) error {
	g.keyboardWrapper.Update()

	if err := g.stateManager.Update(); err != nil {
		return err
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	if err := g.stateManager.Draw(r); err != nil {
		return err
	}

	return nil
}
