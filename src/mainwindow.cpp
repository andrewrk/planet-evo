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
    ui->graphicsView->setScene(&scene);

    restart();

    QTimer *timer = new QTimer(this);
    connect(timer, SIGNAL(timeout()), this, SLOT(stepWorld()));
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
    scene.clear();
    world = new World(&scene);

    // add some bouncy carbon particles
    for (int i = 0; i < 20; i++) {
        Vec2 pos(randRange(0.0, world->size.x), randRange(0.0, world->size.y));
        Particle *p = new Particle(CarbonParticle, pos);
        p->vel = Vec2::direction(randf() * 2 * PI);
        p->vel.setLength(randRange(-0.2, 0.2));
        world->addParticle(p);
    }
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
