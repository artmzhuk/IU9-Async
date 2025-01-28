package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const numThreads = 10
const numNumbers = 100

type node struct {
	val    int
	thread int
	next   *node
}

func checkExists(val int, list *node) bool {
	current := list
	for current != nil {
		if current.val == val {
			return true
		}
		current = current.next
	}
	return false
}

func writeNumber(val, thread int, list *node) {
	current := list
	for current.next != nil {
		current = current.next
	}
	current.next = &node{
		val:    val,
		thread: thread,
		next:   nil,
	}
}

func worker(wg *sync.WaitGroup, readMtx, writeMtx *sync.Mutex, thread int, list *node) {
	defer wg.Done()
	for i := 0; i < numNumbers; i++ {
		number := rand.Intn(1000)
		readMtx.Lock()
		exists := checkExists(number, list)
		readMtx.Unlock()
		if !exists {
			writeMtx.Lock()
			if !checkExists(number, list) {
				writeNumber(number, thread, list)
			}
			writeMtx.Unlock()
		} else {
			continue
		}
	}
}

func main() {
	writeMtx := sync.Mutex{}
	readMtx := sync.Mutex{}
	wg := &sync.WaitGroup{}
	list := &node{
		val:  0,
		next: nil,
	}
	wg.Add(numThreads)
	startTime := time.Now()
	for i := 0; i < numThreads; i++ {
		go worker(wg, &readMtx, &writeMtx, i, list)
	}
	wg.Wait()
	resTime := time.Since(startTime)
	isPresent := make(map[int]bool, 1000)
	listLength := 0
	lastRoutine := -1
	res := ""
	current := list
	for current != nil {
		if isPresent[current.val] {
			panic("smth wrong")
		}
		isPresent[current.val] = true
		fmt.Printf("%d(%d) ", current.val, current.thread)
		if lastRoutine != current.thread {
			res += strconv.Itoa(current.thread) + " -> "
			lastRoutine = current.thread
		}
		current = current.next
		listLength++

	}
	fmt.Println()
	fmt.Println("длина списка:", listLength, "из", numNumbers*numThreads, resTime)
	fmt.Println(res)
}
