#include "world.h"
#include "cell.h"
#include "util.h"

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
    // add some carbon particles
    if (time % 40 == 0) {
        Vec2 pos(randRange(0.0, size.x), randRange(0.0, size.y));
        Particle *p = new Particle(CarbonParticle, pos);
        p->vel = Vec2::direction(randf() * 2 * PI);
        p->vel.setLength(randRange(-0.2, 0.2));
        addParticle(p);
    }

    // add some photons
    Vec2 center = size / 2;
    if (time % 10 == 0) {
        Vec2 offset = Vec2::direction(randRange(0.0, PI * 2));
        offset.setLength(size.length() * 0.5);
        Particle *p = new Particle(LightParticle, center + offset);
        p->vel = offset.setLength(-1);
        addParticle(p);
    }

    // step and update positions
    for (int i = 0; i < particles.size(); i++) {
        SceneParticle sp = particles.at(i);
        sp.particle->step(this);
        sp.circle->setPos(sp.particle->pos.x, sp.particle->pos.y);
    }

    // resolve collisions
//    for (int i = 0; i < particles.size(); i++) {
//        SceneParticle sp = particles.at(i);
//        Particle *p = sp.particle;
//        auto collisions = scene->collidingItems(sp.circle);
//        if (collisions.size() > 0) {
//            QGraphicsItem *item = collisions.first();
//            Particle *otherPart = (Particle *)item->data(0).value<void *>();
//            // calculate normal
//            Vec2 normal = otherPart->pos - p->pos;
//            normal.normalize();
//            // calculate relative velocity
//            Vec2 rv = otherPart->vel - p->vel;
//            // calculate relative velocity in terms of the normal direction
//            double vel_along_normal = rv.dot(normal);
//            // do not resolve if velocities are separating
//            if (vel_along_normal > 0) continue;
//            // calculate restitution
//            double e = qMin(p->elasticity(), otherPart->elasticity());
//            // calculate impulse scalar
//            double j = -(1 + e) * vel_along_normal;
//            j /= 1 / p->mass() + 1 / otherPart->mass();
//            // apply impulse
//            Vec2 impulse = normal * j;
//            p->vel -= impulse * (1 / p->mass());
//            otherPart->vel -= impulse * (1 / otherPart->mass());
//        }
//    }

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
