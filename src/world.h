#ifndef WORLD_H
#define WORLD_H

#include <QGraphicsScene>

class World
{
public:
    World(QGraphicsScene *scene);

    void step();

private:
    int time;
    QGraphicsScene *scene;
};

#endif // WORLD_H
