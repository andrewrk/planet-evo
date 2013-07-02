package main

import (
	"fmt"
	"math"
	"os"
)

type Instruction struct {
	OpCode byte
	Value  byte
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
	MutationChanceParam
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
	MutationChanceOp

	// meta
	DnaOpCount
)
const ParameterOpCodeStart = CellDivisionEnergyForNewCellOp
const ParameterOpCodeCount = DnaOpCount - ParameterOpCodeStart

type ParameterInfo struct {
	Type    ParameterType
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
	ValueSourceRandom

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
	{MutationChanceParam, 52},
}

func (p *Particle) getMutationChance() float64 {
	param := p.getParamValByOp(MutationChanceOp)
	const max = 0.05
	const min = 0.0001
	percent := float64(param) / 255
	return min + (percent * (max - min))
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

func (w *World) CreateSingleCelledPlant() Dna {
	dna := Dna{
		Index: 0,
		Instructions: []Instruction{
			{byte(ProgramEndBehaviorOp), byte(BorCContinue)},
			{byte(CellDivisionContingencyPlanOp), byte(BorCBlock)},
			{byte(CellDivisionEnergyForNewCellOp), byte(ValueSource4)},
			{byte(CellDivisionDirectionOp), byte(Direction8N)},
			{byte(CellDivisionNewCellTypeOp), 1},
			{byte(CellDivisionOp), 0},
			{byte(CellDivisionEnergyForNewCellOp), byte(ValueSource1)},
			{byte(CellDivisionOp), 0},
			{byte(CellDivisionOp), 0},
			{byte(CellDivisionOp), 0},
			{byte(CellDivisionOp), 0},
		},
	}
	return dna
}

func (w *World) CreateRandomDna() Dna {
	const instrCount = 10
	dna := Dna{
		Index:        0,
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

type CodeLabel struct {
	Pc   int
	Addr int
}

func (dna *Dna) Clone(p *Particle, w *World) Dna {
	newDna := *dna
	newDna.Instructions = make([]Instruction, 0, len(dna.Instructions))
	mutationChance := p.getMutationChance()
	// identify all code labels so we can preserve them
	labels := make([]CodeLabel, 0)
	for pc, instr := range dna.Instructions {
		opCode := DnaOp(instr.OpCode) % DnaOpCount
		paramIndex := opCode - ParameterOpCodeStart
		if paramIndex >= 0 {
			paramType := ParameterInfos[paramIndex].Type
			if paramType == CodeLabelParam {
				addr := int(instr.Value) + pc
				labels = append(labels, CodeLabel{pc, addr})
			}
		}
	}
	// copy instructions one at a time
	for _, instr := range dna.Instructions {
		// roll the dice
		roll := w.Rand.Float64()
		if roll <= mutationChance {
			switch w.Rand.Intn(3) {
			case 0: // 1/3 chance insert a byte here
				newDna.Instructions = append(newDna.Instructions, Instruction{
					byte(w.Rand.Intn(256)),
					byte(w.Rand.Intn(256)),
				})
				oldPc := len(newDna.Instructions)
				newDna.Instructions = append(newDna.Instructions, instr)
				// adjust labels - add 1 to every PC >= oldPc
				for i, label := range labels {
					if label.Pc >= oldPc {
						labels[i].Pc += 1
						labels[i].Addr += 1
					} else if label.Addr >= oldPc {
						labels[i].Addr += 1
					}
				}
			case 1: // 1/3 chance don't copy this byte
				oldPc := len(newDna.Instructions) + 1
				// adjust labels - subtract 1 from every PC >= oldPc
				for i, label := range labels {
					if label.Pc >= oldPc {
						labels[i].Pc -= 1
						labels[i].Addr -= 1
					} else if label.Addr >= oldPc {
						labels[i].Addr -= 1
					}
				}
			case 2: // 1/3 chance mangle the byte
				newDna.Instructions = append(newDna.Instructions, Instruction{
					byte(w.Rand.Intn(256)),
					byte(w.Rand.Intn(256)),
				})
			}
		} else {
			// copy the instruction correctly
			newDna.Instructions = append(newDna.Instructions, instr)
		}
	}
	// apply label adjustments
	for _, label := range labels {
		newDna.Instructions[label.Pc].Value = byte(label.Addr - label.Pc)
	}
	return newDna
}

func (dna *Dna) PerfectClone() Dna {
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
		dx, dy := getXYForDir(dir)
		x := int(math.Floor(p.Position.X)) + dx
		y := int(math.Floor(p.Position.Y)) + dy
		return int(w.Particles[w.Index(x, y)].Type)
	case ValueSourceNeighborIsSameCellType:
		dir := p.getParamValByOp(ValueSourceDirectionOp)
		dx, dy := getXYForDir(dir)
		x := int(math.Floor(p.Position.X)) + dx
		y := int(math.Floor(p.Position.Y)) + dy
		if w.Particles[w.Index(x, y)].Type == p.Type {
			return 1
		} else {
			return 0
		}
	case ValueSourceNumber:
		return p.getParamValByOp(ValueSourceNumberOp)
	case ValueSourceLabel:
		return p.getParamValByOp(ValueSourceLabelOp)
	case ValueSourceRandom:
		return w.Rand.Intn(256)
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
	fmt.Fprintf(os.Stderr, "Cell at %d, %d split into %s\n", x, y, ParticleClasses[newCellType].Name)
	baby := Particle{
		Type:         newCellType,
		Position:     iv(x, y),
		Organic:      true,
		Energy:       energy,
		IntactDna:    p.IntactDna.Clone(p, w),
		ExecutingDna: p.ExecutingDna.Clone(p, w),
		OrganismAge:  p.OrganismAge,
	}
	baby.InitParamValues()
	baby.ExecutingDna.Index = pc
	w.InsertParticle(baby, dx, dy)
}

func (w *World) InsertParticle(p Particle, dx int, dy int) {
	index := w.AltIndexVec2f(p.Position)
	destPart := w.Particles[index]
	if destPart.Type == NullParticle {
		w.Particles[index] = p
		return
	}
	w.Particles[index] = p
	destPart.Position.X += float64(dx)
	destPart.Position.Y += float64(dy)
	w.InsertParticle(destPart, dx, dy)
}

func (p *Particle) saveToRegister(op DnaOp, val int) {
	switch p.getParamValByOp(op) {
	default:
		panic("unrecognized register")
	case RegisterNone:
		return
	case RegisterX:
		p.RegisterX = val
	case RegisterY:
		p.RegisterY = val
	}
}

func (p *Particle) performCalc(op DnaOp, left int, right int) int {
	switch p.getParamValByOp(op) {
	case OperationLeft:
		return left
	case Operation0:
		return 0
	case Operation1:
		return 1
	case Operation2:
		return 2
	case OperationLeftPlus1:
		return left + 1
	case OperationLeftPlus2:
		return left + 2
	case OperationRight:
		return right
	case OperationRightPlus1:
		return right + 1
	case OperationRightPlus2:
		return right + 2
	case OperationLeftPlusRight:
		return left + right
	case OperationLeftMinusRight:
		return left - right
	case OperationRightMinusLeft:
		return right - left
	case OperationNegLeftNegRight:
		return -left - right
	case OperationLeftDivRight:
		if right == 0 {
			return 0
		}
		return left / right
	case OperationRightDivLeft:
		if left == 0 {
			return 0
		}
		return right / left
	case OperationLeftMult2:
		return left * 2
	case OperationRightMult2:
		return right * 2
	case OperationLeftModRight:
		if right == 0 {
			return 0
		}
		return left % right
	case OperationRightModLeft:
		if left == 0 {
			return 0
		}
		return right % left
	case OperationMin:
		if left < right {
			return left
		}
		return right
	case OperationMax:
		if left > right {
			return left
		}
		return right
	case OperationLeftMultRight:
		return left * right
	}
	panic("unrecognized calculation")
}

func (p *Particle) StepDna(w *World) {
	pc := p.ExecutingDna.Index
	if pc == -1 {
		// DNA program is permanently halted
		return
	}
	if p.Waiting > 0 {
		p.Waiting -= 1
		return
	}
	pc = pc % len(p.ExecutingDna.Instructions)
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
			newCellEnergy := float64(p.getValueSource(CellDivisionEnergyForNewCellOp, w))
			energyRequired := newCellEnergy + 1
			if p.Energy >= energyRequired {
				dir := p.getParamValByOp(CellDivisionDirectionOp)
				x, y := getXYForDir(dir)
				newCellType := p.getNewCellType(CellDivisionNewCellTypeOp)
				if newCellType == NullParticle {
					break
				}
				p.Energy -= newCellEnergy
				doWeFork := p.getParamValByOp(CellDivisionDoWeForkOp)
				newCellPc := pc
				if doWeFork == 1 {
					newCellPc = p.getParamValByOp(CellDivisionForkLabelOp)
				}
				p.Split(w, x, y, newCellType, newCellPc, newCellEnergy)
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
			fmt.Fprintf(os.Stderr, "Cell at %v triggered death.\n", p.Position)
			p.Die()
		case JumpOp:
			left := p.getValueSource(JumpOperandLeftOp, w)
			right := p.getValueSource(JumpOperandRightOp, w)
			comp := p.performComparison(JumpComparisonOp, left, right)
			newAddr := p.getParamValByOp(JumpLabelOp)
			if comp {
				pc = newAddr
			}
		case WaitOp:
			p.Waiting = p.getValueSource(WaitSourceOp, w) % 16
		case UpdateRegisterOp:
			val := p.getValueSource(UpdateRegisterSourceOp, w)
			p.saveToRegister(UpdateRegisterDestOp, val)
		case CalcOp:
			left := p.getValueSource(JumpOperandLeftOp, w)
			right := p.getValueSource(JumpOperandRightOp, w)
			val := p.performCalc(CalcOperationOp, left, right)
			p.saveToRegister(CalcDestOp, val)
		case ModifyDnaOp:
			addr := p.getParamValByOp(ModifyDnaLabelOp)
			addr = addr % len(p.ExecutingDna.Instructions)
			val := p.getValueSource(ModifyDnaSourceOp, w)
			p.ExecutingDna.Instructions[addr].Value = byte(val)
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
		switch p.ParamValues[ProgramEndBehaviorOp-ParameterOpCodeStart] {
		default:
			panic("unknown BlockOrContinue value")
		case BorCBlock:
			pc = -1
		case BorCContinue:
			pc = 0
		}
	}
	p.ExecutingDna.Index = pc
}
