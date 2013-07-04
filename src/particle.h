#ifndef PARTICLE_H
#define PARTICLE_H

#include <QString>
#include <QColor>
#include "vec2.h"

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
    {"Null", 0, QColor("#000000"), 0, 0, 0, 1},
    {"Carbon", 1, QColor("#374B65"), 0, 0.99, 0, 1},
    {"Oxygen", 1, QColor("#94B4DD"), 0, 0.99, 0, 1},
    {"Dirt", 10, QColor("#6B3000"), 0, 0.1, 0.001, 1},
    {"Water", 10, QColor("#21009D"), 0, 0.7, 0.0001, 1},
    {"Light", 0, QColor("#FFF433"), 0, 0, 0, 1},

    // organic particles
    {"Chloro", 4, QColor("#0A7A00"), 5, 0.1, 0.001, 1},
    {"Fiber", 6,  QColor("#B75900"), 2, 0.5, 0.002, 1},
    {"Zygote", 5, QColor("#EFEFEF"), 10, 0.7, 0.0005, 1},

};

class Particle
{
public:
    Particle(ParticleType type, Vec2 pos);

    ParticleType type;
    Vec2 pos;
    Vec2 vel;

    virtual void step();
    virtual bool organic();

    QString name();
    double mass();
    double friction();
    QColor color();
    double maxEnergy();
    double radius();
};

#endif // PARTICLE_H
