#include "particle.h"
#include "world.h"

Particle::Particle(ParticleType type, Vec2 pos) :
    type(type),
    pos(pos),
    vel(0, 0)
{
}

void Particle::step(World *w)
{
    // apply velocity
    pos += vel;

    // apply friction
    double frict = friction();
    if (frict >= vel.length()) {
        vel.clear();
    } else {
        vel.retract(frict);
    }

    age += 1;
    int max_age = maxAge();
    if (max_age >= 0 && age > max_age) {
        w->destroyParticle(this);
    }
}

bool Particle::organic()
{
    return false;
}

QString Particle::name() const
{
    return PARTICLE_CLASSES[type].name;
}

double Particle::mass() const
{
    return PARTICLE_CLASSES[type].mass;
}

double Particle::friction() const
{
    return PARTICLE_CLASSES[type].friction;
}

QColor Particle::color() const
{
    return PARTICLE_CLASSES[type].color;
}

double Particle::maxEnergy() const
{
    return PARTICLE_CLASSES[type].max_energy;
}

double Particle::elasticity() const
{
    return PARTICLE_CLASSES[type].elasticity;
}

double Particle::radius() const
{
    return PARTICLE_CLASSES[type].radius;
}

int Particle::maxAge() const
{
    return PARTICLE_CLASSES[type].max_age;
}

