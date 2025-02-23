package main

import (
	"github.com/go-co-op/gocron"
	kline2 "jobcenter/kline"
	"time"
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	kline := kline2.NewKline()
	s.Every(1).Minute().Do(func() {
		kline.Do("1m")
	})
	s.Every(1).Hour().Do(func() {
		kline.Do("1H")
	})
	s.StartBlocking()
}
