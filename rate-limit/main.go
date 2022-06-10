package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/rate-limit/fixedwindow"
	"github.com/rate-limit/leakybucket"
	"github.com/rate-limit/slidingwindow"
	"github.com/rate-limit/tokenbucket"
)

func main() {
	// TestFixedWindowRateLimiter()
	// TestLeakyBucketRateLimiter()
	TestSlidingWindowRateLimiter()
	// TestTokenBucketRateLimiter()
}

func TestFixedWindowRateLimiter() {
	fwrl := fixedwindow.FixedWindowRateLimiter(3, 10)
	for i := 0; i < 100; i++ {
		ok := fwrl.Allow()
		if ok {
			log.Printf("%d instance allowed", i)
		} else {
			log.Printf("%d instance not allowed", i)
		}
		sleept := rand.Intn(5)
		time.Sleep(time.Duration(sleept*100) * time.Millisecond)
	}
}

func TestLeakyBucketRateLimiter() {
	lbrl := leakybucket.LeakyBucketRateLimiter(10, 10)
	hold := make(chan string)

	processCtr := 0

	go (func() {
		for {
			data := lbrl.Consume()
			if data == nil {
				continue
			}
			processCtr++
			log.Print(data, processCtr)
			time.Sleep(500 * time.Millisecond)
		}
	})()

	for i := 0; i < 100; i++ {
		ok := lbrl.Allow("process task")
		if ok {
			log.Printf("%d instance allowed", i)
		} else {
			log.Printf("%d instance not allowed", i)
		}
		sleept := rand.Intn(5)
		time.Sleep(time.Duration(sleept*100) * time.Millisecond)
	}

	<-hold
	log.Print("finished")
}

func TestSlidingWindowRateLimiter() {
	swrl := slidingwindow.SlidingWindowRateLimiter(10, 3)

	for i := 0; i < 100; i++ {
		ok := swrl.Allow()
		if ok {
			log.Printf("%d instance allowed", i)
		} else {
			log.Printf("%d instance not allowed", i)
		}
		sleept := rand.Intn(5)
		time.Sleep(time.Duration(sleept*20) * time.Millisecond)
	}
}

func TestTokenBucketRateLimiter() {
	tbrl := tokenbucket.TokenBucketRateLimiter(2, 7)

	for i := 0; i < 100; i++ {
		ok := tbrl.Allow()
		if ok {
			log.Printf("%d instance allowed", i)
		} else {
			log.Printf("%d instance not allowed", i)
		}
		sleept := rand.Intn(5)
		time.Sleep(time.Duration(sleept*20) * time.Millisecond)
	}
}
