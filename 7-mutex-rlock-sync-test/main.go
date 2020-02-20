package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.RWMutex
	mu.Lock()

	count := 10000
	wg.Add(count)

	times := make([]time.Time, count, count)
	for i := 0; i < count; i++ {
		i := i
		go func() {
			mu.RLock()
			defer mu.RUnlock()
			t := time.Now()
			times[i] = t
			wg.Done()
		}()
	}

	mu.Unlock()

	timesCnt := make(map[time.Time]int)
	for _, t := range times {
		timesCnt[t]++
	}
	for t, cnt := range timesCnt {
		if cnt > 1 {
			fmt.Println(t.String(), cnt)
		}
	}
}
