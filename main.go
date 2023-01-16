package main

import (
	"flag"
	"log"
	"time"
)

var (
	DbPath      string
	KeywordPath string
	TimeSleep   time.Duration
	Help        bool
)

func main() {
	flag.BoolVar(&Help, "h", false, "Help")
	flag.StringVar(&DbPath, "d", "hotsearch.db", "The Database file Path")
	flag.StringVar(&KeywordPath, "k", "keywords.csv", "The keyword file Path")
	flag.DurationVar(&TimeSleep, "t", time.Hour/3, "For example, -time 30s creates a timer of 30 seconds.")
	flag.Parse()

    log.Println(TimeSleep.Seconds())
	if Help {
		flag.PrintDefaults()
		return
	}

	run()
}
