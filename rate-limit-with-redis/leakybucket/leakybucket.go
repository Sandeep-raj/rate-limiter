package leakybucket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/rate-limiter/rate-limit-with-redis/redisclient"
)

type leakyBucket struct {
	rdblock *redisclient.RdbLock
	size    int
	key     string
}

var Bucket *leakyBucket

func InitLeakyBucket(size int, key string) {
	Bucket = &leakyBucket{
		rdblock: &redisclient.RdbLock{
			Expiration: 1000,
			Key:        "leaky-bucket",
		},
		size: size,
		key:  key,
	}
}

func (bucket *leakyBucket) Allow() bool {
	uuidStr := uuid.New().String()
	bucket.rdblock.AquireLock(fmt.Sprintf("%s-%s", bucket.rdblock.Key, uuidStr))
	defer bucket.rdblock.ReleaseLock(fmt.Sprintf("%s-%s", bucket.rdblock.Key, uuidStr))

	bucketLen, err := redisclient.Rdb.LLen(context.TODO(), bucket.key).Result()

	if err != nil {
		log.Print("Error while getting length of the bucket")
		return false
	}

	if bucketLen < int64(bucket.size) {
		redisclient.Rdb.RPush(context.TODO(), bucket.key, uuidStr)
		return true
	}

	return false
}

func (bucket *leakyBucket) Consume() {
	val, err := redisclient.Rdb.LPop(context.TODO(), bucket.key).Result()
	time.Sleep(500 * time.Millisecond)
	if err != nil {
		if err.Error() != redis.Nil.Error() {
			log.Print("error getting the data")
		}
		return
	}

	log.Printf("%s data consumed", val)
}
