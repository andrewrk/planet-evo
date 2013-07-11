#include "mainwindow.h"
#include "ui_mainwindow.h"

#include <QDebug>
#include <QTimer>
#include "vec2.h"
#include "util.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow),
    world(NULL)
{
    ui->setupUi(this);
    ui->actionExit->setShortcut(QKeySequence("Alt+F4"));

    restart();

    QTimer *timer = new QTimer(this);
    bool ok;
    ok = connect(timer, SIGNAL(timeout()), this, SLOT(stepWorld()));
    Q_ASSERT(ok);
    timer->start(16);
}

MainWindow::~MainWindow()
{
    delete ui;
    delete world;
}

void MainWindow::on_actionExit_triggered()
{
    this->close();
}

void MainWindow::stepWorld()
{
    for (int i = 0; i < speed; i++) {
        world->step();
    }
    ui->graphicsView->update();
}

void MainWindow::on_actionSlower_triggered()
{
    updateSpeed(speed - 1);
}

void MainWindow::updateSpeed(int new_speed)
{
    speed = new_speed;
    if (speed < 0) speed = 0;
    ui->statusBar->showMessage(QString("Speed: %1x").arg(speed), 3000);
}

void MainWindow::restart()
{
    uint seed = 1234;
    qsrand(seed);
    qDebug() << "Using seed" << seed;
    delete world;
    world = new World();
    ui->graphicsView->setWorld(world);
}

void MainWindow::on_actionFaster_triggered()
{
    updateSpeed(speed + 1);
}

void MainWindow::on_actionDouble_triggered()
{
    updateSpeed(speed * 2);
}

void MainWindow::on_actionHalf_triggered()
{
    updateSpeed(speed / 2);
}

void MainWindow::on_actionTogglePause_triggered()
{
    if (speed > 0) {
        old_speed = speed;
        updateSpeed(0);
    } else {
        updateSpeed(old_speed);
    }
}

void MainWindow::on_actionRestart_triggered()
{
    restart();
}

void MainWindow::on_graphicsView_mousePress(QMouseEvent *event)
{
    event->accept();

    Vec2 pt = ui->graphicsView->mapToWorld(event->pos());
    world->spawnRandomCreature(pt);
}
