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
	handle := NewHandleData()
	go UpdateKeywords(handle, k)
	go handle.Handle()

	tableName := []string{
		"bilibili",
		"douyin",
		"weibo",
		"weibu",
		"cnvd",
		"keywords",
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)

	b := pool.NewBufferPool()
	c := pool.NewClient()
	database := utils.NewDatabase(DbPath)

	targets := target.NewTargets(b, c)
	cnvd := target.NewCnvd(b, c)

	noFindTable := database.FindTable(tableName...)
	for i := 0; i < len(noFindTable); i++ {
		database.CreateTable(noFindTable[i])
	}

	handle.PutRun(targets, cnvd)
	for {
		select {
		case <-s:
			Exit(handle, database, b, c)
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
		case <-time.Tick(TimeSleep):
			handle.PutRun(targets, cnvd)
		}
	}
}

func UpdateKeywords(h *HandleData, k *utils.Keyword) {
	k.Exist = new(utils.Check)
	for {
		log.LogPut("get Keywords")
		h.Keywords = k.Keywords(KeywordPath)
		<-time.Tick(time.Second / 2)
	}
}

func Exit(c ...Close) {
	for i := 0; i < len(c); i++ {
		c[i].Close()
	}
	log.LogPut("Done……")
	os.Exit(0)
}
