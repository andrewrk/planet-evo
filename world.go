package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Vec2f struct {
	X float64
	Y float64
}

func (v *Vec2f) Add(other *Vec2f) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2f) Subtract(other *Vec2f) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vec2f) Clear() {
	v.X = 0
	v.Y = 0
}

func (v *Vec2f) Negate() {
	v.X = -v.X
	v.Y = -v.Y
}

func (v *Vec2f) Scale(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2f) FloorEql(other *Vec2f) bool {
	return math.Floor(v.X) == math.Floor(other.X) && math.Floor(v.Y) == math.Floor(other.Y)
}

type Dna struct {
	Instructions []uint16
	Index        int // position in code to execute next
}

type ParticleType int

type ParticleClass struct {
	Name          string
	Mass          float64
	Density       float64
	BlockSunlight bool
	BlockAir      bool
	Color         uint32
	MaxEnergy     float64
}

var ParticleClasses = []ParticleClass{
	// non-organic particles
	{"Null", 0, 0, false, false, 0xff000000, 0},
	{"Carbon", 1, 1, false, true, 0xff374B65, 0},
	{"Oxygen", 1, 1, false, true, 0xff94B4DD, 0},
	{"Dirt", 10, 1, true, true, 0xff6B3000, 0},
	{"Water", 10, 1, false, true, 0xff21009D, 0},
	{"Light", 0, 0, false, false, 0xffFFF433, 0},

	// organic particles
	{"Chloro", 4, 1, true, true, 0xff0A7A00, 5},
	{"Fiber", 6, 1, true, true, 0xffB75900, 2},
}

const (
	NullParticle ParticleType = iota
	CarbonParticle
	OxygenParticle
	DirtParticle
	WaterParticle
	LightParticle
	ChloroParticle
	FiberParticle
)

type Particle struct {
	Type     ParticleType
	Position Vec2f
	Velocity Vec2f

	Organic      bool
	IntactDna    Dna // original DNA
	ExecutingDna Dna // starts as a copy of IntactDna
}

func (p *Particle) Color() uint32 {
	return ParticleClasses[p.Type].Color
}

type World struct {
	Width     int
	Height    int
	Particles []Particle
	AltOffset int
	Time      int64
	Seed      int64
	Rand      *rand.Rand
}

func NewWorld(width int, height int, seed int64) *World {
	w := World{
		Width:     width,
		Height:    height,
		Particles: make([]Particle, width*height*2),
		Seed:      seed,
		Rand:      rand.New(rand.NewSource(seed)),
	}
	w.Rand.Seed(seed)

	waterTop := int(float64(height) * 0.5)
	dirtTop := int(float64(height) * 0.9)

	// introduce a bunch of carbon particles
	carbonParticleCount := waterTop * w.Width / 50
	for i := 0; i < carbonParticleCount; i++ {
		x := w.Rand.Intn(w.Width)
		y := w.Rand.Intn(waterTop)
		vx := 0.2 - w.Rand.Float64()*0.4
		vy := 0.2 - w.Rand.Float64()*0.4
		w.Particles[w.Index(x, y)] = Particle{
			Type:     CarbonParticle,
			Position: iv(x, y),
			Velocity: Vec2f{vx, vy},
		}
	}

	// add water
	for y := waterTop; y < dirtTop; y++ {
		for x := 0; x < width; x++ {
			w.Particles[w.Index(x, y)] = Particle{
				Type:     WaterParticle,
				Position: iv(x, y),
			}
		}
	}

	// add dirt
	for y := dirtTop; y < height; y++ {
		for x := 0; x < width; x++ {
			w.Particles[w.Index(x, y)] = Particle{
				Type:     DirtParticle,
				Position: iv(x, y),
			}
		}
	}

	return &w
}

func iv(x int, y int) Vec2f {
	return Vec2f{float64(x), float64(y)}
}

func (w *World) Step() {
	w.ClearAlt()

	// apply velocity to particles
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			oldIndex := w.Index(x, y)
			sourcePart := w.Particles[oldIndex]
			if sourcePart.Type == NullParticle {
				continue
			}
			newPart := sourcePart
			newPart.Position.Add(&newPart.Velocity)
			w.ApplyParticle(newPart)
		}
	}

	// send a light particle down
	x := w.Rand.Intn(w.Width)
	w.ApplyParticle(Particle{
		Type:     LightParticle,
		Position: iv(x, 1),
		Velocity: Vec2f{0, 1},
	})

	w.Time += 1
	w.Flip()
}

func (w *World) ResolveDisplace(p Particle, target Particle) {
	destIndex := w.AltIndexVec2f(p.Position)
	w.Particles[destIndex] = p
	target.Position.Subtract(&p.Velocity)
	if target.Position.Y >= 0 && target.Position.Y < float64(w.Height) {
		w.ApplyParticle(target)
	} else {
		fmt.Fprintf(os.Stderr, "Error resolving displace at %v\n", p.Position)
	}
}

func (w *World) ResolveReplace(p Particle) {
	destIndex := w.AltIndexVec2f(p.Position)
	w.Particles[destIndex] = p
}

func (w *World) ResolveCollide(p Particle, target Particle) {
	pos := p.Position
	destIndex := w.AltIndexVec2f(p.Position)
	w.Particles[destIndex].Type = NullParticle

	m1 := ParticleClasses[p.Type].Mass
	m2 := ParticleClasses[target.Type].Mass
	v1 := p.Velocity
	v2 := target.Velocity

	v1x := (v1.X*(m1-m2) + 2*m2*v2.X) / (m1 + m2)
	v1y := (v1.Y*(m1-m2) + 2*m2*v2.Y) / (m1 + m2)
	v2x := (v2.X*(m2-m1) + 2*m1*v1.X) / (m1 + m2)
	v2y := (v2.Y*(m2-m1) + 2*m1*v1.Y) / (m1 + m2)

	p.Velocity.X = v1x
	p.Velocity.Y = v1y
	target.Velocity.X = v2x
	target.Velocity.Y = v2y

	const maxCycles = 1000
	for i := 0; i < maxCycles; i++ {
		if !p.Position.FloorEql(&target.Position) {
			w.ApplyParticle(target)
			w.ApplyParticle(p)
			return
		}
		p.Position.Add(&p.Velocity)
		target.Position.Add(&target.Velocity)
	}
	fmt.Fprintf(os.Stderr, "Error resolving collision at %v\n", pos)
}

// handle particle collisions
func (w *World) ApplyParticle(p Particle) {
	// wrap X
	if p.Position.X >= float64(w.Width) {
		p.Position.X -= float64(w.Width)
	}
	if p.Position.X < 0 {
		p.Position.X += float64(w.Width)
	}
	// bounce off top
	if p.Position.Y < 0 {
		p.Position.Y = 0
		p.Velocity.Y = -p.Velocity.Y
	}
	// bounce off bottom
	if p.Position.Y >= float64(w.Height) {
		p.Position.Y = float64(w.Height) - 1
		p.Velocity.Y = -p.Velocity.Y
	}

	destIndex := w.AltIndexVec2f(p.Position)
	destPart := w.Particles[destIndex]

	// no particle in destination point; simply place the particle there.
	if destPart.Type == NullParticle {
		w.Particles[destIndex] = p
		return
	}

	switch p.Type {
	case LightParticle:
		if ParticleClasses[destPart.Type].BlockSunlight {
			// non-light particle wins
		} else {
			// move the non-light particle in the negative direction of the light
			w.ResolveDisplace(p, destPart)
		}
	case CarbonParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveDisplace(destPart, p)
		case ParticleClasses[destPart.Type].BlockAir:
			w.ResolveCollide(p, destPart)
		default:
			w.ResolveDisplace(p, destPart)
		}
	case WaterParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveDisplace(destPart, p)
		default:
			w.ResolveCollide(p, destPart)
		}
	case DirtParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveReplace(p)
		default:
			w.ResolveCollide(p, destPart)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unhandled particle placement: %+v\n", p)
	}
}

func (w *World) Flip() {
	w.AltOffset = w.NextAltOffset()
}

func (w *World) ClearAlt() {
	start := w.NextAltOffset()
	end := start + w.Width*w.Height
	for i := start; i < end; i++ {
		w.Particles[i].Type = NullParticle
	}
}

func (w *World) NextAltOffset() int {
	if w.AltOffset == 0 {
		return w.Width * w.Height
	}
	return 0
}

func (w *World) Index(x int, y int) int {
	return w.AltOffset + y*w.Width + x
}

func (w *World) AltIndex(x int, y int) int {
	return w.NextAltOffset() + y*w.Width + x
}

func (w *World) IndexVec2f(v Vec2f) int {
	return w.Index(int(v.X), int(v.Y))
}

func (w *World) AltIndexVec2f(v Vec2f) int {
	return w.AltIndex(int(v.X), int(v.Y))
}

func (w *World) ColorAt(x int, y int) uint32 {
	p := w.Particles[w.Index(x, y)]
	return ParticleClasses[p.Type].Color
}
