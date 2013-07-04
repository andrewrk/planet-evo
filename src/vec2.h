#ifndef VEC2_H
#define VEC2_H





class Vec2
{
public:
    double x;
    double y;

    Vec2(double x, double y) : x(x), y(y) {}
    Vec2() : x(0), y(0) {}

    double lengthSqrd() {
        return x * x + y * y;
    }

    double length() {
        return sqrt(lengthSqrd());
    }

    Vec2& scale(double scalar) {
        x *= scalar;
        y *= scalar;
        return *this;
    }

    Vec2& operator+=(const Vec2 &other) {
        x += other.x;
        y += other.y;
        return *this;
    }
    Vec2& operator-=(const Vec2 &other) {
        x -= other.x;
        y -= other.y;
        return *this;
    }
    Vec2& operator*=(double scalar) {
        x *= scalar;
        y *= scalar;
        return *this;
    }
    Vec2& operator/=(double scalar) {
        x /= scalar;
        y /= scalar;
        return *this;
    }

    Vec2& normalize() {
        double len = length();
        x /= len;
        y /= len;
        return *this;
    }

    // set x and y to zero and return itself
    Vec2& clear() {
        x = 0;
        y = 0;
        return *this;
    }

    // reduce the length by a certain amount
    Vec2& retract(double amount) {
        return setLength(length() - amount);
    }

    // increase the length by a certain amount
    Vec2& extend(double amount) {
        return setLength(length() + amount);
    }

    Vec2& setLength(double length) {
        return normalize().scale(length);
    }

    double angle() {
        return atan2(y, x);
    }

    static Vec2 direction(double radians) {
        return Vec2(cos(radians), sin(radians));
    }
};


inline Vec2 operator+(Vec2 &left, const Vec2& right) {
    left += right;
    return left;
}


inline Vec2 operator-(Vec2 left, const Vec2& right) {
    left -= right;
    return left;
}

inline Vec2 operator*(Vec2 left, double scalar) {
    left *= scalar;
    return left;
}

inline Vec2 operator/(Vec2 left, double scalar) {
    left /= scalar;
    return left;
}

#endif // VEC2_H
