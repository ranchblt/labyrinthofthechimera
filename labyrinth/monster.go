package labyrinth

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type monster struct {
	maxHealth       int
	health          int
	sprite          Sprite
	rgba            *image.RGBA
	topLeft         *coord
	topLeftStart    *coord
	active          bool
	moveClass       movementClass
	speed           int
	speedMultiplier int
	powerupHeal     int
	powerups        []*powerup
	attacked        bool
}

type movementClass string

const (
	straightLine = movementClass("straightLine")
)

func (m *monster) Update() error {
	if m.offScreen() {
		m.active = false
		return nil
	}

	monsterSpeed := m.getMoveSpeed()
	switch m.moveClass {
	case straightLine:
		m.straightLineMove(monsterSpeed)
	default:
		panic("unknown move class")
	}

	m.removePowerups()
	return nil
}

func (m *monster) removePowerups() error {
	activePowerups := []*powerup{}

	for _, p := range m.powerups {
		if p.active {
			activePowerups = append(activePowerups, p)

		}
	}

	m.powerups = activePowerups

	return nil
}

func (m *monster) Draw(r *ebiten.Image) error {
	if m.active {
		err := r.DrawImage(m.sprite.CurrentFrame(), &ebiten.DrawImageOptions{
			ImageParts: m,
		})
		if err != nil {
			return err
		}
		if err := m.drawHearts(r); err != nil {
			return err
		}
	}

	return nil
}

// drawHearts puts the hearts above the monter's head based on their hp
func (m *monster) drawHearts(r *ebiten.Image) error {
	scale := .3
	w, h := resources.HeartImage.Size()
	w = int(float64(w) * scale)
	h = int(float64(h) * scale)
	heartStartX := m.topLeft.X()
	for i := 0; i < m.health; i++ {
		blankSpace := 5 * i
		offset := i * w
		h := &Stationary{
			Image: resources.HeartImage,
			topLeft: &coord{
				x: heartStartX + offset + blankSpace,
				y: m.topLeft.Y() - h,
			},
			scale: scale,
		}
		err := r.DrawImage(h.Image, &ebiten.DrawImageOptions{
			ImageParts: h,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *monster) straightLineMove(monsterSpeed int) {
	m.topLeft.x -= monsterSpeed
}

func (m *monster) getMoveSpeed() int {
	monsterSpeed := m.speed
	if containsPowerupClass(m.powerups, fastFireballPowerup) {
		monsterSpeed *= m.speedMultiplier
	}
	return monsterSpeed
}

func (m *monster) Len() int {
	return m.sprite.Len()
}

func (m *monster) Dst(i int) (x0, y0, x1, y1 int) {
	return defaultDST(i, m.topLeft, m.sprite.CurrentFrame())
}

func (m *monster) Src(i int) (x0, y0, x1, y1 int) {
	width, height := m.sprite.CurrentFrame().Size()
	return 0, 0, width, height
}

func (m *monster) offScreen() bool {
	w, _ := m.sprite.CurrentFrame().Size()
	return m.topLeft.X()+w > ScreenWidth
}

func (m *monster) RGBAImage() *image.RGBA {
	if m.rgba == nil {
		m.rgba = toRGBA(m.sprite.CurrentFrame())
	}
	return m.rgba
}

func (m *monster) AttackRange() bool {
	if !m.attacked {
		a := m.topLeft.X() <= 0
		if !m.attacked && a {
			m.attacked = a
			return true
		}
	}

	return false
}

func (m *monster) hit(fireball *fireball) {
	// TODO take the damage of the fireball into account somehow...
	m.health -= fireball.damage
	if m.health <= 0 {
		// TODO should have some kind of "I DIED" animation somehow...
		m.active = false
	}
}

func (m *monster) powerup(p *powerup) {
	m.powerups = append(m.powerups, p)
	p.Activate()
	if p.class == powerFireballPowerup {
		m.health += m.powerupHeal
		if m.maxHealth < m.health {
			m.health = m.maxHealth
		}
	}
}

func copyMonster(m *monster) *monster {
	monster := &monster{
		maxHealth: m.maxHealth,
		health:    m.maxHealth,
		sprite:    m.sprite,
		topLeft: &coord{
			x: m.topLeftStart.x,
			y: m.topLeftStart.y,
		},
		topLeftStart: &coord{
			x: m.topLeftStart.x,
			y: m.topLeftStart.y,
		},
		active:          true,
		moveClass:       m.moveClass,
		speed:           m.speed,
		speedMultiplier: m.speedMultiplier,
		powerupHeal:     m.powerupHeal,
	}
	return monster
}
