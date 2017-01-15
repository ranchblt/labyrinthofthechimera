package labyrinth

import (
	"sync"

	"github.com/ranchblt/labyrinthofthechimera/settings"

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
	resources       *resources
	config          *settings.Config
}

// resources is where all the assets are stored
type resources struct {
	wizardImage *ebiten.Image
	heartImage  *ebiten.Image
}

// NewGame returns a new labyrinth game.
func NewGame(debug *bool) *Game {
	lvl := zap.ErrorLevel
	if *debug {
		lvl = zap.DebugLevel
	}

	logger := zap.New(zap.NewTextEncoder(zap.TextNoTime()), lvl)
	keyboardWrapper := NewKeyboardWrapper()

	g := &Game{
		logger:          logger,
		keyboardWrapper: keyboardWrapper,
		resources:       &resources{},
		config:          settings.New(),
	}

	g.load(logger)

	wizard := newWizard(g.resources.wizardImage, g.config.WizardMoveSpeed)

	stateManager := statemanager.New()
	stateManager.Add(&gameState{
		keyboardWrapper: g.keyboardWrapper,
		wizard:          wizard,
		heartImage:      g.resources.heartImage,
		lives:           g.config.Lives,
		maxLives:        g.config.Lives,
	})
	stateManager.SetActive(gameStateID)

	g.stateManager = stateManager

	return g
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

func (g *Game) load(logger zap.Logger) {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(g *Game) {
		defer wg.Done()

		initImages(g.resources)
		logger.Debug("images loaded")
	}(g)

	wg.Wait()
}

func initImages(res *resources) {
	i, err := openImage("wizard.png")
	handleErr(err)

	res.wizardImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	i, err = openImage("heart.png")
	handleErr(err)

	res.heartImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

}
