package redisclient

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
)

var Rdb *redis.Client
var Lock *RdbLock

type RdbLock struct {
	Expiration int
	Key        string
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		//Addr:     "host.docker.internal:6379",
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func InitRedisLock(exp int, key string) {
	if Rdb == nil {
		InitRedis()
	}

	Lock = &RdbLock{
		Expiration: exp,
		Key:        key,
	}
}

func (rdbl *RdbLock) AquireLock(val string) {
	aquire := false
	for !aquire {
		boolCmd := Rdb.SetNX(context.TODO(), rdbl.Key, val, time.Duration(rdbl.Expiration)*time.Millisecond)
		if boolCmd.Val() {
			aquire = true
		} else {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
}

func (rdbl *RdbLock) ReleaseLock(val string) {
	currVal, err := Rdb.Get(context.TODO(), rdbl.Key).Result()

	if err == nil && currVal == val {
		Rdb.Del(context.TODO(), rdbl.Key)
		log.Print("Lock Released")
	} else {
		log.Printf("Error releasing the lock err - %+v currval - %s val - %s", err, currVal, val)
	}
}
