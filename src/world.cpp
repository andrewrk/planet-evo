#include "world.h"
#include "cell.h"
#include "util.h"

#include <QtDebug>

World::World(QGraphicsScene *scene) :
    scene(scene)
{
    time = 0;
    radius = 480;
    scene->setSceneRect(-radius, -radius, radius * 2, radius * 2);
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
    p->view_object = sp.circle;
}

void World::destroyParticle(Particle *p)
{
    QGraphicsItem * circle = (QGraphicsItem*) p->view_object;
    scene->removeItem(circle);
}
