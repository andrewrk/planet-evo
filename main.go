package main

import (
	"fmt"
	"github.com/scottferg/Go-SDL/sdl"
	"os"
)

type Dna struct {
	Code  []byte
	Index int // position in code to execute next
}

type ParticleType int

const (
	NoParticle ParticleType = iota
)

const (
	SkyColor = 0x94B4DDff
)

type OrganicParticle struct {
	IntactDna    Dna // original DNA
	ExecutingDna Dna // starts as a copy of IntactDna
	BasicParticle
}

func (p *OrganicParticle) GetColor() uint32 {
	return p.Color
}

type BasicParticle struct {
	Mass    float32
	Density float32
	Color   uint32
}

func (p *BasicParticle) GetColor() uint32 {
	return p.Color
}

type Particle interface {
	GetColor() uint32
}

type World struct {
	Width     int
	Height    int
	Particles []Particle
	Time      int64
}

type View struct {
	OffsetX int
	OffsetY int
	Zoom    float32 // minimum 1 which is 1 pixel per particle
}

func NewWorld(width int, height int) *World {
	w := World{
		Width:     width,
		Height:    height,
		Particles: make([]Particle, width*height),
	}
	return &w
}

func (w *World) Step() {

}

func (w *World) ParticleAt(x int, y int) Particle {
	return w.Particles[y*w.Width+x]
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
					screen.FillRect(pix, SkyColor)
				} else {
					screen.FillRect(pix, particle.GetColor())
				}
			}
		}
		screen.Flip()
	}
}
