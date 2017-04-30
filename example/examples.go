package main

import (
	tr "task-runner"
	"fmt"
	"time"
	"strconv"
	"os"
)

var printString = func(intVal int) string {
	return strconv.Itoa(intVal)
}

type AdditionTask struct {
	x int
	y int
	result int
}

func (t AdditionTask) Execute() {
	fmt.Println("Executing Addition Task")
	fmt.Println("Adding " + printString(t.x) + " and " + printString(t.y))
	fmt.Println("Result is " + printString(t.x + t.y))
}

func NewAdditionTask(x, y int) *AdditionTask {
	task := new(AdditionTask)
	task.x = x
	task.y = y
	return task
}

type MultiplicationTask struct {
	x int
	y int
	result int
}

func (t MultiplicationTask) Execute() {
	fmt.Println("Executing Multiplication Task")
	fmt.Println("Multiplying " + printString(t.x) + " and " + printString(t.y))
	fmt.Println("Result is " + printString(t.x * t.y))
}

func NewMultiplicationTask(x, y int) *MultiplicationTask {
	task := new(MultiplicationTask)
	task.x = x
	task.y = y
	return task
}

func main()  {
	runner := tr.StartTaskRunner(2, os.Stdout)
	for i := 0; i < 10; i++ {
		t1 := NewAdditionTask(i, (10 - i))
		t2 := NewMultiplicationTask(i, (10 - i))
		runner.EnqueueTask(*t1)
		runner.EnqueueTask(*t2)
	}

	task_chan := runner.CreateBufferedTaskRunner(5)
	t1 := NewAdditionTask(5, 5)
	t2 := NewMultiplicationTask(5, 5)
	task_chan <- t1
	task_chan <- t2

	// Give 20 seconds to finish
	time.Sleep(time.Second * 20)

	runner.Stop()

	// Give 5 secs to gracefully stop
	time.Sleep(time.Second * 5)
}
