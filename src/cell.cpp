#include "cell.h"
#include "util.h"

Cell::Cell(ParticleType type, Vec2 pos, Dna intact_dna, Dna executing_dna) :
    Particle(type, pos),
    intact_dna(intact_dna),
    executing_dna(executing_dna)
{
    for (int i = 0; i < PARAMETER_OP_CODE_COUNT; i++) {
        param_values[i] = PARAMETER_INFOS[i].default_value;
    }
}

double Cell::getMutationChance() const
{
    int param = this->getParamValByOp(MutationChanceOp);
    const double max = 0.05;
    const double min = 0.0001;
    double percent = ((double)param) / 255.0;
    return min + (percent * (max - min));
}

int Cell::getParamValByOp(DnaOp op) const
{
    int param_index = op - PARAMETER_OP_CODE_START;
    int val = param_values[param_index];
    return val % PARAM_CAPS[PARAMETER_INFOS[param_index].type];
}

int Cell::getValueSource(DnaOp op, World *w)
{
    switch (getParamValByOp(op)) {
    case ValueSourceRegisterX:
        return register_x;
    case ValueSourceRegisterY:
        return register_y;
    case ValueSourceCellAge:
        return age;
    case ValueSourceCellType:
        return type;
    case ValueSourceCellEnergy:
        return floor(energy);
    case ValueSourceOrganismAge:
        return organism_age;
    case ValueSource0:
        return 0;
    case ValueSource1:
        return 1;
    case ValueSource2:
        return 2;
    case ValueSource3:
        return 3;
    case ValueSource4:
        return 4;
    case ValueSource5:
        return 5;
    case ValueSourceNeighborParticleType:
    {
        double rads = getRadiansByOp(ValueSourceDirectionOp);
        Vec2 v = Vec2::direction(rads);
        v.setLength(radius() * 1.5);
        Particle *p = w->getParticleAt(pos + v);
        return p == NULL ? 0 : p->type;
    }
    case ValueSourceNeighborIsSameCellType:
    {
        double rads = getRadiansByOp(ValueSourceDirectionOp);
        Vec2 v = Vec2::direction(rads);
        v.setLength(radius() * 1.5);
        Particle *p = w->getParticleAt(pos + v);
        ParticleType t = p == NULL ? NullParticle : p->type;
        return t == type;
    }
    case ValueSourceNumber:
        return getParamValByOp(ValueSourceNumberOp);
    case ValueSourceLabel:
        return getParamValByOp(ValueSourceLabelOp);
    case ValueSourceRandom:
        return rand() % 256;
    default:
        qFatal("unknown value source");
        throw;
    }
}

double Cell::getRadiansByOp(DnaOp op)
{
    double val = getParamValByOp(op);
    double percent = val / 255.0;
    return percent * 2 * PI;
}

bool Cell::performComparison(DnaOp op, int left, int right)
{
    switch (getParamValByOp(op)) {
    case ComparisonNever:
        return false;
    case ComparisonAlways:
        return true;
    case ComparisonLeftZero:
        return left == 0;
    case ComparisonLeftNonZero:
        return left != 0;
    case ComparisonRightZero:
        return right == 0;
    case ComparisonRightNonZero:
        return right != 0;
    case ComparisonLeftGTRight:
        return left > right;
    case ComparisonLeftLTRight:
        return left < right;
    case ComparisonLeftEQRight:
        return left == right;
    case ComparisonLeftNERight:
        return left != right;
    case ComparisonLeftGERight:
        return left >= right;
    case ComparisonLeftLERight:
        return left <= right;
    default:
        qFatal("unrecognized Comparison");
        throw;
    }
}

int Cell::performCalc(DnaOp op, int left, int right)
{
    switch (getParamValByOp(op)) {
    case OperationLeft:
        return left;
    case Operation0:
        return 0;
    case Operation1:
        return 1;
    case Operation2:
        return 2;
    case OperationLeftPlus1:
        return left + 1;
    case OperationLeftPlus2:
        return left + 2;
    case OperationRight:
        return right;
    case OperationRightPlus1:
        return right + 1;
    case OperationRightPlus2:
        return right + 2;
    case OperationLeftPlusRight:
        return left + right;
    case OperationLeftMinusRight:
        return left - right;
    case OperationRightMinusLeft:
        return right - left;
    case OperationNegLeftNegRight:
        return -left - right;
    case OperationLeftDivRight:
        return right == 0 ? 0 : left / right;
    case OperationRightDivLeft:
        return left == 0 ? 0 : right / left;
    case OperationLeftMult2:
        return left * 2;
    case OperationRightMult2:
        return right * 2;
    case OperationLeftModRight:
        return right == 0 ? 0 : left % right;
    case OperationRightModLeft:
        return left == 0 ? 0 : right % left;
    case OperationMin:
        return left < right ? left : right;
    case OperationMax:
        return left > right ? left : right;
    case OperationLeftMultRight:
        return left * right;
    default:
        qFatal("unrecognized calculation");
        throw;
    }
}

void Cell::saveToRegister(DnaOp op, int val)
{
    switch (getParamValByOp(op)) {
    case RegisterNone:
        return;
    case RegisterX:
        register_x = val;
        return;
    case RegisterY:
        register_y = val;
        return;
    default:
        qFatal("unrecognized register");
        throw;
    }

}

ParticleType Cell::getNewCellType(DnaOp op)
{
    int new_cell_enum_value = getParamValByOp(op);
    return new_cell_enum_value == 0 ? NullParticle : (ParticleType)(new_cell_enum_value - 1 + FIRST_ORGANIC_PARTICLE);
}

void Cell::stepDna(World *w)
{
    int pc = executing_dna.index;
    if (pc == -1) {
        // DNA program is permanently halted
        return;
    }
    if (waiting > 0) {
        waiting -= 1;
        return;
    }
    pc = pc % executing_dna.instructions.size();
    Instruction instr = executing_dna.instructions.at(pc);
    pc += 1;
    DnaOp op_code = (DnaOp)(instr.op_code % DnaOpCount);
    if (op_code < PARAMETER_OP_CODE_START) {
        switch (op_code) {
        default:
            qFatal("unhandled op code");
            throw;
        case NoOp:
            // done. that was easy.
            break;
        case CellDivisionOp:
        {
            int new_cell_energy = getValueSource(CellDivisionEnergyForNewCellOp, w);
            int energy_required = new_cell_energy + 1;
            if (energy > energy_required) {
                double dir = getRadiansByOp(CellDivisionDirectionOp);
                ParticleType new_cell_type = getNewCellType(CellDivisionNewCellTypeOp);
                if (new_cell_type == NullParticle) break;
                energy -= new_cell_energy;
                int do_we_fork = getParamValByOp(CellDivisionDoWeForkOp);
                int new_cell_pc = do_we_fork == 1 ? getParamValByOp(CellDivisionForkLabelOp) : pc;
                split(Vec2::direction(dir), new_cell_type, new_cell_pc, new_cell_energy, w);
            } else {
                int plan = getParamValByOp(CellDivisionContingencyPlanOp);
                switch (plan) {
                default:
                    qFatal("urecognized contingency plan");
                    throw;
                case BorCBlock:
                    pc -= 1;
                    break;
                case BorCContinue:
                    break;
                }
            }
        }
            break;
        case CellDeathOp:
            die();
            break;
        case JumpOp:
        {
            int left = getValueSource(JumpOperandLeftOp, w);
            int right = getValueSource(JumpOperandRightOp, w);
            bool comp = performComparison(JumpComparisonOp, left, right);
            int new_addr = getParamValByOp(JumpLabelOp);
            if (comp) pc = new_addr;
        }
            break;
        case WaitOp:
            waiting = getValueSource(WaitSourceOp, w) % 16;
            break;
        case UpdateRegisterOp:
        {
            int val = getValueSource(UpdateRegisterSourceOp, w);
            saveToRegister(UpdateRegisterDestOp, val);
        }
            break;
        case CalcOp:
        {
            int left = getValueSource(JumpOperandLeftOp, w);
            int right = getValueSource(JumpOperandRightOp, w);
            int val = performCalc(CalcOperationOp, left, right);
            saveToRegister(CalcDestOp, val);
        }
            break;
        case ModifyDnaOp:
        {
            int addr = getParamValByOp(ModifyDnaLabelOp);
            addr = addr % executing_dna.instructions.size();
            int val = getValueSource(ModifyDnaSourceOp, w);
            executing_dna.instructions[addr].value = val;
        }
            break;
        }
    } else {
        // set a parameter value
        int param_index = op_code - PARAMETER_OP_CODE_START;
        int param_type = PARAMETER_INFOS[param_index].type;
        int value_cap = PARAM_CAPS[param_type];
        int new_value = instr.value % value_cap;
        if (param_type == CodeLabelParam) {
            new_value += pc - 1;
        }
        param_values[param_index] = new_value;
    }
    if (pc >= executing_dna.instructions.size()) {
        switch (param_values[ProgramEndBehaviorOp - PARAMETER_OP_CODE_START]) {
        case BorCBlock:
            pc = -1;
        case BorCContinue:
            pc = 0;
        default:
            qFatal("unknown blockorcontinue value");
            throw;
        }
    }
    executing_dna.index = pc;
}

void Cell::die()
{
    energy = 0;
    alive = false;
}

void Cell::split(Vec2 dir, ParticleType new_cell_type, int pc, double energy, World *w)
{
    // how is babby formed
    double new_radius = PARTICLE_CLASSES[new_cell_type].radius;
    Vec2 new_pos = pos + dir * (radius() + new_radius);
    Cell *baby = new Cell(new_cell_type, new_pos, intact_dna.clone(this), executing_dna.clone(this));
    baby->energy = energy;
    baby->organism_age = organism_age;
    baby->executing_dna.index = pc;

    w->addParticle(baby);
}

void Cell::step(World *w)
{
    if (alive) {
        stepDna(w);
        age += 1;
        organism_age += 1;
        energy -= 0.001;
        if (energy <= 0) die();
    }
    Particle::step();
}

bool Cell::organic()
{
    return true;
}

void Cell::gainEnergy(int amt)
{
    energy += amt;
    int max = maxEnergy();
    if (energy > max) energy = max;
}
