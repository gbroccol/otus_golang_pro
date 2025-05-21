package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		n = 1
	}

	tasksCh := make(chan Task)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case tasksCh <- task:
			}
		}
	}()

	errCount := startWorkers(ctx, n, m, tasksCh, cancel)

	if m > 0 && errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func startWorkers(ctx context.Context, n, m int, tasksCh <-chan Task, cancel context.CancelFunc) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCount := 0

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-tasksCh:
					if !ok {
						return
					}
					if err := task(); err != nil {
						mu.Lock()
						errCount++
						if m > 0 && errCount >= m {
							cancel()
						}
						mu.Unlock()
					}
				}
			}
		}()
	}

	wg.Wait()
	return errCount
}
