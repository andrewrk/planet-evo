#ifndef CIRCLE_H
#define CIRCLE_H

#include "shape.h"

class Circle : public Shape {
public:
        typedef glm::float_t radius_t;

public:
        inline Circle(Circle const &other);
        inline Circle(position_t const &position = position_t(0.),
                radius_t radius = 1.);

        inline virtual ~Circle() = default;

        radius_t radius = 1.;
};

#include "circle.inl"
#endif
