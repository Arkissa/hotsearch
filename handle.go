package main

import (
	"fmt"
	"hotsearch/log"
	"strings"
	"time"
)

type Run interface {
	Do() map[string][]string
}

type Close interface {
	Close()
}

type HandleData struct {
	Targets  chan Run
	Data     chan map[string][]string
	Keywords []string
}

func NewHandleData() *HandleData {
	return &HandleData{
		Targets: make(chan Run, 5),
		Data:    make(chan map[string][]string, 5),
	}
}

func (h *HandleData) PutRun(run ...Run) {
	for i := 0; i < len(run); i++ {
		h.Targets <- run[i]
	}
}

func (h *HandleData) Handle() {
	log.LogPut("[INFO] Wait Target\n")
	for target := range h.Targets {
		result := make(map[string][]string)
		for t, words := range target.Do() {
			log.LogPut("[INFO] Start Handle %s\n", t)
			for w := 0; w < len(words); w++ {
				for k := 0; k < len(h.Keywords); k++ {
					if strings.Contains(words[w], h.Keywords[k]) {
						result["keywords"] = append(result["keywords"], fmt.Sprintf("%s:%s", t, words[w]))
					}
				}
				result[t] = append(result[t], words[w])
			}
		}
		h.Data <- result
	}
}

func (h *HandleData) Close() {
	for len(h.Targets) != 0 || len(h.Data) != 0 {
		time.Sleep(1e6)
	}
	close(h.Targets)
	close(h.Data)
}
