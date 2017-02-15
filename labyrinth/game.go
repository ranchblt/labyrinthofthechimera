package labyrinth

import (
	"math/rand"
	"sync"
	"time"

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
	upKeys          []ebiten.Key
	downKeys        []ebiten.Key
	rand            *rand.Rand
}

// resources is where all the assets are stored
type resources struct {
	wizardImage      *ebiten.Image
	wizardShootImage *ebiten.Image
	heartImage       *ebiten.Image
	fireballSprite   []*ebiten.Image
	monsterImage     *ebiten.Image
	powerSpeedImage  *ebiten.Image
	monsterSprite    []*ebiten.Image
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

	fbc := &fireballCreator{
		images:        g.resources.fireballSprite,
		moveSpeed:     g.config.FireballSpeed,
		fastMoveSpeed: g.config.FastFireballSpeed,
		damage:        g.config.FireballDamage,
		powerDamage:   g.config.PowerFireballDamage,
	}

	wizard := newWizard(
		g.resources.wizardImage,
		g.resources.wizardShootImage,
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
		resources:        g.resources,
		wizard:           wizard,
		heartImage:       g.resources.heartImage,
		fastPowerupImage: g.resources.powerSpeedImage,
		rand:             g.rand,
		monsterImage:     g.resources.monsterImage,
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

		initImages(g.resources)
		logger.Debug("images loaded")
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

func initImages(res *resources) {
	i, err := openImage("wizard.png")
	handleErr(err)

	res.wizardImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	i, err = openImage("heart.png")
	handleErr(err)

	res.heartImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	g, err := openGif("fireball.gif")
	handleErr(err)

	res.fireballSprite = make([]*ebiten.Image, len(g.Image))

	for i, f := range g.Image {
		frame, err := ebiten.NewImageFromImage(f, ebiten.FilterNearest)
		handleErr(err)
		res.fireballSprite[i] = frame
	}

	i, err = openImage("monster.png")
	handleErr(err)

	res.monsterImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	i, err = openImage("powerup-speed.png")
	handleErr(err)

	res.powerSpeedImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	g, err = openGif("Snake.gif")
	handleErr(err)

	res.monsterSprite = make([]*ebiten.Image, len(g.Image))

	for i, f := range g.Image {
		frame, err := ebiten.NewImageFromImage(f, ebiten.FilterNearest)
		handleErr(err)
		res.monsterSprite[i] = frame
	}

	i, err = openImage("wizard-shooting.png")
	handleErr(err)

	res.wizardShootImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)
}
