#-------------------------------------------------
#
# Project created by QtCreator 2013-07-03T21:12:36
#
#-------------------------------------------------

QT       += core gui

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = planet-evo
TEMPLATE = app

QMAKE_CXXFLAGS += -std=c++11

SOURCES +=\
    src/world.cpp \
    src/main.cpp \
    src/mainwindow.cpp \
    src/particle.cpp \
    src/cell.cpp \
    src/dna.cpp \
    src/evographicsview.cpp

HEADERS  += \
    src/world.h \
    src/mainwindow.h \
    src/particle.h \
    src/cell.h \
    src/dna.h \
    src/vec2.h \
    src/util.h \
    src/evographicsview.h \
    src/shape.h \
    src/shape.inl \
    src/circle.h \
    src/circle.inl \
    src/collision-detection.h \
    src/collision-detection.inl

FORMS    += \
    src/mainwindow.ui
