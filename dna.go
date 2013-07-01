package main

import (
	"fmt"
	"math"
)

type Instruction struct {
	OpCode byte
	Value byte
}

type Dna struct {
	Instructions []Instruction
	Index        int // position in code to execute next
}

type DnaOp int
type ParameterType int

const (
	ValueSourceParam ParameterType = iota
	RegisterParam
	ComparisonParam
	OperationParam
	CodeLabelParam
	Direction8Param
	NewCellTypeParam
	BlockOrContinueParam
	NumberParam
	BooleanParam
)

var ParamCaps = []int{
	int(ValueSourceParamCount),
	int(RegisterParamCount),
	int(ComparisonParamCount),
	int(OperationParamCount),
	65535,
	8,
	int(OrganicParticleCount) + 1,
	2,
	256,
	2,
}

const (
	// core instructions
	NoOp DnaOp = iota
	CellDivisionOp
	CellDeathOp
	JumpOp
	WaitOp
	UpdateRegisterOp
	CalcOp
	ModifyDnaOp

	// parameter setting
	CellDivisionEnergyForNewCellOp
	CellDivisionDirectionOp
	CellDivisionNewCellTypeOp
	CellDivisionForkLabelOp
	CellDivisionDoWeForkOp
	CellDivisionContingencyPlanOp
	JumpOperandLeftOp
	JumpOperandRightOp
	JumpComparisonOp
	JumpLabelOp
	WaitSourceOp
	UpdateRegisterSourceOp
	UpdateRegisterDestOp
	CalcOperandLeftOp
	CalcOperandRightOp
	CalcOperationOp
	CalcDestOp
	ModifyDnaLabelOp
	ModifyDnaSourceOp
	ValueSourceDirectionOp
	ValueSourceNumberOp
	ValueSourceLabelOp
	ProgramEndBehaviorOp

	// meta
	DnaOpCount
)
const ParameterOpCodeStart = CellDivisionEnergyForNewCellOp
const ParameterOpCodeCount = DnaOpCount - ParameterOpCodeStart

type ParameterInfo struct {
	Type ParameterType
	Default int
}

const (
	Direction8N int = iota
	Direction8S
	Direction8W
	Direction8E
	Direction8NW
	Direction8NE
	Direction8SW
	Direction8SE
)

const (
	BorCBlock int = iota
	BorCContinue
)

const (
	ValueSourceRegisterX int = iota
	ValueSourceRegisterY
	ValueSourceCellAge
	ValueSourceCellType
	ValueSourceCellEnergy
	ValueSourceOrganismAge
	ValueSource0
	ValueSource1
	ValueSource2
	ValueSource3
	ValueSource4
	ValueSource5
	ValueSourceNeighborParticleType
	ValueSourceNeighborIsSameCellType
	ValueSourceNumber
	ValueSourceLabel

	// meta
	ValueSourceParamCount
)

const (
	RegisterNone int = iota
	RegisterX
	RegisterY

	// meta
	RegisterParamCount
)

const (
	ComparisonNever int = iota
	ComparisonAlways
	ComparisonLeftZero
	ComparisonLeftNonZero
	ComparisonRightZero
	ComparisonRightNonZero
	ComparisonLeftGTRight
	ComparisonLeftLTRight
	ComparisonLeftEQRight
	ComparisonLeftNERight
	ComparisonLeftGERight
	ComparisonLeftLERight

	// meta
	ComparisonParamCount
)

const (
	OperationLeft int = iota
	Operation0
	Operation1
	Operation2
	OperationLeftPlus1
	OperationLeftPlus2
	OperationRight
	OperationRightPlus1
	OperationRightPlus2
	OperationLeftPlusRight
	OperationLeftMinusRight
	OperationRightMinusLeft
	OperationNegLeftNegRight
	OperationLeftDivRight
	OperationRightDivLeft
	OperationLeftMult2
	OperationRightMult2
	OperationLeftModRight
	OperationRightModLeft
	OperationMin
	OperationMax
	OperationLeftMultRight

	// meta
	OperationParamCount
)

var ParameterInfos = []ParameterInfo{
	{ValueSourceParam, 5}, //CellDivisionEnergyForNewCellOp
	{Direction8Param, Direction8N},
	{NewCellTypeParam, 0},
	{CodeLabelParam, 0},
	{BooleanParam, 0}, // CellDivisionDoWeForkOp
	{BlockOrContinueParam, BorCBlock},
	{ValueSourceParam, ValueSourceRegisterX}, // JumpOperandLeftOp
	{ValueSourceParam, ValueSourceRegisterX},
	{ComparisonParam, ComparisonLeftNonZero},
	{CodeLabelParam, 0},
	{ValueSourceParam, ValueSource2},
	{ValueSourceParam, ValueSourceRegisterX},
	{RegisterParam, RegisterNone},
	{ValueSourceParam, ValueSourceRegisterX},
	{ValueSourceParam, ValueSourceRegisterX},
	{OperationParam, OperationLeft},
	{RegisterParam, RegisterNone},
	{CodeLabelParam, 1},
	{ValueSourceParam, ValueSourceRegisterX},
	{Direction8Param, Direction8N},
	{NumberParam, 0},
	{CodeLabelParam, 1},
	{BlockOrContinueParam, BorCContinue},
}

func getXYForDir(dir int) (x int, y int) {
	switch dir {
	case Direction8N:
		return 0, -1
	case Direction8S:
		return 0, 1
	case Direction8W:
		return -1, 0
	case Direction8E:
		return 1, 0
	case Direction8NW:
		return -1, 1
	case Direction8NE:
		return 1, 1
	case Direction8SW:
		return -1, -1
	case Direction8SE:
		return 1, -1
	}
	panic("unknown direction")
}

func (w *World) CreateRandomDna() Dna {
	const instrCount = 10
	dna := Dna{
		Index: 0,
		Instructions: make([]Instruction, instrCount),
	}
	for i := range dna.Instructions {
		dna.Instructions[i] = w.CreateRandomInstruction()
	}
	return dna
}

func (w *World) CreateRandomInstruction() Instruction {
	return Instruction{
		byte(w.Rand.Intn(256)),
		byte(w.Rand.Intn(256)),
	}
}

func (dna *Dna) Clone() Dna {
	newDna := *dna
	newDna.Instructions = make([]Instruction, len(dna.Instructions))
	copy(newDna.Instructions, dna.Instructions)
	return newDna
}

func (p *Particle) getParamValByOp(op DnaOp) int {
	paramIndex := op - ParameterOpCodeStart
	val := p.ParamValues[paramIndex]
	return val % ParamCaps[ParameterInfos[paramIndex].Type]
}

func (p *Particle) getValueSource(op DnaOp, w *World) int {
	switch p.getParamValByOp(op) {
	case ValueSourceRegisterX:
		return p.RegisterX
	case ValueSourceRegisterY:
		return p.RegisterY
	case ValueSourceCellAge:
		return p.Age
	case ValueSourceCellType:
		return int(p.Type)
	case ValueSourceCellEnergy:
		return int(math.Floor(p.Energy))
	case ValueSourceOrganismAge:
		return p.OrganismAge
	case ValueSource0:
		return 0
	case ValueSource1:
		return 1
	case ValueSource2:
		return 2
	case ValueSource3:
		return 3
	case ValueSource4:
		return 4
	case ValueSource5:
		return 5
	case ValueSourceNeighborParticleType:
		dir := p.getParamValByOp(ValueSourceDirectionOp)
		x, y := getXYForDir(dir)
		return int(w.Particles[w.Index(x, y)].Type)
	case ValueSourceNeighborIsSameCellType:
		dir := p.getParamValByOp(ValueSourceDirectionOp)
		x, y := getXYForDir(dir)
		if w.Particles[w.Index(x, y)].Type == p.Type {
			return 1
		} else {
			return 0
		}
	case ValueSourceNumber:
		return p.getParamValByOp(ValueSourceNumberOp)
	case ValueSourceLabel:
		return p.getParamValByOp(ValueSourceLabelOp)
	}
	panic("unrecognized ValueSource")
}

func (p *Particle) performComparison(op DnaOp, left int, right int) bool {
	switch p.getParamValByOp(op) {
	case ComparisonNever:
		return false
	case ComparisonAlways:
		return true
	case ComparisonLeftZero:
		return left == 0
	case ComparisonLeftNonZero:
		return left != 0
	case ComparisonRightZero:
		return right == 0
	case ComparisonRightNonZero:
		return right != 0
	case ComparisonLeftGTRight:
		return left > right
	case ComparisonLeftLTRight:
		return left < right
	case ComparisonLeftEQRight:
		return left == right
	case ComparisonLeftNERight:
		return left != right
	case ComparisonLeftGERight:
		return left >= right
	case ComparisonLeftLERight:
		return left <= right
	}
	panic("unrecognized Comparison")
}

func (p *Particle) Die() {
	p.Energy = 0
	p.Dead = true
}

func (p *Particle) getNewCellType(op DnaOp) ParticleType {
	newCellEnumValue := p.getParamValByOp(op)
	if newCellEnumValue == 0 {
		return NullParticle
	}
	return ParticleType(newCellEnumValue) - 1 + FirstOrganicParticle
}

func (p *Particle) Split(w *World, dx int, dy int, newCellType ParticleType, pc int, energy float64) {
	x := int(math.Floor(p.Position.X)) + dx
	y := int(math.Floor(p.Position.Y)) + dy
	// how is babby formed
	baby := Particle{
		Type: newCellType,
		Position: iv(x, y),
		Organic: true,
		Energy: energy,
		IntactDna: p.IntactDna.Clone(),
		ExecutingDna: p.ExecutingDna.Clone(),
	}
	baby.InitParamValues()
	baby.ExecutingDna.Index = pc
	w.ApplyParticle(baby)
}

func (p *Particle) StepDna(w *World) {
	if !p.Organic || p.Dead {
		return
	}
	pc := p.ExecutingDna.Index
	if pc == -1 {
		// DNA program is permanently halted
		return
	}
	instr := p.ExecutingDna.Instructions[pc]
	pc += 1
	var opCode DnaOp = DnaOp(instr.OpCode) % DnaOpCount
	if opCode < ParameterOpCodeStart {
		switch opCode {
		default:
			panic(fmt.Sprintf("unhandled op code: %d", opCode))
		case NoOp:
			// done. that was easy.
		case CellDivisionOp:
			energyRequired := float64(p.getValueSource(CellDivisionEnergyForNewCellOp, w))
			if p.Energy >= energyRequired {
				p.Energy -= energyRequired
				dir := p.getParamValByOp(CellDivisionDirectionOp)
				x, y := getXYForDir(dir)
				newCellType := p.getNewCellType(CellDivisionNewCellTypeOp)
				doWeFork := p.getParamValByOp(CellDivisionDoWeForkOp)
				newCellPc := pc
				if doWeFork == 1 {
					newCellPc = p.getParamValByOp(CellDivisionForkLabelOp)
				}
				p.Split(w, x, y, newCellType, newCellPc, energyRequired)
			} else {
				plan := p.getParamValByOp(CellDivisionContingencyPlanOp)
				switch plan {
				default:
					panic("unrecognized contingency plan")
				case BorCBlock:
					pc -= 1
				case BorCContinue:
					break
				}
			}
		case CellDeathOp:
			p.Die()
		case JumpOp:
			left := p.getValueSource(JumpOperandLeftOp, w)
			right := p.getValueSource(JumpOperandRightOp, w)
			comp := p.performComparison(JumpComparisonOp, left, right)
			newAddr := p.getParamValByOp(JumpLabelOp)
			if comp {
				pc = newAddr
			}
		//case WaitOp:
		//case UpdateRegisterOp:
		//case CalcOp:
		//case ModifyDnaOp:
		}
	} else {
		// set a parameter value
		paramIndex := opCode - ParameterOpCodeStart
		paramType := ParameterInfos[paramIndex].Type
		valueCap := ParamCaps[paramType]
		newValue := int(instr.Value) % valueCap
		if paramType == CodeLabelParam {
			newValue += pc - 1
		}
		p.ParamValues[paramIndex] = newValue

	}
	if pc >= len(p.ExecutingDna.Instructions) {
		switch p.ParamValues[ProgramEndBehaviorOp - ParameterOpCodeStart] {
		default:
			panic("unknown BlockOrContinue value")
		case BorCBlock:
			pc = -1
		case BorCContinue:
			pc = 0
		}
	}
	p.ExecutingDna.Index = pc
	p.Age += 1
	p.OrganismAge += 1
}
