#ifndef EVOGRAPHICSVIEW_H
#define EVOGRAPHICSVIEW_H

#include <QGraphicsView>

class EvoGraphicsView : public QGraphicsView
{
    Q_OBJECT
public:
    explicit EvoGraphicsView(QWidget *parent = 0);
    
signals:
    void mousePress(QMouseEvent *event);
    void mouseRelease(QMouseEvent *event);

protected:
    virtual void mousePressEvent(QMouseEvent *event);
    virtual void mouseReleaseEvent(QMouseEvent *event);
};

#endif // EVOGRAPHICSVIEW_H
