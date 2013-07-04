#include "particle.h"

Particle::Particle(ParticleType type, Vec2 pos) :
    type(type),
    pos(pos)
{
}

void Particle::step()
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
}

bool Particle::organic()
{
    return false;
}

QString Particle::name()
{
    return PARTICLE_CLASSES[type].name;
}

double Particle::mass()
{
    return PARTICLE_CLASSES[type].mass;
}

double Particle::friction()
{
    return PARTICLE_CLASSES[type].friction;
}

QColor Particle::color()
{
    return PARTICLE_CLASSES[type].color;
}

double Particle::maxEnergy()
{
    return PARTICLE_CLASSES[type].max_energy;
}

double Particle::radius()
{
    return PARTICLE_CLASSES[type].radius;
}

