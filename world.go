package main

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
}

func NewWorld(width int, height int) *World {
	w := World{
		Width:     width,
		Height:    height,
		Particles: make([]Particle, width*height),
	}
	waterTop := int(float32(height) * 0.5)
	dirtTop := int(float32(height) * 0.9)
	for y := waterTop; y < dirtTop; y++ {
		for x := 0; x < width; x++ {
			w.SetParticleAt(x, y, &BasicParticle{
				Type: WaterParticle,
				Position: Vec2f{float32(x), float32(y)},
			});
		}
	}
	for y := dirtTop; y < height; y++ {
		for x := 0; x < width; x++ {
			w.SetParticleAt(x, y, &BasicParticle{
				Type: DirtParticle,
				Position: Vec2f{float32(x), float32(y)},
			});
		}
	}

	return &w
}

func (w *World) Step() {

}

func (w *World) ParticleAt(x int, y int) Particle {
	return w.Particles[y*w.Width+x]
}

func (w *World) SetParticleAt(x int, y int, p Particle) {
	w.Particles[y * w.Width + x] = p
}
