package main

import (
	"context"
	"fmt"
)

type ITask interface {
	Start(in, out chan any)
	Stop()
}

type ClientPool struct {
	Tasks       map[int]ITask
	chIn, chOut chan any
	Ctx         context.Context
}

func NewClientPool(ctx context.Context, chIn, chOut chan any) *ClientPool {
	t := make(map[int]ITask)
	return &ClientPool{t, chIn, chOut, ctx}
}

func (cp *ClientPool) Add(id int, task ITask) {
	cp.Tasks[id] = task
}

func (cp *ClientPool) Start() {
	go func() {
		<-cp.Ctx.Done()
		cp.Stop()
		close(cp.chOut)
	}()

	for k, v := range cp.Tasks {
		go v.Start(cp.chIn, cp.chOut)
		fmt.Println("pool: Start task #", k)
	}
}

func (cp *ClientPool) Stop() {
	defer fmt.Println("End pool stop")
	for k, v := range cp.Tasks {
		v.Stop()
		fmt.Println("pool: Stop task #", k)
	}

}
