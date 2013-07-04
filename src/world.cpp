#include "world.h"

World::World(QGraphicsScene *scene) :
    scene(scene)
{
    time = 0;
}

void World::step()
{

    time += 1;
}
