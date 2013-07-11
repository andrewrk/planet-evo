#ifndef EVOGRAPHICSVIEW_H
#define EVOGRAPHICSVIEW_H

#include <QWidget>
#include <QTransform>
#include "world.h"
#include "vec2.h"

class EvoGraphicsView : public QWidget
{
    Q_OBJECT
public:
    explicit EvoGraphicsView(QWidget *parent = 0);

    void setWorld(World *w);
    Vec2 mapToWorld(QPoint pt);

    
signals:
    void mousePress(QMouseEvent *event);
    void mouseRelease(QMouseEvent *event);

protected:
    virtual void mousePressEvent(QMouseEvent *event);
    virtual void mouseReleaseEvent(QMouseEvent *event);
    virtual void paintEvent(QPaintEvent *) override;

private:
    World *world;
    QTransform transform;
};

#endif // EVOGRAPHICSVIEW_H
