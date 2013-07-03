package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

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

	// introduce a bunch of carbon particles
	carbonParticleCount := w.Width * w.Height / 50
	for i := 0; i < carbonParticleCount; i++ {
		x := w.Rand.Intn(w.Width)
		y := w.Rand.Intn(w.Height)
		vx := 0.2 - w.Rand.Float64()*0.4
		vy := 0.2 - w.Rand.Float64()*0.4
		w.Particles[w.Index(x, y)] = Particle{
			Type:     CarbonParticle,
			Position: iv(x, y),
			Velocity: Vec2f{vx, vy},
		}
	}

	// add water
	waterParticleCount := w.Width * w.Height / 50
	for i := 0; i < waterParticleCount; i++ {
		x := w.Rand.Intn(w.Width)
		y := w.Rand.Intn(w.Height)
		w.Particles[w.Index(x, y)] = Particle{
			Type:     WaterParticle,
			Position: iv(x, y),
		}
	}

	// add dirt
	dirtParticleCount := w.Width * w.Height / 50
	for i := 0; i < dirtParticleCount; i++ {
		x := w.Rand.Intn(w.Width)
		y := w.Rand.Intn(w.Height)
		w.Particles[w.Index(x, y)] = Particle{
			Type:     DirtParticle,
			Position: iv(x, y),
		}
	}

	return &w
}

func iv(x int, y int) Vec2f {
	return Vec2f{float64(x) + 0.5, float64(y) + 0.5}
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
			newPart.Step(w)
			w.ApplyParticle(newPart)
		}
	}

	w.Time += 1
	w.Flip()
}

func (w *World) ResolveReplace(p Particle) {
	destIndex := w.AltIndexVec2f(p.Position)
	w.Particles[destIndex] = p
}

func (w *World) ResolveInsert(p Particle, v Vec2f) {
	p.MoveToPixBorder(v)
	p.Position.X = mod(p.Position.X, float64(w.Width))
	p.Position.Y = mod(p.Position.Y, float64(w.Height))

	index := w.AltIndexVec2f(p.Position)
	destPart := w.Particles[index]
	w.Particles[index] = p
	if destPart.Type == NullParticle {
		return
	}
	w.ResolveInsert(destPart, v)
}

func (w *World) ResolveCollide(p Particle, target Particle) {
	destIndex := w.AltIndexVec2f(p.Position)
	w.Particles[destIndex].Type = NullParticle

	m1 := ParticleClasses[p.Type].Mass
	m2 := ParticleClasses[target.Type].Mass
	e1 := ParticleClasses[p.Type].Elasticity
	e2 := ParticleClasses[target.Type].Elasticity
	v1 := p.Velocity
	v2 := target.Velocity

	v1x := e1 * (v1.X*(m1-m2) + 2*m2*v2.X) / (m1 + m2)
	v1y := e1 * (v1.Y*(m1-m2) + 2*m2*v2.Y) / (m1 + m2)
	v2x := e2 * (v2.X*(m2-m1) + 2*m1*v1.X) / (m1 + m2)
	v2y := e2 * (v2.Y*(m2-m1) + 2*m1*v1.Y) / (m1 + m2)

	p.Velocity.X = v1x
	p.Velocity.Y = v1y
	target.Velocity.X = v2x
	target.Velocity.Y = v2y

	p.CheckLowVelocity()
	target.CheckLowVelocity()

	if p.Velocity.LengthSqrd() == 0 && target.Velocity.LengthSqrd() == 0 {
		// the particles collided and then stopped moving. we need to
		// force the moving particle back from whence it came
		if v1.IsZero() {
			v2.Negate()
			w.ResolveInsert(target, v2)
			w.ApplyParticle(p)
		} else {
			v1.Negate()
			w.ApplyParticle(target)
			w.ResolveInsert(p, v1)
		}
		return
	}

	// based on the new velocities, calculate how much we must add to the
	// the positions until a particle exits this pixel
	p1x := p.Position.X
	p1y := p.Position.Y
	p2x := target.Position.X
	p2y := target.Position.Y

	// take the decimal part only so these values represent the position
	// within the pixel
	p1x -= math.Floor(p1x)
	p1y -= math.Floor(p1y)
	p2x -= math.Floor(p2x)
	p2y -= math.Floor(p2y)

	// adjust based on the velocities so that these values represent the
	// distance that must be traveled to exit the pixel
	if v1x > 0 {
		p1x = 1 - p1x
	}
	if v1y > 0 {
		p1y = 1 - p1y
	}
	if v2x > 0 {
		p2x = 1 - p2x
	}
	if v2y > 0 {
		p2y = 1 - p2y
	}

	// calculate the time it will take for each to exit the pixel
	var t1x, t1y, t2x, t2y float64
	if v1x == 0 {
		t1x = math.Inf(1)
	} else {
		t1x = p1x / math.Abs(v1x)
	}
	if v1y == 0 {
		t1y = math.Inf(1)
	} else {
		t1y = p1y / math.Abs(v1y)
	}
	if v2x == 0 {
		t2x = math.Inf(1)
	} else {
		t2x = p2x / math.Abs(v2x)
	}
	if v2y == 0 {
		t2y = math.Inf(1)
	} else {
		t2y = p2y / math.Abs(v2y)
	}

	// whichever one is smaller is the answer
	t := math.Min(t1x, math.Min(t1y, math.Min(t2x, t2y)))
	p.Position.X += v1x * t
	p.Position.Y += v1y * t
	target.Position.X += v2x * t
	target.Position.Y += v2y * t

	// run another frame if necessary
	if p.Position.FloorEql(&target.Position) {
		p.Position.Add(&p.Velocity)
		target.Position.Add(&target.Velocity)
	}

	w.ApplyParticle(target)
	w.ApplyParticle(p)
}

// modulus the way God intended it.
func mod(x float64, y float64) float64 {
	val := math.Mod(x, y)
	if val < 0 {
		return val + y
	}
	return val
}

// handle particle collisions
func (w *World) ApplyParticle(p Particle) {
	// wrap X and Y
	p.Position.X = mod(p.Position.X, float64(w.Width))
	p.Position.Y = mod(p.Position.Y, float64(w.Height))

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
			w.Particles[destIndex].Absorb(p)
		} else {
			// move the non-light particle in the negative direction of the light
			panic("light")
		}
	case CarbonParticle:
		switch {
		case destPart.Type == LightParticle:
			panic("light")
		case ParticleClasses[destPart.Type].BlockAir:
			w.ResolveCollide(p, destPart)
		default:
			panic("carbon")
		}
	case WaterParticle:
		switch {
		case destPart.Type == LightParticle:
			panic("displace")
		default:
			w.ResolveCollide(p, destPart)
		}
	case DirtParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveReplace(p)
			w.Particles[destIndex].Absorb(destPart)
		default:
			w.ResolveCollide(p, destPart)
		}
	case FiberParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveReplace(p)
			w.Particles[destIndex].Absorb(destPart)
		default:
			w.ResolveCollide(p, destPart)
		}
	case ChloroParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveReplace(p)
			w.Particles[destIndex].Absorb(destPart)
		default:
			w.ResolveCollide(p, destPart)
		}
	case ZygoteParticle:
		switch {
		case destPart.Type == LightParticle:
			w.ResolveReplace(p)
			w.Particles[destIndex].Absorb(destPart)
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

func (w *World) SpawnRandomCreature(x int, y int) {
	//dna := w.CreateRandomDna()
	dna := w.CreateSingleCelledPlant()
	p := Particle{
		Type:         ZygoteParticle,
		Position:     iv(x, y),
		Energy:       ParticleClasses[ZygoteParticle].MaxEnergy,
		IntactDna:    dna,
		ExecutingDna: dna.PerfectClone(),
	}
	p.InitParamValues()
	w.Particles[w.Index(x, y)] = p

}
