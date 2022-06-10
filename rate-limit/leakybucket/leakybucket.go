package leakybucket

import (
	"sync"

	"github.com/golang-collections/collections/queue"
)

type leakyBucket struct {
	interval   int
	max_reqs   int
	bucket     queue.Queue
	bucketLock *sync.Mutex
}

func LeakyBucketRateLimiter(interval int, maxReqs int) leakyBucket {
	return leakyBucket{
		interval:   interval,
		max_reqs:   maxReqs,
		bucket:     queue.Queue{},
		bucketLock: &sync.Mutex{},
	}
}

func (lb *leakyBucket) Allow(req interface{}) bool {

	lb.bucketLock.Lock()
	defer lb.bucketLock.Unlock()

	if lb.bucket.Len() < lb.max_reqs {
		lb.bucket.Enqueue(req)
		return true
	}

	return false
}

func (lb *leakyBucket) Consume() interface{} {
	lb.bucketLock.Lock()
	defer lb.bucketLock.Unlock()

	if lb.bucket.Len() > 0 {
		return lb.bucket.Dequeue()
	}

	return nil
}
