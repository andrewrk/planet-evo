#ifndef SHAPE_HPP
#define SHAPE_HPP

#include <glm/glm.hpp>

class Shape {
public:
        typedef glm::vec2 position_t;

public:
        inline Shape(position_t const &position = position_t(0.));
        inline Shape(Shape const &other);

        virtual ~Shape() = 0;

        position_t position = position_t(0.);
};

#endif

#include "shape.inl"
