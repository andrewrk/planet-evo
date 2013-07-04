#ifndef CELL_H
#define CELL_H

#include "particle.h"
#include "dna.h"
#include "world.h"
#include "vec2.h"
#include <QVector>

class Cell : public Particle
{
public:
    Cell(ParticleType type, Vec2 pos, Dna intact_dna, Dna executing_dna);

    bool alive;
    double energy;
    int age;
    int organism_age;
    Dna intact_dna; // original DNA
    Dna executing_dna; // starts as a copy of intact_dna
    int param_values[PARAMETER_OP_CODE_COUNT];
    int register_x;
    int register_y;
    int waiting; // until this many steps are done, do nothing

    double getMutationChance() const;
    int getParamValByOp(DnaOp op) const;
    int getValueSource(DnaOp op, World *w);
    double getRadiansByOp(DnaOp op);
    bool performComparison(DnaOp op, int left, int right);
    int performCalc(DnaOp op, int left, int right);
    void saveToRegister(DnaOp op, int val);
    ParticleType getNewCellType(DnaOp op);
    void stepDna(World *w);

    void die();
    void split(Vec2 dir, ParticleType new_cell_type, int pc, double energy);

    void step(World *w);

    bool organic();

    void gainEnergy(int amt);
};

#endif // CELL_H
