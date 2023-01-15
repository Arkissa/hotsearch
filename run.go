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
	go UpdateKeyword(KeywordPath, handle, k)
	go handle.Handle()

	tableName := []string{
		"bilibili",
		"douyin",
		"weibo",
		"weibu",
		"cnvd",
		"keywords",
	}

	b := pool.NewBufferPool()
	c := pool.NewClient()
	database := utils.NewDatabase(DbPath)

	targets := target.NewTargets(b, c)
	cnvd := target.NewCnvd(b, c)

	go Exit(handle, database, b, c)
	go handle.PutRun(targets, cnvd)

	noFindTable := database.FindTable(tableName...)
	for i := 0; i < len(noFindTable); i++ {
		database.CreateTable(noFindTable[i])
	}

	for {
		for website, datas := range <-handle.Data {
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
	}
}

func UpdateKeyword(n string, h *HandleData, k *utils.Keyword) {
	k.Exist = new(utils.Check)
	for {
		h.Keywords = k.Keywords(n)
		<-time.Tick(time.Second / 2)
	}
}

func Exit(c ...Close) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)
	<-s
	for i := 0; i < len(c); i++ {
		c[i].Close()
	}
	log.LogPut("Done……")
	os.Exit(0)
}
