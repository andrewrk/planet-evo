package main

import (
	"fmt"
	"os"
	"math"
)

type ParticleType int

type ParticleClass struct {
	Name          string
	Mass          float64
	Density       float64
	BlockSunlight bool
	BlockAir      bool
	Color         uint32
	MaxEnergy     float64
	Elasticity    float64
	Friction      float64
	DeadType      ParticleType
}

var ParticleClasses = []ParticleClass{
	// non-organic particles
	{"Null", 0, 0, false, false, 0xff000000, 0, 0, 0, NullParticle},
	{"Carbon", 1, 0.2, false, true, 0xff374B65, 0, 0.99, 0, NullParticle},
	{"Oxygen", 1, 0.1, false, true, 0xff94B4DD, 0, 0.99, 0, NullParticle},
	{"Dirt", 10, 4, true, true, 0xff6B3000, 0, 0.1, 0.001, NullParticle},
	{"Water", 10, 1, false, true, 0xff21009D, 0, 0.7, 0.0001, NullParticle},
	{"Light", 0, 0, false, false, 0xffFFF433, 0, 0, 0, NullParticle},

	// organic particles
	{"Chloro", 4, 0.9, true, true, 0xff0A7A00, 5, 0.1, 0.001, DirtParticle},
	{"Fiber", 6, 0.8, true, true, 0xffB75900, 2, 0.5, 0.002, DirtParticle},
	{"Zygote", 5, 2, true, true, 0xffEFEFEF, 10, 0.7, 0.0005, DirtParticle},
}

const (
	// non-organic particles
	NullParticle ParticleType = iota
	CarbonParticle
	OxygenParticle
	DirtParticle
	WaterParticle
	LightParticle

	// organic particles
	ChloroParticle
	FiberParticle
	ZygoteParticle

	// meta
	ParticleCount
)

const FirstOrganicParticle = ChloroParticle
const OrganicParticleCount = ParticleCount - FirstOrganicParticle

type Particle struct {
	Type     ParticleType
	Position Vec2f
	Velocity Vec2f

	// only if particle is organic
	Dead         bool
	Energy       float64
	Age          int
	OrganismAge  int
	IntactDna    Dna // original DNA
	ExecutingDna Dna // starts as a copy of IntactDna
	ParamValues  [ParameterOpCodeCount]int
	RegisterX    int
	RegisterY    int
	Waiting      int // until this many steps are done, do nothing
}

func (p *Particle) Organic() bool {
	return p.Type >= FirstOrganicParticle
}

func (p *Particle) Color() uint32 {
	return ParticleClasses[p.Type].Color
}

func (p *Particle) InitParamValues() {
	for i := range p.ParamValues {
		p.ParamValues[i] = ParameterInfos[i].Default
	}
}

func (p *Particle) Step(w *World) {
	if p.Organic() && !p.Dead {
		p.StepDna(w)
		p.Age += 1
		p.OrganismAge += 1
		p.Energy -= 0.001
		if p.Energy <= 0 {
			fmt.Fprintf(os.Stderr, "Cell at %v ran out of energy.\n", p.Position)
			p.Die()
		}
	}
	// apply velocity
	p.Position.Add(&p.Velocity)
	// apply friction
	friction := ParticleClasses[p.Type].Friction
	if friction >= math.Abs(p.Velocity.X) {
		p.Velocity.X = 0
	} else {
		p.Velocity.X -= sign(p.Velocity.X) * friction
	}
	if friction >= math.Abs(p.Velocity.Y) {
		p.Velocity.Y = 0
	} else {
		p.Velocity.Y -= sign(p.Velocity.Y) * friction
	}
	p.CheckLowVelocity()
}

func (p *Particle) CheckLowVelocity() {
	const lowVelocityThreshold = 0.000001
	if p.Velocity.LengthSqrd() < lowVelocityThreshold {
		p.Velocity.Clear()
	}
}

func (p *Particle) Absorb(other Particle) {
	if p.Type == ChloroParticle && other.Type == LightParticle {
		fmt.Fprintf(os.Stderr, "Cell at %v +1 energy. now has %v\n", p.Position, p.Energy)
		p.GainEnergy(1)
	}
}

func (p *Particle) GainEnergy(amt int) {
	p.Energy += 1
	max := ParticleClasses[p.Type].MaxEnergy
	if p.Energy > max {
		p.Energy = max
	}
}

func (p *Particle) MoveToPixBorder(v Vec2f) {
	oldPos := p.Position

	// take the decimal part only so these values represent the position
	// within the pixel
	x := p.Position.X - math.Floor(p.Position.X)
	y := p.Position.Y - math.Floor(p.Position.Y)

	// adjust based on the velocities so that these values represent the
	// distance that must be traveled to exit the pixel
	if v.X > 0 {
		x = 1 - x
	}
	if v.Y > 0 {
		y = 1 - y
	}

	// calculate the time it will take for each to exit the pixel
	var tx, ty float64
	if v.X == 0 {
		tx = math.Inf(1)
	} else {
		tx = x / math.Abs(v.X)
	}
	if v.Y == 0 {
		ty = math.Inf(1)
	} else {
		ty = y / math.Abs(v.Y)
	}

	// whichever one is smaller is the answer
	t := math.Min(tx, ty)
	p.Position.X += v.X * t
	p.Position.Y += v.Y * t

	// run another frame if necessary
	if p.Position.FloorEql(&oldPos) {
		p.Position.Add(&v)
	}
}
