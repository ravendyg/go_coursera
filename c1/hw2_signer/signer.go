package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ExecutePipeline -
func ExecutePipeline(jobs ...job) {
	wg := sync.WaitGroup{}
	var in chan interface{}
	out := make(chan interface{})

	for _, jb := range jobs {
		localJob := jb
		wg.Add(1)

		go func(int_in, int_out chan interface{}) {
			localJob(int_in, int_out)
			close(int_out)
			wg.Done()
		}(in, out)

		// replace channels
		in = out
		out = make(chan interface{})
	}

	wg.Wait()
}

// SingleHash -
func SingleHash(in, out chan interface{}) {
	wgOut := sync.WaitGroup{}

	for input := range in {
		wgOut.Add(1)
		data := fmt.Sprintf("%d", input)
		md5 := DataSignerMd5(data)
		go func() {
			wg := sync.WaitGroup{}
			wg.Add(2)

			var part1 string
			var part2 string
			go func() {
				part1 = DataSignerCrc32(data)
				wg.Done()
			}()
			go func() {
				part2 = DataSignerCrc32(md5)
				wg.Done()
			}()

			wg.Wait()
			result := part1 + "~" + part2
			out <- result
			wgOut.Done()
		}()
	}

	wgOut.Wait()
}

// MultiHash -
func MultiHash(in, out chan interface{}) {
	storage := make(map[int]string)
	var mutex = &sync.Mutex{}
	wgOut := sync.WaitGroup{}
	var k int

	// in case one of runCalculations would take longer
	ensureRightOrder := func(acc *[]string, j int) {
		mutex.Lock()

		storage[j] = strings.Join(*acc, "")

		keys := make([]int, 0)
		for m := range storage {
			keys = append(keys, m)
		}
		sort.Ints(keys)
		key := keys[0]
		out <- storage[key]
		delete(storage, key)

		mutex.Unlock()
	}

	runCalculations := func(data string, j int) {
		wg := sync.WaitGroup{}

		acc := make([]string, 6)
		for i := 0; i <= 5; i++ {
			wg.Add(1)
			go func(j int) {
				c := strconv.Itoa(j) + data
				acc[j] = DataSignerCrc32(c)
				wg.Done()
			}(i)
		}
		wg.Wait()

		ensureRightOrder(&acc, j)

		wgOut.Done()
	}

	for input := range in {
		wgOut.Add(1)

		go runCalculations(input.(string), k)

		k++
	}

	wgOut.Wait()
}

// CombineResults -
func CombineResults(in, out chan interface{}) {
	acc := make([]string, 0)
	for input := range in {
		data := input.(string)
		acc = append(acc, data)
	}
	sort.Strings(acc)
	out <- strings.Join(acc, "_")
}
