package utils

import (
	"encoding/csv"
	"hotsearch/log"
	"os"
	"strings"
)

type Status interface {
	Check(name, suffix string) bool
}

type Check struct{}

type Keyword struct {
	Exist Status
}

func (c *Keyword) Keywords(name string) (keywords []string) {
	if c.Exist.Check(name, "csv") {
		return keywords
	} else {
		file, _ := os.Open(name)
		defer file.Close()

		reader := csv.NewReader(file)
		keywords, err := reader.Read()
		if err != nil {
			log.LogOutErr("Reader Keyword err", err)
			return keywords
		}

		return keywords
	}
}

func (is *Check) Check(name, suffix string) bool {
	fileName := strings.Split(name, ".")
	if fileName[len(fileName)-1] != suffix {
		return false
	}

	_, err := os.Stat(name)
	return os.IsNotExist(err)
}
