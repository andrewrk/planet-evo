#ifndef CIRCLE_HPP
#define CIRCLE_HPP

#include "shape.h"

class Circle : public Shape {
public:
        typedef glm::float_t radius_t;

public:
        Circle(Circle const &other);
        Circle(position_t const &position = position_t(0.),
                radius_t radius = 1.);

        virtual ~Circle() = default;

        radius_t radius = 1.;
};

#endif

#include "circle.inl"
