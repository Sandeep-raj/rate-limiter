package tokenbucket

import (
	"sync"
	"time"
)

type tokenbucket struct {
	interval        int
	max_reqs        int
	rate            float32
	ctr             float32
	last_reset_time int64
	bucketlock      *sync.Mutex
}

func TokenBucketRateLimiter(interval int, maxReqs int) tokenbucket {
	currTime := time.Now().Unix()

	return tokenbucket{
		interval:        interval,
		max_reqs:        maxReqs,
		rate:            float32(maxReqs) / float32(interval),
		ctr:             float32(maxReqs),
		last_reset_time: currTime,
		bucketlock:      &sync.Mutex{},
	}
}

func (tbrl *tokenbucket) Allow() bool {
	tbrl.bucketlock.Lock()
	defer tbrl.bucketlock.Unlock()

	currTime := time.Now().Unix()

	if currTime == tbrl.last_reset_time && tbrl.ctr < 1 {
		return false
	}

	if currTime != tbrl.last_reset_time {
		diff := currTime - tbrl.last_reset_time

		tbrl.ctr += float32(diff) * tbrl.rate
		if tbrl.ctr > float32(tbrl.max_reqs) {
			tbrl.ctr = float32(tbrl.max_reqs)
		}

		tbrl.last_reset_time = currTime
	}

	tbrl.ctr = tbrl.ctr - 1
	return true
}
