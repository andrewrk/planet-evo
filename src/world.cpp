#include "world.h"
#include "cell.h"

#include <QtDebug>

World::World(QGraphicsScene *scene) :
    scene(scene)
{
    time = 0;
    size.x = 640;
    size.y = 480;
}

void World::step()
{
    for (int i = 0; i < particles.size(); i++) {
        SceneParticle sp = particles.at(i);
        sp.particle->step(this);
        sp.circle->setPos(sp.particle->pos.x, sp.particle->pos.y);
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
    QGraphicsItem * item = scene->itemAt(pt.x, pt.y);
    if (item == NULL) return NULL;

    return (Particle*) item->data(0).value<void *>();
}

void World::addParticle(Particle *p)
{
    SceneParticle sp;
    sp.particle = p;
    QRect r(-p->radius(), -p->radius(), p->radius() * 2, p->radius() * 2);
    sp.circle = scene->addEllipse(r, QColor("#000000"), p->color());
    sp.circle->setPos(p->pos.x, p->pos.y);
    sp.circle->setData(0, QVariant::fromValue((void *)p));
    particles.append(sp);
}
