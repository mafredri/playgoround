package main

import (
	"sync"
	"testing"
)

var ch = make(chan struct{}, 1000)

func init() {
	var j int
	go func() {
		for range ch {
			j++
		}
	}()
}

func BenchmarkChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}

func BenchmarkNormal(b *testing.B) {
	var j int
	var mu sync.Mutex
	for i := 0; i < b.N; i++ {
		mu.Lock()
		j++
		mu.Unlock()
	}
}
