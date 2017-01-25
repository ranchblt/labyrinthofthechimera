package labyrinth

import (
	"errors"
	"strings"

	"github.com/hajimehoshi/ebiten"
)

type KeyboardWrapper struct {
	keyState map[ebiten.Key]int
	keys     []ebiten.Key
}

func NewKeyboardWrapper() *KeyboardWrapper {
	var kw = &KeyboardWrapper{
		keyState: map[ebiten.Key]int{},
		keys: []ebiten.Key{
			ebiten.Key0,
			ebiten.Key1,
			ebiten.Key2,
			ebiten.Key3,
			ebiten.Key4,
			ebiten.Key5,
			ebiten.Key6,
			ebiten.Key7,
			ebiten.Key8,
			ebiten.Key9,
			ebiten.KeyA,
			ebiten.KeyB,
			ebiten.KeyC,
			ebiten.KeyD,
			ebiten.KeyE,
			ebiten.KeyF,
			ebiten.KeyG,
			ebiten.KeyH,
			ebiten.KeyI,
			ebiten.KeyJ,
			ebiten.KeyK,
			ebiten.KeyL,
			ebiten.KeyM,
			ebiten.KeyN,
			ebiten.KeyO,
			ebiten.KeyP,
			ebiten.KeyQ,
			ebiten.KeyR,
			ebiten.KeyS,
			ebiten.KeyT,
			ebiten.KeyU,
			ebiten.KeyV,
			ebiten.KeyW,
			ebiten.KeyX,
			ebiten.KeyY,
			ebiten.KeyZ,
			ebiten.KeyAlt,
			ebiten.KeyBackspace,
			ebiten.KeyCapsLock,
			ebiten.KeyComma,
			ebiten.KeyControl,
			ebiten.KeyDelete,
			ebiten.KeyDown,
			ebiten.KeyEnd,
			ebiten.KeyEnter,
			ebiten.KeyEscape,
			ebiten.KeyF1,
			ebiten.KeyF2,
			ebiten.KeyF3,
			ebiten.KeyF4,
			ebiten.KeyF5,
			ebiten.KeyF6,
			ebiten.KeyF7,
			ebiten.KeyF8,
			ebiten.KeyF9,
			ebiten.KeyF10,
			ebiten.KeyF11,
			ebiten.KeyF12,
			ebiten.KeyHome,
			ebiten.KeyInsert,
			ebiten.KeyLeft,
			ebiten.KeyPageDown,
			ebiten.KeyPageUp,
			ebiten.KeyPeriod,
			ebiten.KeyRight,
			ebiten.KeyShift,
			ebiten.KeySpace,
			ebiten.KeyTab,
			ebiten.KeyUp,
		},
	}
	return kw
}

func (kw *KeyboardWrapper) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

func (kw *KeyboardWrapper) KeyPushed(key ebiten.Key) bool {
	return kw.keyState[key] == 1
}

// LastKeysPushed returns all keys that were pushed
// useful if you don't know what the user input might be for
// changing the key bindings.
// Careful using this looping through a map is slow
func (kw *KeyboardWrapper) LastKeysPushed() []ebiten.Key {
	keysPressed := []ebiten.Key{}
	for key, value := range kw.keyState {
		if value > 0 {
			keysPressed = append(keysPressed, key)
		}
	}
	return keysPressed
}

func (kw *KeyboardWrapper) Update() {
	for _, key := range kw.keys {
		if kw.IsKeyPressed(key) {
			kw.keyState[key]++
		} else {
			kw.keyState[key] = 0
		}
	}
}

// Parse converts a string version of a key to a ebiten.Key
func (kw *KeyboardWrapper) Parse(id string) (ebiten.Key, error) {
	switch strings.ToUpper(id) {
	case "0":
		return ebiten.Key0, nil
	case "1":
		return ebiten.Key1, nil
	case "2":
		return ebiten.Key2, nil
	case "3":
		return ebiten.Key3, nil
	case "4":
		return ebiten.Key4, nil
	case "5":
		return ebiten.Key5, nil
	case "6":
		return ebiten.Key6, nil
	case "7":
		return ebiten.Key7, nil
	case "8":
		return ebiten.Key8, nil
	case "9":
		return ebiten.Key9, nil
	case "A":
		return ebiten.KeyA, nil
	case "B":
		return ebiten.KeyB, nil
	case "C":
		return ebiten.KeyC, nil
	case "D":
		return ebiten.KeyD, nil
	case "E":
		return ebiten.KeyE, nil
	case "F":
		return ebiten.KeyF, nil
	case "G":
		return ebiten.KeyG, nil
	case "H":
		return ebiten.KeyH, nil
	case "I":
		return ebiten.KeyI, nil
	case "J":
		return ebiten.KeyJ, nil
	case "K":
		return ebiten.KeyK, nil
	case "L":
		return ebiten.KeyL, nil
	case "M":
		return ebiten.KeyM, nil
	case "N":
		return ebiten.KeyN, nil
	case "O":
		return ebiten.KeyO, nil
	case "P":
		return ebiten.KeyP, nil
	case "Q":
		return ebiten.KeyQ, nil
	case "R":
		return ebiten.KeyR, nil
	case "S":
		return ebiten.KeyS, nil
	case "T":
		return ebiten.KeyT, nil
	case "U":
		return ebiten.KeyU, nil
	case "V":
		return ebiten.KeyV, nil
	case "W":
		return ebiten.KeyW, nil
	case "X":
		return ebiten.KeyX, nil
	case "Y":
		return ebiten.KeyY, nil
	case "Z":
		return ebiten.KeyZ, nil
	case "ALT":
		return ebiten.KeyAlt, nil
	case "BS":
		return ebiten.KeyBackspace, nil
	case "CAPS":
		return ebiten.KeyCapsLock, nil
	case "COMMA":
		return ebiten.KeyComma, nil
	case "CTRL":
		return ebiten.KeyControl, nil
	case "DEL":
		return ebiten.KeyDelete, nil
	case "DOWN":
		return ebiten.KeyDown, nil
	case "END":
		return ebiten.KeyEnd, nil
	case "ENTER":
		return ebiten.KeyEnter, nil
	case "ESC":
		return ebiten.KeyEscape, nil
	case "F1":
		return ebiten.KeyF1, nil
	case "F2":
		return ebiten.KeyF2, nil
	case "F3":
		return ebiten.KeyF3, nil
	case "F4":
		return ebiten.KeyF4, nil
	case "F5":
		return ebiten.KeyF5, nil
	case "F6":
		return ebiten.KeyF6, nil
	case "F7":
		return ebiten.KeyF7, nil
	case "F8":
		return ebiten.KeyF8, nil
	case "F9":
		return ebiten.KeyF9, nil
	case "F10":
		return ebiten.KeyF10, nil
	case "F11":
		return ebiten.KeyF11, nil
	case "F12":
		return ebiten.KeyF12, nil
	case "HOME":
		return ebiten.KeyHome, nil
	case "INS":
		return ebiten.KeyInsert, nil
	case "LEFT":
		return ebiten.KeyLeft, nil
	case "PAGEDOWN":
		return ebiten.KeyPageDown, nil
	case "PAGEUP":
		return ebiten.KeyPageUp, nil
	case "PERIOD":
		return ebiten.KeyPeriod, nil
	case "RIGHT":
		return ebiten.KeyRight, nil
	case "SHIFT":
		return ebiten.KeyShift, nil
	case "SPACE":
		return ebiten.KeySpace, nil
	case "TAB":
		return ebiten.KeyTab, nil
	case "UP":
		return ebiten.KeyUp, nil
	default:
		return ebiten.Key(-1), errors.New("Not a valid key")
	}
}
