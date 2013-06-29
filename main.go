package main

import (
	"fmt"
	"github.com/scottferg/Go-SDL/sdl"
	"os"
)

const (
	EmptyColor = 0xff000000
)

type View struct {
	OffsetX int
	OffsetY int
	Zoom    float32 // minimum 1 which is 1 pixel per particle
}

func main() {
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		fmt.Fprintf(os.Stderr, "SDL init error\n")
		return
	}
	screen := sdl.SetVideoMode(1000, 600, 32, sdl.DOUBLEBUF|sdl.HWSURFACE|sdl.HWACCEL)
	if screen == nil {
		fmt.Fprintf(os.Stderr, "SDL setvideomode error\n")
		return
	}

	sdl.WM_SetCaption("Planet Evo", "")

	w := NewWorld(1000, 600)

	for {
		event := <-sdl.Events
		switch event.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		}

		w.Step()
		pix := &sdl.Rect{0, 0, 1, 1}
		for y := 0; y < w.Height; y++ {
			pix.Y = int16(y)
			for x := 0; x < w.Width; x++ {
				pix.X = int16(x)
				particle := w.ParticleAt(x, y)
				if particle == nil {
					screen.FillRect(pix, EmptyColor)
				} else {
					screen.FillRect(pix, particle.GetColor())
				}
			}
		}
		screen.Flip()
	}
}
