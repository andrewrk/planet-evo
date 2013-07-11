#ifndef SHAPE_INL
#define SHAPE_INL

Shape::Shape(position_t const &p)
        : position { p }
{
}

Shape::Shape(Shape const &other)
        : position { other.position }
{
}

Shape::~Shape()
{
}

#endif
