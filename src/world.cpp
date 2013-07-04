#include "world.h"
#include "cell.h"

World::World(QGraphicsScene *scene) :
    scene(scene)
{
    time = 0;
}

void World::step()
{

    time += 1;
}

void World::spawnRandomCreature(Vec2 pt)
{
    Dna dna = Dna::createSingleCelledPlant();
    Cell c = Cell(ZygoteParticle, pt, dna, dna.perfectClone());
    // TODO:
}

Particle *World::getParticleAt(Vec2 pt)
{
    // TODO
    return NULL;
}
