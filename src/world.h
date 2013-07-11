#ifndef WORLD_H
#define WORLD_H

#include <QList>

#include "dna.h"
#include "vec2.h"


class World
{
public:
    double radius;
    QList<Particle *> particles;

    World();

    void step();
    void spawnRandomCreature(Vec2 pt);
    Particle *getParticleAt(Vec2 pt);

    void addParticle(Particle *p);
    void destroyParticle(Particle *target);
private:
    int time;
};

#endif // WORLD_H
