package main

import (
	"fmt"
	"os"
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
}

var ParticleClasses = []ParticleClass{
	// non-organic particles
	{"Null", 0, 0, false, false, 0xff000000, 0, 0},
	{"Carbon", 1, 0.2, false, true, 0xff374B65, 0, 0.99},
	{"Oxygen", 1, 0.1, false, true, 0xff94B4DD, 0, 0.99},
	{"Dirt", 10, 4, true, true, 0xff6B3000, 0, 0.1},
	{"Water", 10, 1, false, true, 0xff21009D, 0, 0.9},
	{"Light", 0, 0, false, false, 0xffFFF433, 0, 0},

	// organic particles
	{"Chloro", 4, 0.9, true, true, 0xff0A7A00, 5, 0.1},
	{"Fiber", 6, 0.8, true, true, 0xffB75900, 2, 0.5},
	{"Zygote", 5, 2, true, true, 0xffEFEFEF, 10, 0.7},
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
	Organic  bool

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

func (p *Particle) Color() uint32 {
	return ParticleClasses[p.Type].Color
}

func (p *Particle) InitParamValues() {
	for i := range p.ParamValues {
		p.ParamValues[i] = ParameterInfos[i].Default
	}
}

func (p *Particle) Step(w *World) {
	if p.Organic && !p.Dead {
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

	// apply gravity
	//const gravityConstant = 0.01
	//gravity := gravityConstant
	//// apply upward force depending on mass and density
	//density := ParticleClasses[p.Type].Density
	//mass := ParticleClasses[p.Type].Mass
	//md := mass * density
	//if md > 0 && md < 1 {
	//	gravity *= md
	//}
	//p.Velocity.Y += gravity
	//if p.Velocity.Y > 1 {
	//	// can't go faster than the speed of light!
	//	p.Velocity.Y = 1
	//}
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
