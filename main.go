package main

import (
	"fmt"
	"github.com/scottferg/Go-SDL/sdl"
	"os"
	"time"
)

type View struct {
	OffsetX int
	OffsetY int
	Zoom    float32 // minimum 1 which is 1 pixel per particle
}

var world *World
var speed = 1
var zoom = 3

func main() {
	const width = 100
	const height = 100

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		fmt.Fprintf(os.Stderr, "SDL init error\n")
		return
	}
	screen := sdl.SetVideoMode(width * zoom, height * zoom, 32, sdl.DOUBLEBUF|sdl.HWSURFACE|sdl.HWACCEL)
	if screen == nil {
		fmt.Fprintf(os.Stderr, "SDL setvideomode error\n")
		return
	}

	sdl.WM_SetCaption("Planet Evo", "")

	seed := int64(1234)
	fmt.Printf("Using seed %d\n", seed)
	world = NewWorld(width, height, seed)

	go handleEvents()

	frame := make(chan int, 1)
	go func() {
		for {
			frame <- 1
			time.Sleep(16 * time.Millisecond)
		}
	}()

	for {
		for i := 0; i < speed; i++ {
			world.Step()
		}
		pix := &sdl.Rect{0, 0, uint16(zoom), uint16(zoom)}
		for y := 0; y < world.Height; y++ {
			pix.Y = int16(y * zoom)
			for x := 0; x < world.Width; x++ {
				pix.X = int16(x * zoom)
				screen.FillRect(pix, world.ColorAt(x, y))
			}
		}
		<-frame
		screen.Flip()
	}
}

func handleEvents() {
	for {
		event := <-sdl.Events
		switch e := event.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyboardEvent:
			if e.State == 1 {
				switch e.Keysym.Sym {
				case sdl.K_RIGHTBRACKET:
					if e.Keysym.Mod & sdl.KMOD_LSHIFT != 0 {
						speed *= 2
					} else {
						speed += 1
					}
					fmt.Fprintf(os.Stderr, "Speed: %d\n", speed)
				case sdl.K_LEFTBRACKET:
					if e.Keysym.Mod & sdl.KMOD_LSHIFT != 0 {
						speed /= 2
					} else {
						speed -= 1
					}
					if speed < 0 {
						speed = 0
					}
					fmt.Fprintf(os.Stderr, "Speed: %d\n", speed)
				}
			}
		case sdl.MouseButtonEvent:
			if e.Button == 1 && e.State == 1 {
				world.SpawnRandomCreature(int(e.X) / zoom, int(e.Y) / zoom)
			} else if e.Button == 2 && e.State == 1 {
				fmt.Fprintf(os.Stderr, "(%d, %d)\n", int(e.X) / zoom, int(e.Y) / zoom)
			}
		}
	}
}
