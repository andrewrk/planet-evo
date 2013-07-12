#ifndef SHAPE_H
#define SHAPE_H

#define GLM_PRECISION_HIGHP_FLOAT
#include <glm/glm.hpp>

class Shape {
public:
        typedef glm::vec2 position_t;

public:
        inline Shape(position_t const &position = position_t(0.));
        inline Shape(Shape const &other);

        inline virtual ~Shape() = 0;

        position_t position = position_t(0.);
};

#include "shape.inl"
#endif
