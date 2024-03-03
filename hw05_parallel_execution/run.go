package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var (
	countError  int64
	taskCounter int64
)

var (
	isErrorChan   = make(chan bool)
	isAllTaskDone = make(chan bool)
	doneTaskChan  = make(chan bool)
)

func runTask(t Task) {
	err := t()
	if err != nil {
		atomic.AddInt64(&countError, 1)
	}
	doneTaskChan <- true
	atomic.AddInt64(&taskCounter, -1)
}

func taskDoneChecker(taskSize int) {
	doneTaskCount := 0
	for {
		<-doneTaskChan
		doneTaskCount++
		if doneTaskCount == taskSize {
			isAllTaskDone <- true
			return
		}
	}
}

func checkErrorCount(m int64) bool {
	for {
		if atomic.LoadInt64(&countError) >= m {
			isErrorChan <- true
		}
	}
}

func taskLimiter(tasks []Task, n int64, taskSize int) {
	nbTask := 0
	for {
		if atomic.LoadInt64(&taskCounter) < n {
			atomic.StoreInt64(&taskCounter, 1)
			if nbTask == taskSize-1 {
				return
			}
			go runTask(tasks[nbTask])
			nbTask++
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	go taskDoneChecker(len(tasks))
	go checkErrorCount(int64(m))
	go taskLimiter(tasks, int64(n), len(tasks))
	select {
	case <-isErrorChan:
		return ErrErrorsLimitExceeded
	case <-isAllTaskDone:
		return nil
	}
}
