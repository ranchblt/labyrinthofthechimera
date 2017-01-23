package settings

import (
	"github.com/ranchblt/labyrinthofthechimera/resource"

	"github.com/pelletier/go-toml"
)

// Config are settings that can be tweaked, easier to change in file
// no recompile required.
type Config struct {
	WizardMoveSpeed int
	FireballSpeed   int
	Lives           int
}

// New gets a new config loaded from file
func New() *Config {
	b, err := resource.Asset("config.toml")

	c := &Config{}

	t, err := toml.Load(string(b))
	if err != nil {
		panic("Failed to parse config " + err.Error())
	}

	const moveSpeed = "wizard.move_speed"
	if !t.Has(moveSpeed) {
		panic("config missing " + moveSpeed)
	}
	c.WizardMoveSpeed = int(t.Get(moveSpeed).(int64))

	const lives = "game.lives"
	if !t.Has(lives) {
		panic("config missing " + lives)
	}
	c.Lives = int(t.Get(lives).(int64))

	const fbSpeed = "fireball.move_speed"
	if !t.Has(fbSpeed) {
		panic("config missing " + fbSpeed)
	}
	c.FireballSpeed = int(t.Get(fbSpeed).(int64))

	return c
}
