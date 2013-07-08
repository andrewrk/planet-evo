#ifndef WORLD_H
#define WORLD_H

#include <QGraphicsScene>
#include <QGraphicsItem>
#include <QList>

#include "dna.h"
#include "vec2.h"


class World
{
public:
    Vec2 size;

    World(QGraphicsScene *scene);

    void step();
    void spawnRandomCreature(Vec2 pt);
    Particle *getParticleAt(Vec2 pt);

    void addParticle(Particle *p);
    void destroyParticle(Particle *p);
private:
    int time;
    QGraphicsScene *scene;

    struct SceneParticle {
        Particle *particle;
        QGraphicsItem *circle;
    };

    QList<SceneParticle> particles;
};

#endif // WORLD_H
