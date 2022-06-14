package tokenbucket

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rate-limiter/rate-limit-with-redis/redisclient"
)

type tokenBucket struct {
	lock     *redisclient.RdbLock
	interval int
	maxreqs  int
	key      string
	ctr      string
}

var Bucket *tokenBucket

func InitTokenBucket() {
	Bucket = &tokenBucket{
		lock: &redisclient.RdbLock{
			Expiration: 1000,
			Key:        "token-bucket",
		},
		interval: 2,
		maxreqs:  10,
		key:      "tbucket",
		ctr:      "tbucket-ctr",
	}
}

// func (bucket *tokenBucket) Allow() bool {
// 	uuidStr := uuid.New().String()
// 	bucket.lock.AquireLock(fmt.Sprintf("%s-%s", bucket.lock.Key, uuidStr))
// 	defer bucket.lock.ReleaseLock(fmt.Sprintf("%s-%s", bucket.lock.Key, uuidStr))

// 	currWindow := time.Now().Unix()

// 	windowStr, err := redisclient.Rdb.Get(context.TODO(), bucket.key).Result()

// 	if err != nil {
// 		if err.Error() == redis.Nil.Error() {
// 			windowStr = fmt.Sprintf("%d", currWindow)
// 			_, err = redisclient.Rdb.Set(context.TODO(), bucket.key, windowStr, 0).Result()
// 			if err != nil {
// 				log.Print("error while setting the key")
// 				return false
// 			}
// 			_, err = redisclient.Rdb.Set(context.TODO(), bucket.ctr, bucket.maxreqs, 0).Result()
// 			if err != nil {
// 				log.Print("error while setting the ctr")
// 				return false
// 			}
// 		} else {
// 			log.Print("error while getting the bucket key")
// 			return false
// 		}
// 	}

// 	window, err := strconv.ParseInt(windowStr, 10, 64)
// 	if err != nil {
// 		log.Print("Error parsing the window string")
// 		return false
// 	}

// 	rate := float64(bucket.maxreqs) / float64(bucket.interval)
// 	currCtrStr, err := redisclient.Rdb.Get(context.TODO(), bucket.ctr).Result()

// 	if err != nil {
// 		if err.Error() == redis.Nil.Error() {
// 			_, err = redisclient.Rdb.Set(context.TODO(), bucket.ctr, bucket.maxreqs, 0).Result()
// 			if err != nil {
// 				log.Print("error while setting the ctr")
// 				return false
// 			}
// 		} else {
// 			log.Print("error while setting the ctr")
// 			return false
// 		}
// 	}

// 	currCtr, err := strconv.ParseFloat(currCtrStr, 32)
// 	if err != nil {
// 		log.Print("Error parsing the ctr string")
// 		return false
// 	}

// 	totCtr := currCtr + ((float64(currWindow) - float64(window)) * rate)
// 	if totCtr >= 1 {
// 		if totCtr >= float64(bucket.maxreqs) {
// 			totCtr = float64(bucket.maxreqs)
// 		}

// 		_, err = redisclient.Rdb.Set(context.TODO(), bucket.ctr, totCtr-1, 0).Result()
// 		if err != nil {
// 			log.Print("error setting the counter value")
// 			return false
// 		}

// 		if currWindow > window {
// 			_, err := redisclient.Rdb.Set(context.TODO(), bucket.key, currWindow, 0).Result()
// 			if err != nil {
// 				log.Print("error while setting the window key")
// 				return false
// 			}
// 		}

// 		return true
// 	}

// 	return false
// }

func (bucket *tokenBucket) Allow() bool {
	uuidStr := uuid.New().String()
	bucket.lock.AquireLock(fmt.Sprintf("%s-%s", bucket.lock.Key, uuidStr))
	defer bucket.lock.ReleaseLock(fmt.Sprintf("%s-%s", bucket.lock.Key, uuidStr))

	currWindow := time.Now().Unix()

	redisclient.Rdb.SetNX(context.TODO(), bucket.key, fmt.Sprintf("%d", currWindow), 0)
	redisclient.Rdb.SetNX(context.TODO(), bucket.ctr, bucket.maxreqs, 0)

	windowStr, _ := redisclient.Rdb.Get(context.TODO(), bucket.key).Result()
	currCtrStr, _ := redisclient.Rdb.Get(context.TODO(), bucket.ctr).Result()
	window, err := strconv.ParseInt(windowStr, 10, 64)
	if err != nil {
		log.Print("Error parsing the window string")
		return false
	}

	currCtr, err := strconv.ParseFloat(currCtrStr, 32)
	if err != nil {
		log.Print("Error parsing the ctr string")
		return false
	}

	rate := float64(bucket.maxreqs) / float64(bucket.interval)

	totCtr := currCtr + ((float64(currWindow) - float64(window)) * rate)

	if totCtr >= 1 {
		if totCtr >= float64(bucket.maxreqs) {
			totCtr = float64(bucket.maxreqs)
		}

		_, err = redisclient.Rdb.Set(context.TODO(), bucket.ctr, totCtr-1, 0).Result()
		if err != nil {
			log.Print("error setting the counter value")
			return false
		}

		if currWindow > window {
			_, err := redisclient.Rdb.Set(context.TODO(), bucket.key, currWindow, 0).Result()
			if err != nil {
				log.Print("error while setting the window key")
				return false
			}
		}

		return true
	}

	return false

}
