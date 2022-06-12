package fixedwindow

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/rate-limiter/rate-limit-with-redis/redisclient"
)

type fixedWindow struct {
	rdblock   *redisclient.RdbLock
	size      int
	interval  int
	windowKey string
	ctrKey    string
}

var FixedWindow *fixedWindow

func InitFixedWindow() {
	FixedWindow = &fixedWindow{
		rdblock: &redisclient.RdbLock{
			Expiration: 1000,
			Key:        "fixed-window",
		},
		size:      10,
		interval:  3,
		windowKey: "windowKey",
		ctrKey:    "ctrKey",
	}

}

func (fw *fixedWindow) Allow() bool {
	uuidStr := uuid.New()
	fw.rdblock.AquireLock(fmt.Sprintf("%s-%s", fw.rdblock.Key, uuidStr))
	defer fw.rdblock.ReleaseLock(fmt.Sprintf("%s-%s", fw.rdblock.Key, uuidStr))

	currTime := fmt.Sprintf("%d", time.Now().Unix()/int64(fw.interval))
	windowStr, err := redisclient.Rdb.Get(context.TODO(), fw.windowKey).Result()

	if err != nil && err.Error() != redis.Nil.Error() {
		log.Print("Error while getting window string")
		return false
	}

	if currTime != windowStr && strings.Compare(currTime, windowStr) == 1 {
		_, err = redisclient.Rdb.Set(context.TODO(), fw.windowKey, currTime, 0).Result()
		if err != nil && err.Error() != redis.Nil.Error() {
			log.Print("Error while setting windowKey")
			return false
		}
		redisclient.Rdb.Set(context.TODO(), fw.ctrKey, fw.size, 0)
		if err != nil && err.Error() != redis.Nil.Error() {
			log.Print("Error while setting ctrKey")
			return false
		}
	}

	ctrStr, err := redisclient.Rdb.Get(context.TODO(), fw.ctrKey).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		log.Print("Error while getting counter")
		return false
	}

	if ctrStr != "0" {
		_, err = redisclient.Rdb.Decr(context.TODO(), fw.ctrKey).Result()
		if err != nil && err.Error() != redis.Nil.Error() {
			log.Print("Error while decr ctrKey")
			return false
		}
		return true
	}

	return false
}
