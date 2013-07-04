#ifndef DNA_H
#define DNA_H

#include <QPoint>
#include <QVector>

#include "particle.h"

class Cell;


struct Instruction {
    uchar op_code;
    uchar value;
};

enum DnaOp {
    // core instructions
    NoOp,
    CellDivisionOp,
    CellDeathOp,
    JumpOp,
    WaitOp,
    UpdateRegisterOp,
    CalcOp,
    ModifyDnaOp,

    // parameter setting
    CellDivisionEnergyForNewCellOp,
    CellDivisionDirectionOp,
    CellDivisionNewCellTypeOp,
    CellDivisionForkLabelOp,
    CellDivisionDoWeForkOp,
    CellDivisionContingencyPlanOp,
    JumpOperandLeftOp,
    JumpOperandRightOp,
    JumpComparisonOp,
    JumpLabelOp,
    WaitSourceOp,
    UpdateRegisterSourceOp,
    UpdateRegisterDestOp,
    CalcOperandLeftOp,
    CalcOperandRightOp,
    CalcOperationOp,
    CalcDestOp,
    ModifyDnaLabelOp,
    ModifyDnaSourceOp,
    ValueSourceDirectionOp,
    ValueSourceNumberOp,
    ValueSourceLabelOp,
    ProgramEndBehaviorOp,
    MutationChanceOp,

    // meta
    DnaOpCount

};

enum ParameterType {
    ValueSourceParam,
    RegisterParam,
    ComparisonParam,
    OperationParam,
    CodeLabelParam,
    DirectionParam,
    NewCellTypeParam,
    BlockOrContinueParam,
    NumberParam,
    BooleanParam,
    MutationChanceParam
};

struct ParameterInfo {
    ParameterType type;
    int default_value;
};

enum BorC {
    BorCBlock,
    BorCContinue
};

enum ValueSource {
    ValueSourceRegisterX,
    ValueSourceRegisterY,
    ValueSourceCellAge,
    ValueSourceCellType,
    ValueSourceCellEnergy,
    ValueSourceOrganismAge,
    ValueSource0,
    ValueSource1,
    ValueSource2,
    ValueSource3,
    ValueSource4,
    ValueSource5,
    ValueSourceNeighborParticleType,
    ValueSourceNeighborIsSameCellType,
    ValueSourceNumber,
    ValueSourceLabel,
    ValueSourceRandom,

    // meta
    ValueSourceParamCount
};

enum Register {
    RegisterNone,
    RegisterX,
    RegisterY,

    // meta
    RegisterParamCount

};

enum Comparison {
    ComparisonNever,
    ComparisonAlways,
    ComparisonLeftZero,
    ComparisonLeftNonZero,
    ComparisonRightZero,
    ComparisonRightNonZero,
    ComparisonLeftGTRight,
    ComparisonLeftLTRight,
    ComparisonLeftEQRight,
    ComparisonLeftNERight,
    ComparisonLeftGERight,
    ComparisonLeftLERight,

    // meta
    ComparisonParamCount
};

enum Operation {
    OperationLeft,
    Operation0,
    Operation1,
    Operation2,
    OperationLeftPlus1,
    OperationLeftPlus2,
    OperationRight,
    OperationRightPlus1,
    OperationRightPlus2,
    OperationLeftPlusRight,
    OperationLeftMinusRight,
    OperationRightMinusLeft,
    OperationNegLeftNegRight,
    OperationLeftDivRight,
    OperationRightDivLeft,
    OperationLeftMult2,
    OperationRightMult2,
    OperationLeftModRight,
    OperationRightModLeft,
    OperationMin,
    OperationMax,
    OperationLeftMultRight,

    // meta
    OperationParamCount
};

const int PARAM_CAPS[] = {
    ValueSourceParamCount,
    RegisterParamCount,
    ComparisonParamCount,
    OperationParamCount,
    65535,
    8,
    ORGANIC_PARTICLE_COUNT + 1,
    2,
    256,
    2,
    256,
};

const int PARAMETER_OP_CODE_START = CellDivisionEnergyForNewCellOp;
const int PARAMETER_OP_CODE_COUNT = DnaOpCount - PARAMETER_OP_CODE_START;

const ParameterInfo PARAMETER_INFOS[] = {
    {ValueSourceParam, 5}, //CellDivisionEnergyForNewCellOp
    {DirectionParam, 0},
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
    {DirectionParam, 0},
    {NumberParam, 0},
    {CodeLabelParam, 1},
    {BlockOrContinueParam, BorCContinue},
    {MutationChanceParam, 52},
};

struct CodeLabel {
    int pc;
    int addr;
};

class Dna {
public:
    QVector<Instruction> instructions;
    int index; // position in code to execute next

    Dna clone(const Cell *c);
    Dna perfectClone();

    static QPoint getXYForDir(int dir);
    static Dna createSingleCelledPlant();
    static Dna createRandomDna();
    static Instruction createRandomInstruction();
};

#endif // DNA_H
