package main

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

type task struct {
	id   int
	ok   bool
	done chan struct{}
}

func NewTask(id int) *task {
	d := make(chan struct{})
	return &task{id: id, done: d}
}

func (t *task) OK() bool {
	return t.ok
}

func (t *task) Start(chIn, chOut chan any) {
	var err error
	var comp C
	t.ok = true
	defer func() {
		t.ok = false
		<-t.done
		log.Println("end of task #", t.id)
	}()

	hh := NewHttpHelper().
		URL(URL).
		Param("lang", "ru")

LOOP:
	for {
		select {
		case <-t.done:
			fmt.Printf("task #%v get cancel\n", t.id)
			break LOOP
		case val, ok := <-chIn:
			to := 0
			max_req := 0
		LOOP_IIN:
			for {
				select {
				case <-t.done:
					fmt.Printf("task #%v get cancel\n", t.id)
					break LOOP
				default:
					if ok {
						iin := val.(string)
						hr := hh.Param("bin", iin).Get()
						//fmt.Println("Response:", hr)

						if hr.Err() != nil {
							err = hr.Err()
							fmt.Printf("task #%v:  GET: %v. Status: %v\n", t.id, err, hr.StatusCode)

							if e, ok := err.(*url.Error); ok && e.Timeout() && max_req < MAX_R {
								to += 3
								max_req++
								log.Printf("task #%v: req #%v: sleeping %v sec...\n", t.id, max_req, to)
								time.Sleep(time.Duration(to) * time.Second)
							} else if max_req < MAX_R {
								max_req++
								break LOOP_IIN
							} else {
								break LOOP
							}
						} else {
							if hr.StatusCode >= 400 {
								if max_req >= MAX_R {
									break LOOP_IIN
								}
								to += 3
								max_req++
								log.Printf("task #%v: Status: %v, req #%v: sleeping %v sec...\n", t.id, hr.Status, max_req, to)
								time.Sleep(time.Duration(to) * time.Second)
							} else if hr.OK() {
								hr.JSON(&comp)
								if err = hr.Err(); err != nil || !comp.Success {
									log.Printf("task #%v: error IIN %v to JSON err=%v. %v\n", t.id, iin, err, comp)
								} else {
									//mtx.Lock
									chOut <- comp
									log.Printf("task #%v: OK IIN %v to JSON \n", t.id, iin)
									//mtx.Unlock
								}
								break LOOP_IIN
							}
						}
					} else {
						break LOOP
					}
				}
			}
		}
	}
}

func (t *task) Stop() {
	if t.ok {
		t.done <- struct{}{}
	} else {
		fmt.Println("The task #", t.id, "stoped")

	}
}
