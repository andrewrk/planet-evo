package main

import (
	"fmt"
	"os"
	"github.com/scottferg/Go-SDL/sdl"
)

func main() {
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		fmt.Fprintf(os.Stderr, "SDL init error\n")
		return
	}
	screen := sdl.SetVideoMode(800, 900, 32, 0)
	if screen == nil {
		fmt.Fprintf(os.Stderr, "SDL setvideomode error\n")
		return
	}

	sdl.WM_SetCaption("Planet Evo", "")

	for {
		event := <-sdl.Events
		switch event.(type) {
		case sdl.QuitEvent:
			os.Exit(0);
		}
	}
}

