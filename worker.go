package task_runner

import "log"
/*
 * Worker type.
 * The task queue is the incoming tasks channel.
 * The worker queue is the queue of workers it will add itself to.
 * Stop channel to stop the worker.
 */
type Worker struct {
	ID int
	TaskQueue chan Task
	WorkerQueue chan chan Task
	StopSignal chan bool
	log *log.Logger
}

func NewWorker(id int, workerQueue chan chan Task, logr *log.Logger) Worker {
	worker := Worker{
		ID: id,
		TaskQueue: make(chan Task),
		WorkerQueue: workerQueue,
		StopSignal: make(chan bool),
		log: logr,
	}
	return worker
}

func (w *Worker) Start () {
	go func() {
		w.log.Println("Worker starting")
		for {
			w.log.Println("Looking for work")
			w.WorkerQueue <- w.TaskQueue
			select {
			case task := <- w.TaskQueue:
				w.log.Println("Got work")
				task.Execute()
			case <- w.StopSignal:
				w.log.Println("Stopping")
				return
			}
		}
	} ()
}

func (w * Worker) Stop()  {
	w.log.Println("Got stop signal")
	go func() {
		w.StopSignal <- true
	} ()
}