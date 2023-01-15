package main

import (
	"flag"
)

var (
	DbPath      string
	KeywordPath string
	Help        bool
)

func main() {
	flag.StringVar(&DbPath, "d", "hotsearch.db", "The Database file Path")
	flag.StringVar(&KeywordPath, "k", "keywords.csv", "The keyword file Path")
	flag.BoolVar(&Help, "h", false, "Help")
	flag.Parse()

	if Help {
		flag.PrintDefaults()
		return
	}

	run()
}
