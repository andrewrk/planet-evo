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

func main() {
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		fmt.Fprintf(os.Stderr, "SDL init error\n")
		return
	}
	screen := sdl.SetVideoMode(640, 480, 32, sdl.DOUBLEBUF|sdl.HWSURFACE|sdl.HWACCEL)
	if screen == nil {
		fmt.Fprintf(os.Stderr, "SDL setvideomode error\n")
		return
	}

	sdl.WM_SetCaption("Planet Evo", "")

	seed := int64(1234)
	fmt.Printf("Using seed %d\n", seed)
	world = NewWorld(640, 480, seed)

	go handleEvents()

	frame := make(chan int, 1)
	go func() {
		for {
			frame <- 1
			time.Sleep(16 * time.Millisecond)
		}
	}()

	for {
		world.Step()
		pix := &sdl.Rect{0, 0, 1, 1}
		for y := 0; y < world.Height; y++ {
			pix.Y = int16(y)
			for x := 0; x < world.Width; x++ {
				pix.X = int16(x)
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
		case sdl.MouseButtonEvent:
			if e.Button == 1 && e.State == 1 {
				world.SpawnRandomCreature(int(e.X), int(e.Y))
			}
		}
	}
}
