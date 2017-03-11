package settings

import (
	"github.com/ranchblt/labyrinthofthechimera/resource"

	"github.com/pelletier/go-toml"
)

// Config are settings that can be tweaked, easier to change in file
// no recompile required.
type Config struct {
	WizardMoveSpeed           int
	FireballSpeed             int
	FastFireballSpeed         int
	FireballDamage            int
	PowerFireballDamage       int
	Lives                     int
	UpKeys                    []string
	DownKeys                  []string
	MinPlayAreaHeight         int
	PowerupDespawnTime        int
	MonsterSpeedMultiplier    int
	MonsterPowerupHeal        int
	MonsterPushback           int
	MonsterPushbackMultiplier float64
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

	const fastFbSpeed = "fireball.fast_move_speed"
	checkRequired(t, fastFbSpeed)
	c.FastFireballSpeed = int(t.Get(fastFbSpeed).(int64))

	const fbDamage = "fireball.damage"
	checkRequired(t, fbDamage)
	c.FireballDamage = int(t.Get(fbDamage).(int64))

	const powerFbDamage = "fireball.power_damage"
	checkRequired(t, powerFbDamage)
	c.PowerFireballDamage = int(t.Get(powerFbDamage).(int64))

	const minHeight = "game.min_play_area_height"
	checkRequired(t, minHeight)
	c.MinPlayAreaHeight = int(t.Get(minHeight).(int64))

	const powerDespawn = "powerup.despawn_time"
	checkRequired(t, powerDespawn)
	c.PowerupDespawnTime = int(t.Get(powerDespawn).(int64))

	const monsterSpeedMultiplier = "monster.speed_multiplier"
	checkRequired(t, monsterSpeedMultiplier)
	c.MonsterSpeedMultiplier = int(t.Get(monsterSpeedMultiplier).(int64))

	const monsterPowerHeal = "monster.power_heal"
	checkRequired(t, monsterPowerHeal)
	c.MonsterPowerupHeal = int(t.Get(monsterPowerHeal).(int64))

	const monsterPushback = "monster.pushback"
	checkRequired(t, monsterPushback)
	c.MonsterPushback = int(t.Get(monsterPushback).(int64))

	const monsterPushbackMultiplier = "monster.pushback_multiplier"
	checkRequired(t, monsterPushbackMultiplier)
	c.MonsterPushbackMultiplier = t.Get(monsterPushbackMultiplier).(float64)

	return c
}

func checkRequired(t *toml.TomlTree, key string) {
	if !t.Has(key) {
		panic("config missing " + key)
	}
}
