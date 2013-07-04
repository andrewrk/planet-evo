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
    src/particle.cpp

HEADERS  += \
    src/world.h \
    src/mainwindow.h \
    src/particle.h

FORMS    += \
    src/mainwindow.ui
