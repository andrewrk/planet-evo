#ifndef CIRCLE_INL
#define CIRCLE_INL

Circle::Circle(Circle const &other)
        : Shape(other)
        , radius { other.radius }
{
}

Circle::Circle(position_t const &p, radius_t r)
        : Shape(p)
        , radius { r }
{
}

#endif
