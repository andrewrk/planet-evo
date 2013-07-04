#include "dna.h"
#include "cell.h"

#include "util.h"

Dna Dna::clone(const Cell *c)
{
    Dna new_dna = *this;
    double mutation_chance = c->getMutationChance();
    // identify all the code labels so we can preserve them
    QVector<CodeLabel> labels;
    for (int pc = 0; pc < this->instructions.size(); pc++) {
        Instruction instr = this->instructions.at(pc);
        int op_code = instr.op_code % DnaOpCount;
        int param_index = op_code - PARAMETER_OP_CODE_START;
        if (param_index >= 0) {
            ParameterType param_type = PARAMETER_INFOS[param_index].type;
            if (param_type == CodeLabelParam) {
                int addr = instr.value + pc;
                labels.append(CodeLabel{pc, addr});
            }
        }
    }
    // copy instructions one at a time
    foreach (Instruction instr, this->instructions) {
        // roll the dice
        double roll = randf();
        if (roll <= mutation_chance) {
            switch (qrand() % 3) {
            case 0: // 1/3 chance insert a byte here
            {
                new_dna.instructions.append(Instruction{
                    (uchar)(qrand() % 256),
                    (uchar)(qrand() % 256)
                });
                int old_pc = new_dna.instructions.size();
                new_dna.instructions.append(instr);
                // adjust labels - add 1 to every PC >= old_pc
                for (int i = 0; i < labels.size(); i++) {
                    CodeLabel label = labels.at(i);
                    if (label.pc >= old_pc) {
                        label.pc += 1;
                        label.addr += 1;
                        labels.replace(i, label);
                    } else if (label.addr >= old_pc) {
                        label.addr += 1;
                        labels.replace(i, label);
                    }
                }
            }
                break;
            case 1: // 1/3 chance don't copy this byte
            {
                int old_pc = new_dna.instructions.size() + 1;
                // adjust labels - subtract 1 from every pc >= old_pc
                for (int i = 0; i < labels.size(); i++) {
                    CodeLabel label = labels.at(i);
                    if (label.pc >= old_pc) {
                        label.pc -= 1;
                        label.addr -= 1;
                        labels.replace(i, label);
                    } else if (label.addr >= old_pc) {
                        label.addr -= 1;
                        labels.replace(i, label);
                    }
                }
            }
                break;
            case 2: // 1/3 chance mangle the byte
                new_dna.instructions.append(Instruction{
                    (uchar)(qrand() % 256),
                    (uchar)(qrand() % 256)
                });
                break;
            default:
                qFatal("unexpected random number");
                throw;
            }
        } else {
            // copy instruction correctly
            new_dna.instructions.append(instr);
        }
    }
    // apply label adjustments
    foreach (CodeLabel label, labels) {
        Instruction instr = new_dna.instructions.at(label.pc);
        instr.value = label.addr - label.pc;
        new_dna.instructions.replace(label.pc, instr);
    }
    return new_dna;
}

Dna Dna::perfectClone()
{
    Dna new_dna = *this;
    return new_dna;
}

Dna Dna::createSingleCelledPlant()
{
    Dna dna;
    dna.index = 0;
    dna.instructions.append(Instruction{ProgramEndBehaviorOp, BorCContinue});
    dna.instructions.append(Instruction{CellDivisionContingencyPlanOp, BorCBlock});
    dna.instructions.append(Instruction{CellDivisionEnergyForNewCellOp, ValueSource4});
    dna.instructions.append(Instruction{CellDivisionDirectionOp, 0});
    dna.instructions.append(Instruction{CellDivisionNewCellTypeOp, 1});
    dna.instructions.append(Instruction{CellDivisionOp, 0});
    dna.instructions.append(Instruction{CellDivisionEnergyForNewCellOp, ValueSource1});
    dna.instructions.append(Instruction{CellDivisionOp, 0});
    dna.instructions.append(Instruction{CellDivisionOp, 0});
    dna.instructions.append(Instruction{CellDivisionOp, 0});
    dna.instructions.append(Instruction{CellDivisionOp, 0});
    return dna;
}

Dna Dna::createRandomDna()
{
    const int instr_count = 10;
    Dna dna;
    dna.index = 0;
    for (int i = 0; i < instr_count; i++) {
        dna.instructions.append(createRandomInstruction());
    }
    return dna;
}

Instruction Dna::createRandomInstruction()
{
    Instruction instr;
    instr.op_code = qrand() % 256;
    instr.value = qrand() % 256;
    return instr;
}

