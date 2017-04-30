package test

import "testing"
import (
	tr "task-runner"
	"os"
	"time"
)

type TestTask struct {
	executed *bool
}

func NewTestTask(exec *bool) *TestTask {
	task := new(TestTask)
	task.executed = exec
	return task
}

func (t TestTask) Execute() {
	*t.executed = true
}

func TestTaskRunner_EnqueueTask(t *testing.T)  {
	no_worker := 3
	runner := tr.StartTaskRunner(no_worker, os.Stdout)

	result := false
	task := NewTestTask(&result)
	runner.EnqueueTask(*task)

	time.Sleep(time.Second)
	if ! result {
		t.Error("Test task not executed yet.")
	}

	runner.Stop()

	result = false
	task2 := NewTestTask(&result)
	runner.EnqueueTask(*task2)

	time.Sleep(time.Second)
	if result {
		t.Error("Enqueue should not execute if runner is stopped.")
	}
}


func TestTaskRunner_CreateBufferedTaskRunner(t *testing.T)  {
	no_worker := 3
	runner := tr.StartTaskRunner(no_worker, os.Stdout)

	task_chan := runner.CreateBufferedTaskRunner(10)
	if cap(task_chan) != 10 {
		t.Error("Size of task queue not correct")
	}

	result1 := false
	task1 := NewTestTask(&result1)

	result2 := false
	task2 := NewTestTask(&result2)

	result3 := false
	task3 := NewTestTask(&result3)

	task_chan <- task1
	task_chan <- task2
	task_chan <- task3

	time.Sleep(time.Second * 5) // Give some time to execute

	if ! result1 {
		t.Error("Test task 1 not executed yet.")
	}

	if ! result2 {
		t.Error("Test task 2 not executed yet.")
	}

	if ! result3 {
		t.Error("Test task 3 not executed yet.")
	}

	runner.Stop()
	time.Sleep(time.Second * 5) // Give some time to stop

	// Test size of task queue
	for i := 0; i < 10; i++ {
		t := NewTestTask(nil)
		task_chan <- t
	}

	if len(task_chan) != 10 {
		t.Errorf("Channel length mismatch. Len %d", len(task_chan))
	}
}
