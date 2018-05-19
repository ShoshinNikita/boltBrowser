package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"db"
	"web"
)

type opts struct {
	port   string
	debug  bool
	offset int
}

func main() {
	var flags opts
	flag.StringVar(&flags.port, "port", ":500", "port for website (with ':')")
	flag.BoolVar(&flags.debug, "debug", false, "debug mode")
	flag.IntVar(&flags.offset, "offset", 100, "number of records on single page")
	flag.Parse()

	fmt.Printf("[INFO] Start, port - %s, debug mode - %t, offset - %d\n", flags.port, flags.debug, flags.offset)

	stopSite := make(chan struct{})
	stop := make(chan os.Signal, 1)

	db.SetOffset(flags.offset)

	go web.Start(flags.port, flags.debug, stopSite)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	close(stopSite)
	time.Sleep(100 * time.Millisecond)
	web.CloseDBs()
	fmt.Println("[INFO] Stop program")
}
