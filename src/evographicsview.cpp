#include "evographicsview.h"
#include <QPainter>
#include <QtDebug>

EvoGraphicsView::EvoGraphicsView(QWidget *parent) :
    QWidget(parent)
{
    transform.translate(400, 200);
}

void EvoGraphicsView::setWorld(World *w)
{
    world = w;
}

Vec2 EvoGraphicsView::mapToWorld(QPoint pt)
{
    QTransform inverted = transform.inverted();
    QPointF new_pt = inverted.map(pt);
    return Vec2(new_pt.x(), new_pt.y());
}

void EvoGraphicsView::mousePressEvent(QMouseEvent *event)
{
    emit mousePress(event);
}

void EvoGraphicsView::mouseReleaseEvent(QMouseEvent *event)
{
    emit mouseRelease(event);
}

void EvoGraphicsView::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    painter.setBackground(Qt::white);
    painter.eraseRect(0, 0, this->width(), this->height());
    painter.setPen(QPen(Qt::black, 1));
    painter.setTransform(transform);
    foreach (Particle *p, world->particles) {
        painter.setBrush(p->color());
        painter.drawEllipse(QPointF(p->pos.x, p->pos.y), p->radius(), p->radius());
    }
}
