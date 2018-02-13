package workerpool

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	// ErrTaskQueueIsFull is over the channel cap
	ErrTaskQueueIsFull = errors.New("task queue is full")
)

// WorkerPool defined a pool with specific worker to do the task async
type WorkerPool struct {
	shutDownChannel chan struct{} //Channel to shutdown the workerpool
	taskChannel     chan Task     //task Channel,buffered
	queuedTask      int32         // The number of task queued

	taskQueueCap uint32         // The max number of task can store in the queue.
	workerMgr    *WorkerManager // WorkerManager
}

// NewWorkerPoolWithDefault new a WorkerPool with DefaultWorker
func NewWorkerPoolWithDefault(taskQueueCapacity uint32, workerSize uint32, waitIfNoWorker, allowTempExceedCap bool) *WorkerPool {

	if taskQueueCapacity == 0 || workerSize == 0 {
		return nil
	}

	wp := &WorkerPool{
		shutDownChannel: make(chan struct{}),
		taskChannel:     make(chan Task, taskQueueCapacity),
		taskQueueCap:    taskQueueCapacity,
		workerMgr:       NewWorkerMgr(workerSize, waitIfNoWorker, allowTempExceedCap, nil),
	}
	return wp
}

// NewWorkerPool new a WorkerPool
func NewWorkerPool(taskQueueCapacity uint32, workerMgr *WorkerManager) *WorkerPool {

	wp := &WorkerPool{
		shutDownChannel: make(chan struct{}),
		taskChannel:     make(chan Task, taskQueueCapacity),
		taskQueueCap:    taskQueueCapacity,
		workerMgr:       workerMgr,
	}
	return wp
}

// PushTask put the task into taskList
func (wp *WorkerPool) PushTask(task Task) error {

	if len(wp.taskChannel) >= int(wp.taskQueueCap) {
		return ErrTaskQueueIsFull
	}
	wp.taskChannel <- task
	atomic.AddInt32(&wp.queuedTask, 1)

	return nil
}

// Start return error is there something error
func (wp *WorkerPool) Start() error {
	// check the work pool, if something err return
	go func() {
		fmt.Println("workpool is start")
		for {
			select {
			case <-wp.shutDownChannel:
				return
			default:
				worker, err := wp.workerMgr.GetWorker()
				if err != nil {
					fmt.Println(err)
					return
				}

				if worker != nil {
					t, isClosed := <-wp.taskChannel
					if !isClosed {
						fmt.Println("task channel has benn closed")
						wp.workerMgr.PutWorker(worker)
						return
					}

					atomic.AddInt32(&wp.queuedTask, -1)
					go func(task Task, woker Worker, mgr *WorkerManager) {
						worker.Work(task)
						mgr.PutWorker(worker)
					}(t, worker, wp.workerMgr)
				} else {
					wp.workerMgr.PutWorker(worker)
				}
			}
		}
	}()

	return nil
}

// Stop the pool
func (wp *WorkerPool) Stop() {
	close(wp.shutDownChannel)
	wp.workerMgr.Close()
}
