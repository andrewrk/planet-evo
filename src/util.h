#ifndef UTIL_H
#define UTIL_H

#include <QtGlobal>

const double PI = 3.14159265358979;

inline double randf() {
    return qrand() / (double) RAND_MAX;
}

inline double randRange(double min, double max) {
    return min + randf() * (max - min);
}

inline int randRange(int min, int max) {
    return min + qrand() % (max - min);
}

#endif // UTIL_H
