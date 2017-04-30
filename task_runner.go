package task_runner

import (
	"io"
	"log"
	"os"
)

// Structure of Task which we will be able to run
type Task interface {
	Execute()
}

type TaskRunner struct {
	// Queue of workers
	workerQueue chan chan Task
	// Channel for stop signal
	stopSignal chan bool
	// Buffered TaskRunner mode channel
	taskChannel chan Task
	// Map of workers currently working
	workerMap map [int] *Worker
	// No of workers
	noOfWorkers int
	// Stop receiving tasks
	stopped bool
	// Logger
	err_log *log.Logger
	info_log *log.Logger
}

func StartTaskRunner(no_of_workers int, logger_out io.Writer) *TaskRunner {

	taskRunner := new(TaskRunner)
	taskRunner.noOfWorkers = no_of_workers
	taskRunner.workerQueue = make(chan chan Task, no_of_workers)
	taskRunner.workerMap = make(map[int] *Worker, no_of_workers)
	taskRunner.stopSignal = make(chan bool)
	taskRunner.stopped = false

	// Init logger
	taskRunner.initLogger(logger_out)

	for i := 1; i <= no_of_workers; i++ {
		worker := NewWorker(i, taskRunner.workerQueue, taskRunner.info_log)
		worker.Start()
		taskRunner.workerMap[i] = &worker
	}

	return taskRunner
}

func (runner *TaskRunner) initLogger(logger_op io.Writer)  {
	runner.err_log = log.New(logger_op, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	runner.info_log = log.New(logger_op, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func (runner *TaskRunner) EnqueueTask(newTask Task) {
	if runner.stopped {
		runner.err_log.Println("TaskRunner not running")
		return
	}
	runner.info_log.Println("Got new task")
	go func(task Task) {
		worker := <-runner.workerQueue
		runner.info_log.Println("Assigning task to worker")
		worker <- task
	} (newTask)
}

func (runner *TaskRunner) Stop() {
	runner.info_log.Println("Got stop request")
	go func(t *TaskRunner) {
		runner.info_log.Println("Stopping workers")
		for i := 1; i <= t.noOfWorkers; i++ {
			t.workerMap[i].Stop()
		}
		runner.stopped = true
		t.stopSignal <- true
	} (runner)
}

func (runner *TaskRunner) CreateBufferedTaskRunner(size int) chan Task {
	if size > 1000 && ! (os.Getenv("OVERRIDE_QUEUE_SIZE") == "1") {
		runner.err_log.Println("Queue size too big")
		return nil
	}
	runner.taskChannel = make(chan Task, size)
	runner.startBufferedTaskRunner()
	return runner.taskChannel
}

func (runner *TaskRunner) startBufferedTaskRunner() {
	go func(_runner *TaskRunner) {
		_runner.info_log.Println("Starting listener for Tasks")
		for {
			select {
			case task := <- _runner.taskChannel:
				_runner.info_log.Println("Got new task on channel")
				go func(t Task) {
					worker := <- _runner.workerQueue
					_runner.info_log.Println("Assigning worker")
					worker <- t
				} (task)
			case <- _runner.stopSignal:
				return
			}


		}
	} (runner)
}