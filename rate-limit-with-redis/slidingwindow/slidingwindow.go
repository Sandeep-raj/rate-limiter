package slidingwindow

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/rate-limiter/rate-limit-with-redis/redisclient"
)

type slidingWindow struct {
	lock     *redisclient.RdbLock
	maxreqs  int
	interval int
	key      string
}

var SlidingWindow *slidingWindow

func InitSlidingWindow() {
	SlidingWindow = &slidingWindow{
		lock: &redisclient.RdbLock{
			Expiration: 1000,
			Key:        "sliding-window",
		},
		maxreqs:  10,
		interval: 3,
		key:      "slidingwindow",
	}
}

func (slidingWindow *slidingWindow) Allow() bool {
	uuidStr := uuid.New().String()
	slidingWindow.lock.AquireLock(fmt.Sprintf("%s-%s", slidingWindow.lock.Key, uuidStr))
	defer slidingWindow.lock.ReleaseLock(fmt.Sprintf("%s-%s", slidingWindow.lock.Key, uuidStr))

	currWindow := time.Now().Unix() / int64(slidingWindow.interval)
	prevWindow := currWindow - 1

	e := float32(time.Now().Unix()%int64(slidingWindow.interval)) / float32(slidingWindow.interval)

	currCtrStr, err := redisclient.Rdb.Get(context.TODO(), fmt.Sprintf("%s-%d", slidingWindow.key, currWindow)).Result()

	if err != nil && err.Error() != redis.Nil.Error() {
		log.Print("Error in getting the current window")
		return false
	}

	prevCtrStr, err := redisclient.Rdb.Get(context.TODO(), fmt.Sprintf("%s-%d", slidingWindow.key, prevWindow)).Result()

	if err != nil && err.Error() != redis.Nil.Error() {
		log.Print("Error in getting the previous window")
		return false
	}

	currCtr, _ := strconv.Atoi(currCtrStr)
	prevCtr, _ := strconv.Atoi(prevCtrStr)

	totCtr := float32(prevCtr)*(1-e) + float32(currCtr)
	log.Printf("%d %d %f %d", currCtr, prevCtr, totCtr, time.Now().Unix())

	if totCtr+1 <= float32(slidingWindow.maxreqs) {
		redisclient.Rdb.Set(context.TODO(), fmt.Sprintf("%s-%d", slidingWindow.key, currWindow), currCtr+1, time.Duration(2*slidingWindow.interval)*time.Second)
		return true
	}

	return false
}
