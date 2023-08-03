package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const TASK_TIMEOUT time.Duration = 10500 * time.Millisecond
const STEP_INCREASE_SLEEP int = 3

type Worker struct {
	HttpHelper *HttpHelper
	id         string
}

// type IWorker interface {
// 	Start(ctx context.Context, in, out chan any)
// }

type ClientPool struct {
	NumOfWorkers int
	workers      map[string]*Worker
	ChIn, ChOut  chan any
}

func NewClientPool(numOfWorkers int, chIn, chOut chan any) *ClientPool {
	var nameWorker string
	t := make(map[string]*Worker)
	for j := 0; j < numOfWorkers; j++ {
		nameWorker = fmt.Sprintf("w-%d", j+1)
		t[nameWorker] = NewWorker(nameWorker)
	}
	return &ClientPool{numOfWorkers, t, chIn, chOut}
}

func (cp *ClientPool) Add(id string, worker *Worker) {
	cp.workers[id] = worker
}

func (cp *ClientPool) Start(ctx context.Context) {
	defer fmt.Println("STOP pool.")
	defer close(cp.ChOut)
	var wg sync.WaitGroup

	for _, w := range cp.workers {
		wg.Add(1)
		go w.Start(ctx, cp.ChIn, cp.ChOut, &wg)
	}
	wg.Wait()

}

func NewWorker(id string) *Worker {
	return &Worker{NewHttpHelper(), id}
}

func (w *Worker) Start(ctx context.Context, chInt <-chan any, chOut chan<- any, wg *sync.WaitGroup) {
	defer fmt.Printf("!!! worker %v: STOPED... Number of GOROUTINES: %v\n", w.id, runtime.NumGoroutine())
	defer wg.Done()
	num := 0
	//TODO create connection instance for pass to the task (LongConnection)
	w.HttpHelper = w.HttpHelper.
		URL(URL).
		Param("lang", "ru")
LOOP:
	for {
		select {
		case i := <-chInt:
			ctxTask, _ := context.WithTimeout(ctx, TASK_TIMEOUT)
			num++
			nameTask := fmt.Sprintf("%v/%v", w.id, num)

			resp, e := HttpConn(ctxTask, w.HttpHelper, i.(string), nameTask) //, chanResponce)

			if e != nil {
				log.Printf("worker #%v: error IIN %v to JSON err=%v\n", w.id, i, e)
			} else {
				log.Printf("worker #%v: OK IIN %v to JSON %#v\n", w.id, i, resp.Obj)

				chOut <- resp
			}

		case <-ctx.Done():
			break LOOP
		}
	}
}

func HttpConn(ctx context.Context, hh *HttpHelper, iin string, name string) (C, error) { //, ch chan any) {
	var err error
	var comp C
	max_req, to := 0, 0
	hh = hh.Param("bin", iin)

	for max_req < MAX_R {
		hr := hh.Get(ctx)
		if hr.Err() != nil {
			err = hr.Err()
			if !isTimeout(err) { // not Timeout error
				break
			}
		} else {
			if hr.OK() {
				hr.JSON(&comp)
				if err = hr.Err(); err != nil {
					err = fmt.Errorf("task #%v: error IIN %v to JSON err=%v. %v", name, iin, err, comp)
				} else if !comp.Success {
					err = fmt.Errorf("task #%v: IIN %v Company not success. %v", name, iin, comp)
				}
				break
			} else if hr.StatusCode < 400 { // status code not timeout
				break
			}

		}
		to += STEP_INCREASE_SLEEP
		max_req++
		log.Printf("task #%v: req #%v: sleeping %v sec...\n", name, max_req, to)
		time.Sleep(time.Duration(to) * time.Second)
	}
	return comp, err
}
func isTimeout(err error) bool {
	type Timeout interface {
		Timeout() bool
	}
	e, ok := err.(Timeout)
	return ok && e.Timeout()
}

func LongConnection(ctx context.Context, i int, name string, ch chan any) {
	//the immitation long connection
	defer fmt.Println("--------------END CONN", name)
	d := i*100 + 50
	fmt.Printf("#### task %v: start %vms\n", name, d)
	//working...
	time.Sleep(time.Duration(d) * time.Millisecond)
	ch <- fmt.Sprintf("#### RESULT t#%v: %v ms", name, d)

	fmt.Printf("### task %v: end %v ms\n", name, d)
}
