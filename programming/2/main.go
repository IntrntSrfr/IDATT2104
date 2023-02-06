package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	workers := NewPool(4)
	loop := NewPool(1)

	workers.Start()
	loop.Start()

	workers.Post(func() {
		fmt.Println("This is task A")
	})

	workers.Post(func() {
		fmt.Println("This is task B")
	})

	workers.PostAfter(3000, func() {
		fmt.Println("This is task E")
	})

	loop.Post(func() {
		fmt.Println("This is task C")
	})

	loop.Post(func() {
		fmt.Println("This is task D")
	})

	loop.PostAfter(4000, func() {
		fmt.Println("This is task F")
	})

	// stop after 2 seconds, if nothing else stops it we have a deadlock
	time.AfterFunc(time.Second*5, func() {
		workers.Stop()
		loop.Stop()
	})

	workers.Wait()
	loop.Wait()
	fmt.Println("i am now finished yapperrooo")
}

type Pool struct {
	*sync.Cond
	n     int
	tasks []func()
	wg    *sync.WaitGroup
	stop  bool
}

func NewPool(n int) *Pool {
	return &Pool{
		Cond:  sync.NewCond(&sync.Mutex{}),
		n:     n,
		tasks: []func(){},
		wg:    &sync.WaitGroup{},
	}
}

func (p *Pool) Start() {
	p.wg.Add(p.n)
	for i := 0; i < p.n; i++ {
		go func(i int) {
			defer p.wg.Done()
			for {
				p.Cond.L.Lock()
				for len(p.tasks) == 0 {
					p.Cond.Wait()
					if p.stop {
						p.Cond.L.Unlock()
						return
					}
				}

				var f func()
				if len(p.tasks) > 0 {
					f = p.tasks[0]
					p.tasks = p.tasks[1:]
				}
				p.Cond.L.Unlock()

				if f != nil {
					f()
				}
			}
		}(i)
	}
}

func (p *Pool) Stop() {
	p.L.Lock()
	p.stop = true
	p.L.Unlock()
	p.Cond.Broadcast()
}

func (p *Pool) Post(f func()) {
	p.L.Lock()
	p.tasks = append(p.tasks, f)
	p.Signal()
	p.L.Unlock()
}

func (p *Pool) PostAfter(t time.Duration, f func()) {
	time.AfterFunc(time.Millisecond*t, func() { p.Post(f) })
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
