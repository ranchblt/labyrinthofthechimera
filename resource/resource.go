package resource

import (
	"bytes"
	"image"
	"image/gif"

	"github.com/hajimehoshi/ebiten"
	"github.com/uber-go/zap"
)

// Resources is where all the assets are stored
type Resources struct {
	WizardImage      *ebiten.Image
	WizardShootImage *ebiten.Image
	HeartImage       *ebiten.Image
	FireballSprite   []*ebiten.Image
	MonsterImage     *ebiten.Image
	PowerSpeedImage  *ebiten.Image
	MonsterSprite    []*ebiten.Image
}

var loaded bool

func NewResources(logger zap.Logger) *Resources {
	res := &Resources{}
	if !loaded {
		initImages(res)
		logger.Debug("images loaded")
		loaded = true
	}
	return res
}

func initImages(res *Resources) {
	i, err := openImage("wizard.png")
	handleErr(err)

	res.WizardImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	i, err = openImage("heart.png")
	handleErr(err)

	res.HeartImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	g, err := openGif("fireball.gif")
	handleErr(err)

	res.FireballSprite = make([]*ebiten.Image, len(g.Image))

	for i, f := range g.Image {
		frame, err := ebiten.NewImageFromImage(f, ebiten.FilterNearest)
		handleErr(err)
		res.FireballSprite[i] = frame
	}

	i, err = openImage("monster.png")
	handleErr(err)

	res.MonsterImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	i, err = openImage("powerup-speed.png")
	handleErr(err)

	res.PowerSpeedImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)

	g, err = openGif("Snake.gif")
	handleErr(err)

	res.MonsterSprite = make([]*ebiten.Image, len(g.Image))

	for i, f := range g.Image {
		frame, err := ebiten.NewImageFromImage(f, ebiten.FilterNearest)
		handleErr(err)
		res.MonsterSprite[i] = frame
	}

	i, err = openImage("wizard-shooting.png")
	handleErr(err)

	res.WizardShootImage, err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
	handleErr(err)
}

// openImage gets the image out of go-bindata
func openImage(path string) (image.Image, error) {
	b, err := Asset(path)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	return image, nil
}

// openGif gets the gif out of go-bindata
func openGif(path string) (*gif.GIF, error) {
	b, err := Asset(path)
	if err != nil {
		return nil, err
	}

	image, err := gif.DecodeAll(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	return image, nil
}

// handleErr panics on any error, makes error handling cleaner
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
