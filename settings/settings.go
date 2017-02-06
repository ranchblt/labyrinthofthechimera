package settings

import (
	"github.com/ranchblt/labyrinthofthechimera/resource"

	"github.com/pelletier/go-toml"
)

// Config are settings that can be tweaked, easier to change in file
// no recompile required.
type Config struct {
	WizardMoveSpeed   int
	FireballSpeed     int
	Lives             int
	UpKeys            []string
	DownKeys          []string
	MinPlayAreaHeight int
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
	checkRequired(t, moveSpeed)
	c.WizardMoveSpeed = int(t.Get(moveSpeed).(int64))

	const lives = "game.lives"
	checkRequired(t, lives)
	c.Lives = int(t.Get(lives).(int64))

	const upKeys = "game.up_keys"
	checkRequired(t, upKeys)
	c.UpKeys = []string{}
	temp := t.Get(upKeys).([]interface{})
	for _, i := range temp {
		c.UpKeys = append(c.UpKeys, i.(string))
	}

	const downKeys = "game.down_keys"
	checkRequired(t, downKeys)
	c.DownKeys = []string{}
	temp = t.Get(downKeys).([]interface{})
	for _, i := range temp {
		c.DownKeys = append(c.DownKeys, i.(string))
	}

	const fbSpeed = "fireball.move_speed"
	checkRequired(t, fbSpeed)
	c.FireballSpeed = int(t.Get(fbSpeed).(int64))

	const minHeight = "game.min_play_area_height"
	checkRequired(t, minHeight)
	c.MinPlayAreaHeight = int(t.Get(minHeight).(int64))

	return c
}

func checkRequired(t *toml.TomlTree, key string) {
	if !t.Has(key) {
		panic("config missing " + key)
	}
}
