#include "world.h"
#include "cell.h"

World::World(QGraphicsScene *scene) :
    scene(scene)
{
    time = 0;
    size.x = 640;
    size.y = 480;
    //scene->setSceneRect(0, 0, size.x, size.y);
}

void World::step()
{
    for (int i = 0; i < particles.size(); i++) {
        SceneParticle sp = particles.at(i);
        sp.particle->step();
        sp.circle->setPos(sp.particle->pos.x, sp.particle->pos.y);
    }
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

void World::addParticle(Particle *p)
{
    SceneParticle sp;
    sp.particle = p;
    QRect r(p->pos.x - p->radius() , p->pos.y - p->radius() , p->radius() * 2, p->radius() * 2);
    sp.circle = scene->addEllipse(r, QColor("#000000"), p->color());
    particles.append(sp);
}
