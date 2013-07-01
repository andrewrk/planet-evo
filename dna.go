package main

import (
	"fmt"
	"math"
	"errors"
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
)

var ParamCaps = []int{
	int(ValueSourceParamCount),
	int(RegisterParamCount),
	int(ComparisonParamCount),
	int(OperationParamCount),
	256,
	8,
	int(OrganicParticleCount) + 1,
	2,
	256,
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
	ValueSourceNone int = iota
	ValueSourceRegisterX
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
	{ValueSourceParam, 5},
	{Direction8Param, Direction8N},
	{NewCellTypeParam, 0},
	{CodeLabelParam, 0},
	{BlockOrContinueParam, BorCBlock},
	{ValueSourceParam, ValueSourceRegisterX},
	{ValueSourceParam, ValueSourceRegisterX},
	{ComparisonParam, ComparisonLeftNonZero},
	{CodeLabelParam, 0},
	{ValueSourceParam, ValueSource2},
	{ValueSourceParam, ValueSourceNone},
	{RegisterParam, RegisterNone},
	{ValueSourceParam, ValueSourceNone},
	{ValueSourceParam, ValueSourceNone},
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
	return p.ParamValues[paramIndex]
}

func (p *Particle) getValueSource(op DnaOp, w *World) (val int, err error) {
	switch p.getParamValByOp(op) {
	default:
		panic("unrecognized ValueSource")
	case ValueSourceNone:
		err = errors.New("none")
	case ValueSourceRegisterX:
		val = p.RegisterX
	case ValueSourceRegisterY:
		val = p.RegisterY
	case ValueSourceCellAge:
		val = p.Age
	case ValueSourceCellType:
		val = int(p.Type)
	case ValueSourceCellEnergy:
		val = int(math.Floor(p.Energy))
	case ValueSourceOrganismAge:
		val = p.OrganismAge
	case ValueSource0:
		val = 0
	case ValueSource1:
		val = 1
	case ValueSource2:
		val = 2
	case ValueSource3:
		val = 3
	case ValueSource4:
		val = 4
	case ValueSource5:
		val = 5
	case ValueSourceNeighborParticleType:
		dir := p.getParamValByOp(ValueSourceDirectionOp)
		x, y := getXYForDir(dir)
		val = int(w.Particles[w.Index(x, y)].Type)
	case ValueSourceNeighborIsSameCellType:
		dir := p.getParamValByOp(ValueSourceDirectionOp)
		x, y := getXYForDir(dir)
		if w.Particles[w.Index(x, y)].Type == p.Type {
			val = 1
		} else {
			val = 0
		}
	case ValueSourceNumber:
		val = p.getParamValByOp(ValueSourceNumberOp)
	case ValueSourceLabel:
		val = p.getParamValByOp(ValueSourceLabelOp)
	}
	return
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

func (p *Particle) StepDna(w *World) {
	if !p.Organic {
		return
	}
	pc := p.ExecutingDna.Index
	if pc == -1 {
		// DNA program is permanently halted
		return
	}
	instr := p.ExecutingDna.Instructions[pc]
	var opCode DnaOp = DnaOp(instr.OpCode) % DnaOpCount
	if opCode < ParameterOpCodeStart {
		switch opCode {
		default:
			panic(fmt.Sprintf("unhandled op code: %d", opCode))
		case NoOp:
			// done. that was easy.
		case CellDivisionOp:
			fmt.Println("cell at", p.Position, "attempts cell division")
		case CellDeathOp:
			fmt.Println("cell at", p.Position, "attempts cell death")
		case JumpOp:
			left, err := p.getValueSource(JumpOperandLeftOp, w)
			if err != nil {
				break
			}
			right, err := p.getValueSource(JumpOperandRightOp, w)
			if err != nil {
				break
			}
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
			newValue += pc
		}
		p.ParamValues[paramIndex] = newValue

	}
	pc += 1
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
