package main

import (
	"flag"

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

	web.Initialize()
	go web.Start(flags.port, flags.debug)

	stop := make(chan bool, 0)
	<-stop
}