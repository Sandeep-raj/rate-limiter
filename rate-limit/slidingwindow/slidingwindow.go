package slidingwindow

import (
	"fmt"
	"sync"
	"time"
)

type slidingwindow struct {
	interval   int
	maxreqs    int
	window     map[string]int32
	windowLock *sync.Mutex
}

func SlidingWindowRateLimiter(interval int, maxReqs int) slidingwindow {
	window := make(map[string]int32)
	currTime := time.Now().Unix() / int64(interval)
	window[fmt.Sprintf("%d", currTime)] = 0

	return slidingwindow{
		interval:   interval,
		maxreqs:    maxReqs,
		window:     window,
		windowLock: &sync.Mutex{},
	}
}

func (swrl *slidingwindow) Allow() bool {
	swrl.windowLock.Lock()
	defer swrl.windowLock.Unlock()

	currWindow := time.Now().Unix() / int64(swrl.interval)
	prevWindow := currWindow - 1

	e := float32(time.Now().Unix()%int64(swrl.interval)) / float32((swrl.interval))

	if swrl.window[fmt.Sprintf("%d", currWindow)] != 0 {
		if swrl.window[fmt.Sprintf("%d", currWindow)] >= int32(swrl.maxreqs) {
			return false
		}
	} else {
		prevCtr := swrl.window[fmt.Sprintf("%d", prevWindow)]
		swrl.window = make(map[string]int32)
		swrl.window[fmt.Sprintf("%d", prevWindow)] = prevCtr
		swrl.window[fmt.Sprintf("%d", currWindow)] = 0
	}

	totCtr := (float32(swrl.window[fmt.Sprintf("%d", prevWindow)]) * float32(1-e)) + (float32(swrl.window[fmt.Sprintf("%d", currWindow)]) * float32(e))

	if totCtr+1 <= float32(swrl.maxreqs) {
		swrl.window[fmt.Sprintf("%d", currWindow)] += 1
		return true
	}

	return false
}
