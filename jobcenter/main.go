package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Seconds().Do(func() {
		fmt.Println(time.Now().String())
	})
	s.StartBlocking()
}
