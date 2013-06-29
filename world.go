package main

import (
	"math/rand"
)

type Vec2f struct {
	X float32
	Y float32
}

type Dna struct {
	Code  []byte
	Index int // position in code to execute next
}

type ParticleType int

var ParticleClasses = []ParticleClass{
	{"Carbon", 1, 1, 0xff374B65},
	{"Oxygen", 1, 1, 0xff94B4DD},
	{"Dirt",   1, 1, 0xff6B3000},
	{"Water",  1, 1, 0xff21009D},
	{"Light",  0, 0, 0xffFFF433},
	{"Chloro", 1, 1, 0xff0A7A00},
}

const (
	CarbonParticle ParticleType = iota
	OxygenParticle
	DirtParticle
	WaterParticle
	LightParticle
	ChloroParticle
)

type OrganicParticle struct {
	IntactDna    Dna // original DNA
	ExecutingDna Dna // starts as a copy of IntactDna
	BasicParticle
}

func (p *OrganicParticle) GetColor() uint32 {
	return ParticleClasses[p.Type].Color
}

type BasicParticle struct {
	Type ParticleType
	Position Vec2f
	Velocity Vec2f
}

func (p *BasicParticle) GetColor() uint32 {
	return ParticleClasses[p.Type].Color
}

type ParticleClass struct {
	Name    string
	Mass    float32
	Density float32
	Color   uint32
}

type Particle interface {
	GetColor() uint32
}

type World struct {
	Width     int
	Height    int
	Particles []Particle
	Time      int64
	Seed      int64
	Rand      *rand.Rand
}

func NewWorld(width int, height int, seed int64) *World {
	w := World{
		Width:     width,
		Height:    height,
		Particles: make([]Particle, width*height),
		Seed:      seed,
		Rand:      rand.New(rand.NewSource(seed)),
	}
	w.Rand.Seed(seed)

	waterTop := int(float32(height) * 0.5)
	dirtTop := int(float32(height) * 0.9)
	for y := waterTop; y < dirtTop; y++ {
		for x := 0; x < width; x++ {
			w.SetParticleAt(x, y, &BasicParticle{
				Type: WaterParticle,
				Position: iv(x,y),
			});
		}
	}
	for y := dirtTop; y < height; y++ {
		for x := 0; x < width; x++ {
			w.SetParticleAt(x, y, &BasicParticle{
				Type: DirtParticle,
				Position: iv(x,y),
			});
		}
	}

	return &w
}

func iv(x int, y int) Vec2f {
	return Vec2f{float32(x), float32(y)}
}

func (w *World) Step() {
	if w.Time % 20 == 0 {
		// send a light beam down
		x := w.Rand.Intn(w.Width)
		w.SetParticleAt(x, 0, &BasicParticle{
			Type: LightParticle,
			Position: iv(x, 0),
			Velocity: iv(0, 2),
		});
	}

	w.Time += 1
}

func (w *World) ParticleAt(x int, y int) Particle {
	return w.Particles[y*w.Width+x]
}

func (w *World) SetParticleAt(x int, y int, p Particle) {
	w.Particles[y * w.Width + x] = p
}
