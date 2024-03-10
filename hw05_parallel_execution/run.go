package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func runTask(t Task, errorChan chan<- struct{}, successChan chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	err := t()
	if err != nil {
		errorChan <- struct{}{}
	} else {
		successChan <- struct{}{}
	}
}

func startTask(t Task, wg *sync.WaitGroup, errorChan chan<- struct{}, successChan chan<- struct{}, nbTask *int) {
	wg.Add(1)
	go runTask(t, errorChan, successChan, wg)
	*nbTask++
}

func taskLimiter(tasks []Task, done chan<- error, wg *sync.WaitGroup, maxErrorCount int, maxGoCount int) {
	errorChan := make(chan struct{}, maxGoCount*2)
	successChan := make(chan struct{}, maxGoCount*2)
	var wgTasks sync.WaitGroup
	defer func() {
		wgTasks.Wait()
		close(errorChan)
		close(successChan)
		wg.Done()
	}()
	nbTask := 0
	succesCount := 0
	errorCount := 0
	taskSize := len(tasks)
	for i := 0; i < maxGoCount && i < taskSize; i++ {
		startTask(tasks[nbTask], &wgTasks, errorChan, successChan, &nbTask)
	}
	for {
		select {
		case <-errorChan:
			errorCount++
		case <-successChan:
			succesCount++
		}

		if errorCount == maxErrorCount {
			done <- ErrErrorsLimitExceeded
			return
		}

		if succesCount+errorCount == taskSize {
			done <- nil
			return
		}

		if nbTask < taskSize {
			startTask(tasks[nbTask], &wgTasks, errorChan, successChan, &nbTask)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	wg.Add(1)
	resultChan := make(chan error)
	defer func() {
		close(resultChan)
	}()
	go taskLimiter(tasks, resultChan, &wg, m, n)
	x := <-resultChan
	wg.Wait()
	return x
}
