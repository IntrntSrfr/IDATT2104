package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"
)

type PrimeList struct {
	nums []int
}

func main() {
	var (
		start   int
		end     int
		threads int
	)
	flag.IntVar(&start, "start", 0, "value to start searching from")
	flag.IntVar(&end, "end", 10, "value to end search at")
	flag.IntVar(&threads, "threads", 10, "amount of threads to use")
	flag.Parse()

	var (
		list        = &PrimeList{nums: make([]int, 0)}
		numCh       = make(chan int)
		ctx, cancel = context.WithCancel(context.Background())
		wg          = &sync.WaitGroup{}
		wg2         = &sync.WaitGroup{}
	)

	// run background listener for new primes
	// it terminates when cancel is called
	wg2.Add(1)
	go appendNums(numCh, ctx, list, wg2)

	wg.Add(threads)
	startTime := time.Now()
	for i := 0; i < threads; i++ {
		go findPrimes(start + i, end, threads, numCh, wg)
	}
	wg.Wait() // wait for all prime threads
	cancel()
	wg2.Wait()
	fmt.Println("Time taken:", time.Since(startTime))
	fmt.Println("Prime numbers found:", list.nums)
}

func appendNums(ch <-chan int, ctx context.Context, list *PrimeList, wg *sync.WaitGroup) {
	for {
		select {
		case n := <-ch:
			fmt.Println(n)
			list.nums = append(list.nums, n)
		case <-ctx.Done():
			sort.Sort(sort.IntSlice(list.nums))
			wg.Done()
			return
		}
	}
}

func findPrimes(start, end, delta int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i <= end; i += delta {
		if isPrime(i) {
			ch <- i
		}
	}
}

func isPrime(n int) bool {
	if n < 0 || n == 0 || n == 1 {
		return false
	}

	if n == 2 {
		return true
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
