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
			newPart.Step(w)
			w.ApplyParticle(newPart)
		}
	}

	// send a light particle down
	for i := 0; i < 4; i++ {
		x := w.Rand.Intn(w.Width)
		w.ApplyParticle(Particle{
			Type:     LightParticle,
			Position: iv(x, 1),
			Velocity: Vec2f{0, 1},
		})
	}

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
	} else {
		p1x += 0.00001
	}
	if v1y > 0 {
		p1y = 1 - p1y
	} else {
		p1y += 0.00001
	}
	if v2x > 0 {
		p2x = 1 - p2x
	} else {
		p2x += 0.00001
	}
	if v2y > 0 {
		p2y = 1 - p2y
	} else {
		p2y += 0.00001
	}

	// calculate the time it will take for each to exit the pixel
	t1x := p1x / math.Abs(v1x)
	t1y := p1y / math.Abs(v1y)
	t2x := p2x / math.Abs(v2x)
	t2y := p2y / math.Abs(v2y)

	// whichever one is smaller is the answer
	t := math.Min(t1x, math.Min(t1y, math.Min(t2x, t2y)))
	p.Position.X += v1x * t
	p.Position.Y += v1y * t
	target.Position.X += v2x * t
	target.Position.Y += v2y * t

	w.ApplyParticle(target)
	w.ApplyParticle(p)
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
	// if particles leave out the top, they escape forever
	if p.Position.Y < 0 {
		return
	}
	// bottom is an impenetreble wall
	if p.Position.Y >= float64(w.Height) {
		p.Position.Y = float64(w.Height) - 1
		p.Velocity.Y = -p.Velocity.Y * ParticleClasses[p.Type].Elasticity
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
			w.Particles[destIndex].Absorb(p)
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
		Organic:      true,
		Energy:       ParticleClasses[ZygoteParticle].MaxEnergy,
		IntactDna:    dna,
		ExecutingDna: dna.PerfectClone(),
	}
	p.InitParamValues()
	w.Particles[w.Index(x, y)] = p

}
