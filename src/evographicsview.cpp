#include "evographicsview.h"

EvoGraphicsView::EvoGraphicsView(QWidget *parent) :
    QGraphicsView(parent)
{
}

void EvoGraphicsView::mousePressEvent(QMouseEvent *event)
{
    emit mousePress(event);
}

void EvoGraphicsView::mouseReleaseEvent(QMouseEvent *event)
{
    emit mouseRelease(event);
}
