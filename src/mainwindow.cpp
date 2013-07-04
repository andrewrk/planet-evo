#include "mainwindow.h"
#include "ui_mainwindow.h"

#include <QDebug>
#include <QTimer>

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    ui->actionExit->setShortcut(QKeySequence("Alt+F4"));
    ui->graphicsView->setScene(&scene);

    uint seed = 1234;
    qsrand(seed);
    qDebug() << "Using seed" << seed;
    world = new World(&scene);

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
