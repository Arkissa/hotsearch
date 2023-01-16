package main

import (
	"hotsearch/log"
	"hotsearch/pool"
	"hotsearch/target"
	"hotsearch/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func run() {
	k := new(utils.Keyword)
	k.Exist = new(utils.Check)
	k.Tick = time.NewTicker(time.Second / 2)
	handle := NewHandleData()
    handle.Ticker = time.NewTicker(TimeSleep)
	go handle.Handle()

	tableName := []string{
		"bilibili",
		"douyin",
		"weibo",
		"weibu",
		"cnvd",
		"freebuf",
		"keywords",
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)

	b := pool.NewBufferPool()
	c := pool.NewClient()
	database := utils.NewDatabase(DbPath)

	targets := target.NewTargets(b, c)
	cnvd := target.NewCnvd(b, c)
    
	go handle.PutRun(targets, cnvd)

	noFindTable := database.FindTable(tableName...)
	for i := 0; i < len(noFindTable); i++ {
		database.CreateTable(noFindTable[i])
	}

	for {
		select {
		case <-s:
			die(handle, database, b, c)
            k.Tick.Stop()
            handle.Ticker.Stop()
		case d, ok := <-handle.Data:
			if !ok {
				return
			}
			for website, datas := range d {
				log.LogPut("[INFO] Get Data %s %d\n", website, len(datas))
				for num := 0; num < len(datas); num++ {
					noFindTable = database.FindTable(website)
					for i := 0; i < len(noFindTable); i++ {
						if !database.CreateTable(noFindTable[i]) {
							goto jump
						}
					}

					if !database.InsertData(website, datas[num]) {
						log.LogPut("[WARNING] Insert %s %s Error", website, datas[num])
						goto jump
					}
				}
			jump:
				database.Deduplication(website)
			}
		case <-k.Tick.C:
			handle.Keywords = k.Keywords(KeywordPath)
		}
	}
}

func die(c ...Close) {
	for i := 0; i < len(c); i++ {
		c[i].Close()
	}
	log.LogPut("Done……")
	os.Exit(0)
}
