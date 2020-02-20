package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

const (
	numElements = 10000000
)

type key struct {
	k *string
}

var foo = map[key]int{}

func timeGC() {
	t := time.Now()
	runtime.GC()
	fmt.Printf("gc took: %s\n", time.Since(t))
}

func main() {
	for i := 0; i < numElements; i++ {
		kak := strconv.Itoa(i)
		foo[key{k: &kak}] = i
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}
}
