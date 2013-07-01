package main

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
	{"Zygote", 5, 1, true, true, 0xffEFEFEF, 10},
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
	Organic      bool

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
}


func (p *Particle) Color() uint32 {
	return ParticleClasses[p.Type].Color
}

func (p *Particle) InitParamValues() {
	for i := range p.ParamValues {
		p.ParamValues[i] = ParameterInfos[i].Default
	}
}
