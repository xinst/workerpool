package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xinst/workerpool"
)

// custom your Task
// =============================================
// DownloadTask defined your task
type DownloadTask struct {
	index int
}

// Do method is implied the interface of workerpool.Task
// you can define a result channel to recv the task result if you like
func (dt *DownloadTask) Do() error {
	// random sleep Seconds
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Printf("%d Download task done\n", dt.index)
	return nil
}

// custom your worker
// =============================================
// Download worker
type DownloadWorker struct {
}

// Work fun implied
func (dw *DownloadWorker) Work(task workerpool.Task) {
	fmt.Printf("Download worker do %#v \n", task)
	task.Do()
}

// Close func will be called when the worker is destory
func (dw *DownloadWorker) Close() error {
	return nil
}

func main() {

	mgr := workerpool.NewWorkerMgr(3, false, false, func() workerpool.Worker {
		return &DownloadWorker{}
	})

	// create a WorkerPool with default worker to do the task
	// with 100 taskQueueCap and 3 workers
	wp := workerpool.NewWorkerPool(100, mgr)
	wp.Start()

	// push 10 task
	for index := 0; index < 10; index++ {

		task := &DownloadTask{
			index: index,
		}
		wp.PushTask(task)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sig
}
