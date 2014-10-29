package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/matryer/goblueprints/chapter8/backup"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		interval = flag.Int("interval", 10, "interval between checks (seconds)")
		archive  = flag.String("archive", "archive", "path to archive location")
		pathlist = flag.String("paths", "", "colon separated list of paths to backup")
	)
	flag.Parse()
	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.DefaultArchiver,
		Paths:       make(map[string]string),
	}
	paths := strings.Split(*pathlist, ":")
	if len(paths) < 1 {
		log.Fatalln("must provide at least one path")
	}
	for _, path := range paths {
		m.Paths[path] = ""
	}
	check(m)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	stop := false
	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m)
		case <-signalChan:
			// stop
			fmt.Println()
			log.Printf("Stopping...")
			stop = true
			break
		}
		if stop {
			break
		}
	}
}

func check(m *backup.Monitor) {
	log.Println("Checking...")
	counter, err := m.Now()
	if err != nil {
		log.Fatalln("failed to backup:", err)
	}
	if counter > 0 {
		log.Printf("  Archived %d directories\n", counter)
	} else {
		log.Println("  No changes")
	}
}
