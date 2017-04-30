package task_runner

import (
	"os"
	"testing"
	"time"
)

func TestStartTaskRunner(t *testing.T)  {
	no_worker := 3
	runner := StartTaskRunner(no_worker, os.Stdout)

	if len(runner.workerMap) != no_worker {
		t.Errorf("Incorrect worker map len Expected: %d Found: %d",
			no_worker, len(runner.workerMap))
	}

	if cap(runner.workerQueue) != no_worker {
		t.Errorf("Incorrect worker queue cap Expected: %d Found: %d",
			no_worker, cap(runner.workerQueue))
	}

	runner.Stop()
	time.Sleep(time.Second * 5) // Give 5 secs to gracefully exit.
	if ! runner.stopped {
		t.Error("Task runner not stopped yet.")
	}
}

func TestTaskRunner_CreateBufferedTaskRunner(t *testing.T)  {
	no_worker := 3
	runner := StartTaskRunner(no_worker, os.Stdout)

	// Currently only 1000 are supported.
	// Override using env variable
	task_chan := runner.CreateBufferedTaskRunner(1000000)
	if task_chan != nil {
		t.Errorf("Expected to fail instead returned channel of len %d",
			cap(task_chan))
	}

	task_chan = runner.CreateBufferedTaskRunner(100)

	if cap(task_chan) != 100 {
		t.Errorf("Expected channel capacity: %d Found: %d",
			100, cap(task_chan))
	}

	runner.Stop()
	time.Sleep(time.Second * 5) // Give 5 secs to gracefully exit.
	if ! runner.stopped {
		t.Error("Task runner not stopped yet.")
	}

	runner2 := StartTaskRunner(no_worker, os.Stdout)
	os.Setenv("OVERRIDE_QUEUE_SIZE", "1")
	task_chan2 := runner2.CreateBufferedTaskRunner(1500)

	if cap(task_chan2) != 1500 {
		t.Errorf("Expected channel capacity: %d Found: %d",
			 1500, cap(task_chan2))
	}
}