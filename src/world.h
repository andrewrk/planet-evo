#ifndef WORLD_H
#define WORLD_H

#include <QGraphicsScene>

#include "dna.h"

class World
{
public:
    World(QGraphicsScene *scene);

    void step();

    void spawnRandomCreature(Vec2 pt);

    Particle *getParticleAt(Vec2 pt);

private:
    int time;
    QGraphicsScene *scene;
};

#endif // WORLD_H
