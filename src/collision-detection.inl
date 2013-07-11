#ifndef COLLISION_DETECTION_INL
#define COLLISION_DETECTION_INL

#include <glm/gtx/norm.hpp>

bool isIntersecting(Circle const &first, Circle const &second)
{
        glm::float_t distance_squared =
                glm::length2(first.position - second.position);

        glm::float_t radii = first.radius + second.radius;

        if (distance_squared <= radii * radii) {
                return true;
        }

        return false;
}

#endif
