package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"

	"github.com/xinst/workerpool"
)

// DownloadTask defined your task
type DownloadTask struct {
	index int
	wg *sync.WaitGroup
}

// Do method is implied the interface of workerpool.Task
// you can define a result channel to recv the task result if you like
func (dt *DownloadTask) Do() error {
	defer dt.wg.Done()
	// random sleep Seconds
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Printf("%d Download task done\n", dt.index)
	return nil
}

// DownloadWorker struct
type DownloadWorker struct {
}

// Work fun implied
func (dw *DownloadWorker) Work(task workerpool.Task) {
	//fmt.Printf("Download worker do %#v \n", task)
	task.Do()
}

// Close func will be called when the worker is destroy
func (dw *DownloadWorker) Close() error {
	return nil
}

func advanceTest() {

	mgr := workerpool.NewWorkerMgr(3, false, func() workerpool.Worker {
		return &DownloadWorker{}
	})

	// create a WorkerPool with default worker to do the task
	// with 100 taskQueueCap and 3 workers
	wp := workerpool.NewWorkerPool(100, mgr)
	wp.Start()

	var wg sync.WaitGroup
	wg.Add(10)

	// push 10 task
	for index := 0; index < 10; index++ {

		task := &DownloadTask{
			index: index,
			wg : &wg,
		}
		wp.PushTask(task)
	}
	wg.Wait()
	fmt.Println("all tasks have done")
}

func main() {
	advanceTest()
}
