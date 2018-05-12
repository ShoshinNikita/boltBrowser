package main

import (
	"web"
)

func main() {
	web.Initialize()
	go web.Start()
	stop := make(chan bool, 0)
	<-stop
}