package main

import (
	"flag"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func main() {
	debugPtr := flag.Bool("debug", false, "show debug log")
	jsonPtr := flag.Bool("json", false, "show log in json format")
	flag.Parse()

	if *jsonPtr {
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("Define Log as Json format")
		os.Exit(1)
	}

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
		log.Info("Define Log Debug level")
		runtime.Goexit()
	}

	items := []int{1, 2, 3, 4, 5}

	for _, i := range items {
		log.Info("This is Info log", i)
		if i == 2 {
			os.Exit(1)
		}
	}

	// log.Info("This is Info log")
	// log.Debug("This is Debug log")
	// log.Warning("This is Warning log")
	// log.Error("This is Error log")

}
