package labyrinth

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ranchblt/labyrinthofthechimera/resource"
	"github.com/ranchblt/labyrinthofthechimera/settings"

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/statemanager"
	"github.com/uber-go/zap"
)

const ScreenWidth = 1280
const ScreenHeight = 720

// resources are all the assets loaded
var resources *resource.Resources

// config holds all config values that can be set from file
var config *settings.Config

// Game is a labyrinth game.
type Game struct {
	logger          zap.Logger
	stateManager    statemanager.StateManager
	keyboardWrapper *KeyboardWrapper
	upKeys          []ebiten.Key
	downKeys        []ebiten.Key
	rand            *rand.Rand
}

// NewGame returns a new labyrinth game.
func NewGame(debug *bool) *Game {
	lvl := zap.ErrorLevel
	if *debug {
		lvl = zap.DebugLevel
	}

	logger := zap.New(zap.NewTextEncoder(zap.TextNoTime()), lvl)
	keyboardWrapper := NewKeyboardWrapper()
	config = settings.New()

	g := &Game{
		logger:          logger,
		keyboardWrapper: keyboardWrapper,
	}

	g.load(logger)

	fbc := &fireballCreator{
		images:        resources.FireballSprite,
		moveSpeed:     config.FireballSpeed,
		fastMoveSpeed: config.FastFireballSpeed,
		damage:        config.FireballDamage,
		powerDamage:   config.PowerFireballDamage,
	}

	wizard := newWizard(
		config.WizardMoveSpeed,
		fbc,
		g.upKeys,
		g.downKeys,
		config.MinPlayAreaHeight,
	)

	stateManager := statemanager.New()
	stateManager.Add(&gameState{
		keyboardWrapper:  g.keyboardWrapper,
		wizard:           wizard,
		heartImage:       resources.HeartImage,
		fastPowerupImage: resources.PowerSpeedImage,
		rand:             g.rand,
		monsterImage:     resources.MonsterImage,
		lives:            config.Lives,
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

	wg.Add(2)

	go func(g *Game) {
		defer wg.Done()

		resources = resource.NewResources(logger)
		logger.Debug("resources loaded")
	}(g)

	go func(g *Game) {
		defer wg.Done()

		for _, k := range config.UpKeys {
			key, err := g.keyboardWrapper.Parse(k)
			if err != nil {
				panic("Invalid key in config " + k)
			}
			g.upKeys = append(g.upKeys, key)
		}

		for _, k := range config.DownKeys {
			key, err := g.keyboardWrapper.Parse(k)
			if err != nil {
				panic("Invalid key in config " + k)
			}
			g.downKeys = append(g.downKeys, key)
		}
	}(g)

	randSource := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(randSource)

	wg.Wait()
}
