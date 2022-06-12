package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rate-limiter/rate-limit-with-redis/fixedwindow"
	"github.com/rate-limiter/rate-limit-with-redis/leakybucket"
	"github.com/rate-limiter/rate-limit-with-redis/redisclient"
	"github.com/rate-limiter/rate-limit-with-redis/slidingwindow"
	"github.com/rate-limiter/rate-limit-with-redis/tokenbucket"
)

func init() {
	redisclient.InitRedis()
	//fixedwindow.InitFixedWindow()
	// leakybucket.InitLeakyBucket(10, "lbucket")
	// slidingwindow.InitSlidingWindow()
	tokenbucket.InitTokenBucket()
}

func main() {
	log.Print("Hello World")
	// redisclient.InitRedis()

	// // setres := redisclient.Rdb.Set(context.TODO(), "golang", 34, 0)
	// // log.Print(setres.String())

	// res, err := redisclient.Rdb.Get(context.TODO(), "golag").Result()
	// if err.Error() == redis.Nil.Error() {
	// 	log.Print("xxx")
	// }
	// log.Printf("res - %s, err - %+v", res, err)
	//log.Print(redisclient.Rdb.LLen(context.TODO(), "test").Result())
	// Producer()
	sig := make(chan int)
	//go LeakyBucketProducer()
	//go LeakyBucketConsumer()
	// go SlidingWindowProducer()
	go TokenBucketProducer()
	<-sig
}

func FixedWindowProducer() {
	for i := 0; i < 100; i++ {
		if fixedwindow.FixedWindow.Allow() {
			log.Printf("%s allowed %d", os.Getenv("hostname"), i)
		} else {
			log.Printf("%s not allowed %d", os.Getenv("hostname"), i)
		}

		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	}
}

func LeakyBucketProducer() {
	for i := 0; i < 100; i++ {
		if leakybucket.Bucket.Allow() {
			log.Printf("%s allowed %d", os.Getenv("hostname"), i)
		} else {
			log.Printf("%s not allowed %d", os.Getenv("hostname"), i)
		}

		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	}
}

func LeakyBucketConsumer() {
	for {
		leakybucket.Bucket.Consume()
	}
}

func SlidingWindowProducer() {
	for i := 0; i < 50; i++ {
		if slidingwindow.SlidingWindow.Allow() {
			log.Printf("%s allowed %d", os.Getenv("hostname"), i)
		} else {
			log.Printf("%s not allowed %d", os.Getenv("hostname"), i)
		}

		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	}
}

func TokenBucketProducer() {
	for i := 0; i < 50; i++ {
		if tokenbucket.Bucket.Allow() {
			log.Printf("%s allowed %d", os.Getenv("hostname"), i)
		} else {
			log.Printf("%s not allowed %d", os.Getenv("hostname"), i)
		}

		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	}
}
