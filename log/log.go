package log

import (
	"log"
	"os"
)

func LogOutErr(msg string, errs error) {
	logFile, err := os.OpenFile("hotsearch.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("log file open err: ", err)
		return
	}
	log.SetOutput(logFile)
	log.Printf("[ERROR] %s: %s\n", msg, errs.Error())
}

func LogPut(format string, msg ...any) {
	log.SetOutput(os.Stdout)
	log.Printf(format, msg...)
}
