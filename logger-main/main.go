package main

import (
	"godemo/go_loger/logger"
)

var log logger.Logger

func main() {
	// log = logger.NewConsoleLogger("warring")
	log = logger.NewFileLogger("warring", "./", "my.log", 10*1024)
	for {
		log.Warring("warring %d", 100)
		log.Info("info ...")
		log.Error("error %d %v", 20, "s")
		log.Fatal("fatal ...")
		// time.Sleep(time.Second)
	}
}
