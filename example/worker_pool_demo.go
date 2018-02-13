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

// SimpleTask defined your task
type SimpleTask struct {
	index int
}

// Do method is implied the interface of workerpool.Task
// you can define a result channel to recv the task result if you like
func (dt *SimpleTask) Do() error {
	// random sleep Seconds
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Printf("%d SimpleTask task done\n", dt.index)
	return nil
}

func simpleTest() {
	// create a WorkerPool with default worker to do the task
	// with 100 taskQueueCap and 3 workers
	wp := workerpool.NewWorkerPoolWithDefault(100, 3, false, false)
	wp.Start()

	// push 10 task
	for index := 0; index < 10; index++ {

		task := &SimpleTask{
			index: index,
		}
		wp.PushTask(task)
	}
}

func main() {

	simpleTest()
	//advanceTest()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sig
}
