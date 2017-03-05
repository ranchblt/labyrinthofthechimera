package labyrinth

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten"
)

type fireballClass string

const (
	normalFireball        = fireballClass("normal")
	fastFireball          = fireballClass("fast")
	powerFireball         = fireballClass("power")
	fastPowerFireball     = fireballClass("fastPower")
	upFireball            = fireballClass("up")
	fastUpFireball        = fireballClass("fastUp")
	powerUpFireball       = fireballClass("powerUp")
	fastPowerUpFireball   = fireballClass("fastPowerUp")
	downFireball          = fireballClass("down")
	fastDownFireball      = fireballClass("fastDown")
	powerDownFireball     = fireballClass("powerDown")
	fastPowerDownFireball = fireballClass("fastPowerDown")
)

type direction int

const (
	up   = direction(-1)
	down = direction(1)
)

type fireball struct {
	sprite    Sprite
	rgba      *image.RGBA
	moveSpeed int
	topLeft   *coord
	// think this might be better as a slice actually?
	// might simplify the creation logic in both here & wizard a bit? not sure.
	// need to talk to mike about this. we'd have normal/fast for speed, normal/fast for power,
	// then up/down/straight for movement
	class  fireballClass
	active bool
	damage int
}

func (f *fireball) Update() error {
	if f.offScreen() {
		f.active = false
		return nil
	}

	switch f.class {
	// TODO this may be worth switching to individual cases, but I figure all but the
	// branching shot should be standard left to right movement
	case normalFireball, fastFireball, powerFireball, fastPowerFireball:
		if err := f.standardMovement(); err != nil {
			return err
		}
	case upFireball, fastUpFireball, powerUpFireball, fastPowerUpFireball:
		if err := f.diagonalMovement(up); err != nil {
			return err
		}
	case downFireball, fastDownFireball, powerDownFireball, fastPowerDownFireball:
		if err := f.diagonalMovement(down); err != nil {
			return err
		}
	default:
		return errors.New("Invalid class of fireball")
	}
	return nil
}

func (f *fireball) Draw(r *ebiten.Image) error {
	if f.active {
		return r.DrawImage(f.sprite.CurrentFrame(), &ebiten.DrawImageOptions{
			ImageParts: f,
		})
	}

	return nil
}

func (f *fireball) Len() int {
	return f.sprite.Len()
}

func (f *fireball) standardMovement() error {
	f.topLeft.x += f.moveSpeed
	return nil
}

func (f *fireball) diagonalMovement(direction direction) error {
	// the scaling here should probably be smarter. not sure how.
	// or we could have it go up a bit, then go straight??
	f.topLeft.x += f.moveSpeed
	f.topLeft.y -= (f.moveSpeed / 2) * int(direction)
	return nil
}

func (f *fireball) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultDST(i, f.topLeft, f.sprite.CurrentFrame())
}

func (f *fireball) Src(i int) (x0, y0, x1, y1 int) {
	width, height := f.sprite.CurrentFrame().Size()
	return 0, 0, width, height
}

// offscreen checks if the image is off an edge.
// TLX has crossed right edge
// BLY has crossed top edge
// TLY has crossed bottom edge
func (f *fireball) offScreen() bool {
	_, height := f.sprite.CurrentFrame().Size()

	return f.topLeft.X() > ScreenWidth || f.topLeft.Y()+height < 0 || f.topLeft.Y() > ScreenHeight
}

func (f *fireball) RGBAImage() *image.RGBA {
	if f.rgba == nil {
		f.rgba = toRGBA(f.sprite.CurrentFrame())
	}
	return f.rgba
}

func (f *fireball) hit() {
	// this should probably trigger some kind of animation. could also
	// potentially have powerups that spawn more shots or something. clusterbombs!
	f.active = false
}

type fireballCreator struct {
	images        []*ebiten.Image
	moveSpeed     int
	fastMoveSpeed int
	damage        int
	powerDamage   int
}

func (f *fireballCreator) newFireball(c coord, class fireballClass) *fireball {
	s := NewSprite(f.images, 70)
	go s.Animate()

	var fbSpeed int

	switch class {
	case fastFireball, fastPowerFireball, fastUpFireball, fastPowerUpFireball, fastDownFireball, fastPowerDownFireball:
		fbSpeed = f.fastMoveSpeed
	default:
		fbSpeed = f.moveSpeed
	}

	var fbDamage int
	switch class {
	case powerFireball, fastPowerFireball, powerUpFireball, fastPowerUpFireball, powerDownFireball, fastPowerDownFireball:
		fbDamage = f.powerDamage
	default:
		fbDamage = f.damage
	}

	return &fireball{
		sprite:    s,
		topLeft:   &c,
		moveSpeed: fbSpeed,
		class:     class,
		active:    true,
		damage:    fbDamage,
	}
}
