#ifndef PARTICLE_H
#define PARTICLE_H

#include <QString>
#include <QColor>
#include "vec2.h"

class World;

enum ParticleType {
    // non-organic particles
    NullParticle,
    CarbonParticle,
    OxygenParticle,
    DirtParticle,
    WaterParticle,
    LightParticle,

    // organic particles
    ChloroParticle,
    FiberParticle,
    ZygoteParticle,

    // meta
    ParticleCount
};

const int FIRST_ORGANIC_PARTICLE = ChloroParticle;
const int ORGANIC_PARTICLE_COUNT = ParticleCount - FIRST_ORGANIC_PARTICLE;

struct ParticleClass {
    QString name;
    double mass;
    QColor color;
    double max_energy;
    double elasticity;
    double friction;
    double radius;
};

const ParticleClass PARTICLE_CLASSES[] = {
    {"Null", 0, QColor("#000000"), 0, 0, 0, 4},
    {"Carbon", 1, QColor("#374B65"), 0, 0.99, 0, 4},
    {"Oxygen", 1, QColor("#94B4DD"), 0, 0.99, 0, 4},
    {"Dirt", 10, QColor("#6B3000"), 0, 0.1, 0.001, 4},
    {"Water", 10, QColor("#21009D"), 0, 0.7, 0.0001, 4},
    {"Light", 0, QColor("#FFF433"), 0, 0, 0, 4},

    // organic particles
    {"Chloro", 4, QColor("#0A7A00"), 5, 0.1, 0.001, 4},
    {"Fiber", 6,  QColor("#B75900"), 2, 0.5, 0.002, 4},
    {"Zygote", 5, QColor("#EFEFEF"), 10, 0.7, 0.0005, 4},

};

class Particle
{
public:
    Particle(ParticleType type, Vec2 pos);

    ParticleType type;
    Vec2 pos;
    Vec2 vel;

    virtual void step(World *w);
    virtual bool organic();

    QString name();
    double mass();
    double friction();
    QColor color();
    double maxEnergy();
    double radius();
};

#endif // PARTICLE_H
