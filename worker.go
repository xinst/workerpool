package workerpool

import (
    "errors"
    "sync"
    "sync/atomic"
)

// Task interface is only a common interface, you need to impl your Task such as download task
// with the Do method
type Task interface {
    Do() error
}

// Worker can defined to be a specific Worker，such as download worker
// It only need to impl the work interface and Close interface
type Worker interface {
    Work(Task)
    Close() error
}

// DefaultWorker is empty
type DefaultWorker struct {
}

// Work call the task's Do method
func (dw *DefaultWorker) Work(task Task) {
    task.Do()
}

// Close return nil
func (dw *DefaultWorker) Close() error {
    return nil
}

// WorkerManager create the specific worker to do your task
type WorkerManager struct {
    // freeList the free worker list
    freeList chan Worker

    // hold the worker size
    used int32

    // the fun to new a worker
    newFun func() Worker

    //Allow temporarily exceed cap
    allowTempExceedCap bool

    // sync.group wait
    wg  sync.WaitGroup
}

// NewWorkerMgr create a manager with workerSize, if possible it will create a new worker temporarily
// the newFun defined a func to new a worker impl the worker interface, if it is nil, will use DefaultWorker
// if the waitIfNoWorker is true, it will block at GetWorker() until there is a available worker
// if the allowTempExceedCap is true, it will temporarily out of size the workerSize, it is not under control
func NewWorkerMgr(workerSize uint32, allowTempExceedCap bool, newFun func() Worker) *WorkerManager {
    wm := &WorkerManager{
        freeList:           make(chan Worker, workerSize),
        used:               0,
        allowTempExceedCap: allowTempExceedCap,
    }

    if newFun == nil {
        wm.newFun = func() Worker {
            return &DefaultWorker{}
        }
    } else {
        wm.newFun = newFun
    }

    for index := 0; index < int(workerSize); index++ {
        wm.freeList <- wm.newFun()
    }

    return wm
}

// GetWorker will return an available worker or an error is there is no worker can be use
// or the channel is closed
func (wm *WorkerManager) GetWorker() (Worker, error) {
    // if wm.waitIfNoWorker {
    // 	s, more := <-wm.freeList
    // 	if s == nil || !more {
    // 		return nil, errors.New("pool has been closed")
    // 	}
    // 	return s, nil
    // }

    select {
    case s, more := <-wm.freeList:
        if s == nil || !more {
            return nil, errors.New("pool has been closed")
        }
        atomic.AddInt32(&wm.used, 1)
        wm.wg.Add(1)
        return s, nil
    default:
        if wm.allowTempExceedCap {
            atomic.AddInt32(&wm.used, 1)
            wm.wg.Add(1)
            return wm.newFun(), nil
        }
    }

    return nil, nil
}

// PutWorker can put back the worker under the manager's control
// it is important if you set allowTempExceedCap is true
func (wm *WorkerManager) PutWorker(worker Worker) {
    atomic.AddInt32(&wm.used, -1)
    wm.wg.Done()

    if worker == nil {
        return
    }

    select {
    case wm.freeList <- worker:
    default:
        worker.Close()
    }
}

// Close can close the Channel
func (wm *WorkerManager) Close() {
    close(wm.freeList)
    wm.wg.Wait()
}
