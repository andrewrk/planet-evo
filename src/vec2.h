#ifndef VEC2_H
#define VEC2_H

#include <cmath>

class Vec2
{
public:
    double x;
    double y;

    Vec2(double x, double y) : x(x), y(y) {}
    Vec2() : x(0), y(0) {}

    double lengthSqrd() const {
        return x * x + y * y;
    }

    double length() const {
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
    Vec2& operator*=(const Vec2 &other) {
        x *= other.x;
        y *= other.y;
        return *this;
    }
    Vec2& operator/=(double scalar) {
        x /= scalar;
        y /= scalar;
        return *this;
    }
    Vec2& operator/=(const Vec2 &other) {
        x /= other.x;
        y /= other.y;
        return *this;
    }

    Vec2& normalize() {
        double len = length();
        x /= len;
        y /= len;
        return *this;
    }

    Vec2 normalized() const {
        Vec2 v = *this;
        v.normalize();
        return v;
    }

    double dot(const Vec2 &other) const {
        return x * other.x + y * other.y;
    }

    // set x and y to zero and return itself
    Vec2& clear() {
        x = 0;
        y = 0;
        return *this;
    }

    Vec2 clone() const {
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

    double angle() const {
        return atan2(y, x);
    }

    double distanceTo(const Vec2 &other) {
        return sqrt(pow(other.x - x, 2) + pow(other.y - y, 2));
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

inline Vec2 operator*(Vec2 left, const Vec2& right) {
    left *= right;
    return left;
}

inline Vec2 operator/(Vec2 left, double scalar) {
    left /= scalar;
    return left;
}

inline Vec2 operator/(Vec2 left, const Vec2& right) {
    left /= right;
    return left;
}
#endif // VEC2_H
