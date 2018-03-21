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
		wp.PushTask(task,true)
	}
	wg.Wait()
	fmt.Println("all tasks have done")
}

func main() {
	simpleTest()
}
