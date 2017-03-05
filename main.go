package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten"
	"github.com/ranchblt/labyrinthofthechimera/labyrinth"
)

// Version is autoset from the build script
var Version string

// Build is autoset from the build script
var Build string

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var debug = flag.Bool("debug", false, "Turns on debug lines and debug messaging")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	game := labyrinth.NewGame(debug)
	update := game.Update
	if err := ebiten.Run(update, labyrinth.ScreenWidth, labyrinth.ScreenHeight, 1, "Labrinth of the Chimera "+Version+" "+Build); err != nil {
		log.Fatal(err)
	}
}
