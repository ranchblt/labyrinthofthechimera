package labyrinth

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/uber-go/zap"
)

const ScreenWidth = 1280
const ScreenHeight = 720

// Game is a labyrinth game.
type Game struct {
	logger zap.Logger
}

// NewGame returns a new labyrinth game.
func NewGame(debug *bool) *Game {
	lvl := zap.ErrorLevel
	if *debug {
		lvl = zap.DebugLevel
	}

	logger := zap.New(zap.NewTextEncoder(zap.TextNoTime()), lvl)

	return &Game{
		logger: logger,
	}

}

// Update is the Game loop
func (g *Game) Update(r *ebiten.Image) error {
	return nil
}
