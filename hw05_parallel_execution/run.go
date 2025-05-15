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

	tasksCh := make(chan Task)                              // канал с задачами, из него воркеры будут брать задачи
	ctx, cancel := context.WithCancel(context.Background()) // позволяет остановить все горутины по команде
	defer cancel()                                          // гарантируем отмену

	var wg sync.WaitGroup

	var errCountMutex sync.Mutex // Защищает общий счётчик ошибок errCount
	var errCount int

	// запускаем N воркеров
	for i := 0; i < n; i++ {
		wg.Add(1) // в sync.WaitGroup запоминаем кол-во работающих потоков
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
						errCountMutex.Lock()
						errCount++
						if m > 0 && errCount >= m { //
							cancel()
						}
						errCountMutex.Unlock()
					}
				}
			}
		}()
	}

	// producer - заполняет канал задачами
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

	wg.Wait()

	if m > 0 && errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
