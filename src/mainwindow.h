#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QMouseEvent>
#include "world.h"

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
    
private slots:
    void on_actionExit_triggered();
    void stepWorld();

    void on_actionSlower_triggered();

    void on_actionFaster_triggered();

    void on_actionDouble_triggered();

    void on_actionHalf_triggered();

    void on_actionTogglePause_triggered();

    void on_actionRestart_triggered();

    void on_graphicsView_mousePress(QMouseEvent *event);

private:
    Ui::MainWindow *ui;
    World *world;

    int speed = 1;
    int old_speed = 0;

    void updateSpeed(int new_speed);

    void restart();
};

#endif // MAINWINDOW_H
