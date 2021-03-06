package labyrinth

import (
	"errors"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/labyrinthofthechimera/collision"
)

const gameStateID = "game"

type gameState struct {
	keyboardWrapper     *KeyboardWrapper
	wizard              *wizard
	heartImage          *ebiten.Image
	fastPowerupImage    *ebiten.Image
	powerups            []*powerup
	powerupTimer        *time.Timer
	powerupTimerStarted bool
	rand                *rand.Rand
	// TFE this is just for testing, should not stay this way
	monsterImage *ebiten.Image
	monster      *monster
	monsters     []*monster
	lives        int
	pausing      bool
}

func (s *gameState) OnEnter() error {
	// TFE this is just for testing, should not stay this way
	s.monster = &monster{
		maxHealth: 3,
		health:    3,
		sprite:    NewSprite(resources.MonsterSprite, 300),
		topLeft: &coord{
			x: ScreenWidth - 50,
			y: 360,
		},
		topLeftStart: &coord{
			x: ScreenWidth - 50,
			y: 360,
		},
		active:          true,
		moveClass:       straightLine,
		speed:           1,
		speedMultiplier: config.MonsterSpeedMultiplier,
		powerupHeal:     config.MonsterPowerupHeal,
	}

	go s.monster.sprite.Animate()
	s.monsters = append(s.monsters, s.monster)
	return nil
}

func (s *gameState) OnExit() error {
	return nil
}

func (s *gameState) Draw(r *ebiten.Image) error {

	if err := s.drawLives(r); err != nil {
		return err
	}

	if err := s.wizard.Draw(r); err != nil {
		return err
	}

	for _, p := range s.powerups {
		if err := p.Draw(r); err != nil {
			return err
		}
	}

	for _, m := range s.monsters {
		if err := m.Draw(r); err != nil {
			return err
		}
	}

	return nil
}

func (s *gameState) drawLives(r *ebiten.Image) error {
	w, h := s.heartImage.Size()
	//fmt.Println(ScreenWidth - (s.maxLives * w))
	heartStartX := ScreenWidth
	for i := config.Lives; i != config.Lives-s.lives; i-- {
		h := &Stationary{
			Image: s.heartImage,
			topLeft: &coord{
				x: heartStartX - i*w - (5 * i),
				y: h / 2,
			},
			scale: 1,
		}
		r.DrawImage(h.Image, &ebiten.DrawImageOptions{
			ImageParts: h,
		})
	}

	return nil
}

func (s *gameState) Update() error {
	if s.keyboardWrapper.KeyPushed(ebiten.KeyEscape) {
		return errors.New("User wanted to quit") //Best way to do this?
	}

	if !s.pausing {
		s.updateMonsters()

		if err := s.wizard.Update(s.keyboardWrapper); err != nil {
			return err
		}

		s.updatePowerups()

		s.collisions()
	}

	return nil
}

func (s *gameState) collisions() {
	for _, fireball := range s.wizard.fireballs {
		fireballHitbox := collision.Hitbox{
			Image:   fireball.RGBAImage(),
			TopLeft: fireball.topLeft,
		}
		for _, monster := range s.monsters {
			monsterHitbox := collision.Hitbox{
				Image:   monster.RGBAImage(),
				TopLeft: monster.topLeft,
			}
			if collision.IsColliding(&fireballHitbox, &monsterHitbox) {
				fireball.hit()
				monster.hit(fireball)
				s.pausing = true
				go s.stopPausing()
			}
		}

		for _, powerup := range s.powerups {
			powerupHitbox := collision.Hitbox{
				Image: powerup.RGBAImage(),
				// boy I hope this works.
				TopLeft: powerup.topLeft,
			}
			if collision.IsColliding(&fireballHitbox, &powerupHitbox) {
				fireball.hit()
				s.wizard.activate(powerup)
				powerup.expired = true
			}
		}
	}

	// This needs to be it's own because the loop before will only
	// run if there is a fireball ball out and this needs to be
	// checked all the time.
	for _, powerup := range s.powerups {
		powerupHitbox := collision.Hitbox{
			Image:   powerup.RGBAImage(),
			TopLeft: powerup.topLeft,
		}
		for _, monster := range s.monsters {
			monsterHitbox := collision.Hitbox{
				Image:   monster.RGBAImage(),
				TopLeft: monster.topLeft,
			}

			if collision.IsColliding(&powerupHitbox, &monsterHitbox) {
				// wish this code could go in monster, but I don't know how?
				if powerup.class == branchFireballPowerup {
					copiedMonster := copyMonster(monster)
					s.monsters = append(s.monsters, copiedMonster)
				}
				monster.powerup(powerup)
				powerup.expired = true
			}
		}
	}
}

func (s *gameState) updateMonsters() error {
	monsters := []*monster{}
	for _, m := range s.monsters {
		if err := m.Update(); err != nil {
			return err
		}

		if m.AttackRange() {
			s.lives = s.lives - 1
		}

		if m.active {
			monsters = append(monsters, m)
		}
	}

	// Only keep active monsters in the array
	s.monsters = monsters
	return nil
}

func (s *gameState) updatePowerups() error {
	nonExpiredPowups := []*powerup{}
	for _, p := range s.powerups {
		if err := p.Update(); err != nil {
			return err
		}
		if !p.expired {
			nonExpiredPowups = append(nonExpiredPowups, p)
		}
	}
	s.powerups = nonExpiredPowups

	if !s.powerupTimerStarted {
		// support random? number in future
		s.powerupTimer = time.NewTimer(time.Second)
		go s.spawnPowerup()
		s.powerupTimerStarted = true
	}

	return nil
}

func (s *gameState) ID() string {
	return gameStateID
}

// spawnPowerup puts a powup randomly on the map so the wizard can shoot it.
// This should be inside the playable area. Will have to manage in the future
// so it doesn't spawn on top of a monster.
func (s *gameState) spawnPowerup() {
	<-s.powerupTimer.C
	width, _ := s.wizard.image.Size()
	padding := 50
	p := &powerup{
		image: s.fastPowerupImage,
		class: branchFireballPowerup,
		topLeft: &coord{
			x: s.rand.Intn(ScreenWidth-width-s.wizard.TopLeft.X()-padding*2) + width + s.wizard.TopLeft.X() + padding,
			//y: s.rand.Intn(ScreenHeight-s.config.MinPlayAreaHeight-padding*2) + s.config.MinPlayAreaHeight + padding,
			y: 400,
		},
		timer: time.NewTimer(time.Second * time.Duration(config.PowerupDespawnTime)),
		// TODO have this be configurable?
		durationMillis: 3000,
		boost:          2,
	}
	s.powerups = append(s.powerups, p)
	go p.Despawn()

	s.powerupTimerStarted = false
}

func (s *gameState) stopPausing() {
	t := time.NewTimer(time.Millisecond * time.Duration(config.HitPause))
	<-t.C
	s.pausing = false
}
