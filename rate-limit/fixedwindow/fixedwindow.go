package fixedwindow

import (
	"sync"
	"time"
)

type fixedWindow struct {
	interval    int
	maxrequests int
	currTime    int64
	ctr         int
	windowLock  *sync.Mutex
}

func FixedWindowRateLimiter(interval int, max_reqs int) fixedWindow {
	currTime := time.Now().Unix()
	return fixedWindow{
		interval:    interval,
		maxrequests: max_reqs,
		currTime:    currTime,
		ctr:         max_reqs,
		windowLock:  &sync.Mutex{},
	}
}

func (fw *fixedWindow) Allow() bool {
	fw.windowLock.Lock()
	defer fw.windowLock.Unlock()

	currTime := time.Now().Unix()
	if currTime < fw.currTime+int64(fw.interval) {
		if fw.ctr > 0 {
			fw.ctr--
			return true
		}
	} else {
		fw.currTime = currTime
		fw.ctr = fw.maxrequests

		fw.ctr--
		return true
	}

	return false
}
