# WorkerPool


[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-green.svg)](https://www.gnu.org/licenses/gpl-3.0) [![GoDoc](https://godoc.org/github.com/xinst/workerpool?status.svg)](https://godoc.org/github.com/xinst/workerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/xinst/workerpool)](https://goreportcard.com/report/github.com/xinst/workerpool)
[![Build Status](https://travis-ci.org/xinst/workerpool.svg?branch=master)](https://travis-ci.org/xinst/workerpool)

A customizable worker pool, do task async with workers.  
For simple use, you only need create your task struct and the Do method.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"

	"github.com/xinst/workerpool"
)

// SimpleTask defined your task
type SimpleTask struct {
	index int
	wg *sync.WaitGroup
}

// Do method is implied the interface of workerpool.Task
// you can define a result channel to recv the task result if you like
func (st *SimpleTask) Do() error {
	defer st.wg.Done()
	// random sleep Seconds
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Printf("%d SimpleTask task done\n", st.index)
	return nil
}

func simpleTest() {
	// create a WorkerPool with default worker to do the task
	// with 100 taskQueueCap and 3 workers
	wp := workerpool.NewWorkerPoolWithDefault(100, 3, false)
	wp.Start()

	var wg sync.WaitGroup
	wg.Add(10)
	// push 10 task
	for index := 0; index < 10; index++ {

		task := &SimpleTask{
			index: index,
			wg : &wg,
		}
		wp.PushTask(task)
	}
	wg.Wait()
	fmt.Println("all tasks have done")
}

func main() {
	simpleTest()
}


```

For advancer you can custom your self worker manager and worker and the task

```go
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


```
