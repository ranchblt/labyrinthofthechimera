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

// Game is a labyrinth game.
type Game struct {
	logger          zap.Logger
	stateManager    statemanager.StateManager
	keyboardWrapper *KeyboardWrapper
	config          *settings.Config
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

	g := &Game{
		logger:          logger,
		keyboardWrapper: keyboardWrapper,
		config:          settings.New(),
	}

	g.load(logger)

	fbc := &fireballCreator{
		images:        resources.FireballSprite,
		moveSpeed:     g.config.FireballSpeed,
		fastMoveSpeed: g.config.FastFireballSpeed,
		damage:        g.config.FireballDamage,
		powerDamage:   g.config.PowerFireballDamage,
	}

	wizard := newWizard(
		g.config.WizardMoveSpeed,
		fbc,
		g.upKeys,
		g.downKeys,
		g.config.MinPlayAreaHeight,
	)

	stateManager := statemanager.New()
	stateManager.Add(&gameState{
		keyboardWrapper:  g.keyboardWrapper,
		config:           g.config,
		wizard:           wizard,
		heartImage:       resources.HeartImage,
		fastPowerupImage: resources.PowerSpeedImage,
		rand:             g.rand,
		monsterImage:     resources.MonsterImage,
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

		for _, k := range g.config.UpKeys {
			key, err := g.keyboardWrapper.Parse(k)
			if err != nil {
				panic("Invalid key in config " + k)
			}
			g.upKeys = append(g.upKeys, key)
		}

		for _, k := range g.config.DownKeys {
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
