# WorkerPool


[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/xinst/workerpool?status.svg)](https://godoc.org/github.com/xinst/workerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/xinst/workerpool)](https://goreportcard.com/report/github.com/xinst/workerpool)
[![Build Status](https://travis-ci.org/xinst/workerpool.svg?branch=master)](https://travis-ci.org/xinst/workerpool)

A customizable worker pool, do task async with workers.  
For simple use, you only need create your task struct and the Do method.

```go
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

func main() {

	// create a WorkerPool with default worker to do the task
	// with 100 taskQueueCap and 3 workers
	wp := workerpool.NewWorkerPoolWithDefault(100, 3, false, false)
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

```

For advancer you can custom your self worker manager and worker and the task

```go
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

```
