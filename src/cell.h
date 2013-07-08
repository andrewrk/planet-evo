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

    bool alive = true;
    double energy = 1;
    int organism_age = 0;
    Dna intact_dna; // original DNA
    Dna executing_dna; // starts as a copy of intact_dna
    int param_values[PARAMETER_OP_CODE_COUNT];
    int register_x = 0;
    int register_y = 0;
    int waiting = 0; // until this many steps are done, do nothing

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
    void split(Vec2 dir, ParticleType new_cell_type, int pc, double energy, World *w);

    virtual void step(World *w) override;

    virtual bool organic() override;

    void gainEnergy(int amt);
};

#endif // CELL_H
