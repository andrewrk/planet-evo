#include "world.h"
#include "cell.h"
#include "util.h"
#include "collision-detection.h"

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

    for (Particle *k : particles) {
        for (Particle *l : particles) {
                if (k == l) {
                        continue;
                }

                Circle a(Circle::position_t(k->pos.x, k->pos.y), k->radius());
                Circle b(Circle::position_t(l->pos.x, l->pos.y), l->radius());

                if (isIntersecting(a, b)) {
                        glm::vec2 normal =
                                glm::normalize(b.position - a.position);

                        Vec2 m = l->vel - k->vel;
                        glm::vec2 rv(m.x, m.y);
                        double speed = glm::dot(rv, normal);

                        if (speed > 0) {
                                continue;
                        }

                        double e = glm::min(k->elasticity(), l->elasticity());
                        double j = -(1 + e) * speed;

                        j /= 1 / k->mass() + 1 / l->mass();

                        glm::vec2 impulse = normal * j;

                        glm::vec2 n = impulse * (1 / k->mass());
                        k->vel -= Vec2(n.x, n.y);

                        n = impulse * (1 / l->mass());
                        l->vel -= Vec2(n.x, n.y);
                }
        }
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
