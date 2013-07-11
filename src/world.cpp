#include "world.h"
#include "cell.h"
#include "util.h"

#include <QtDebug>

World::World()
{
    time = 0;
    radius = 480;
}

void World::step()
{
    // add some carbon particles floating around randomly
    if (time % 40 == 0) {
        Vec2 pos = Vec2::direction(randRange(0.0, 2*PI));
        pos.setLength(randRange(0.0, radius));
        Particle *p = new Particle(CarbonParticle, pos);
        p->vel = Vec2::direction(randf() * 2 * PI);
        p->vel.setLength(randRange(-0.2, 0.2));
        addParticle(p);
    }

    // add some photons shooting in from the outside
    if (time % 10 == 0) {
        Vec2 offset = Vec2::direction(randRange(0.0, PI * 2));
        offset.setLength(radius);
        Particle *p = new Particle(LightParticle, offset);
        p->vel = offset.setLength(-1);
        addParticle(p);
    }

    // step and update positions
    foreach (Particle *p, particles) {
        p->step(this);
    }

    time += 1;
}

void World::spawnRandomCreature(Vec2 pt)
{
    Dna dna = Dna::createSingleCelledPlant();
    Cell *c = new Cell(ZygoteParticle, pt, dna, dna.perfectClone());
    c->energy = c->maxEnergy();
    addParticle(c);
}

Particle *World::getParticleAt(Vec2 pt)
{
    foreach (Particle *p, particles) {
        if (p->pos.distanceTo(pt) < p->radius())
            return p;
    }
    return NULL;
}

void World::addParticle(Particle *p)
{
    particles.append(p);
}

void World::destroyParticle(Particle *target)
{
    particles.removeOne(target);
}
