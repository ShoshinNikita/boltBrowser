package main

import (
	"time"
	"syscall"
	"os"
	"flag"
	"fmt"
	"os/signal"

	"web"
)


type opts struct {
	port	string
	debug	bool
}

func main() {
	var flags opts
	flag.StringVar(&flags.port, "port", ":500", "port for website (with ':')")
	flag.BoolVar(&flags.debug, "debug", false, "debug mode")
	flag.Parse()

	fmt.Printf("[INFO] Start, port - %s, debug mode - %t\n", flags.port, flags.debug)

	stopSite := make(chan struct{})
	stop := make(chan os.Signal, 1)

	web.Initialize()
	go web.Start(flags.port, flags.debug, stopSite)
	
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	close(stopSite)
	time.Sleep(100 * time.Millisecond)
	web.CloseDBs()
	fmt.Println("[INFO] Stop program")
}