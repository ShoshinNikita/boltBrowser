package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"db"
	"dbs"
	"versioning"
	"web"
)

type opts struct {
	port     string
	debug    bool
	offset   int
	checkVer bool
}

func main() {
	const currentVersion = "v1.3"

	var flags opts
	flag.StringVar(&flags.port, "port", ":500", "port for website (with ':')")
	flag.BoolVar(&flags.debug, "debug", false, "debug mode")
	flag.IntVar(&flags.offset, "offset", 100, "number of records on single page")
	flag.BoolVar(&flags.checkVer, "checkVer", true, "should program check is there a new version")
	flag.Parse()

	// Checking of ':' before port
	if flags.port[0] != ':' {
		flags.port = ":" + flags.port
	}

	fmt.Printf("boltBrowser %s\n", currentVersion)
	fmt.Printf("[INFO] Start, port - %s, debug mode - %t, offset - %d, check version - %t\n", flags.port, flags.debug, flags.offset, flags.checkVer)

	if flags.checkVer {
		// Checking is there a new version
		data, err := versioning.CheckVersion(currentVersion)
		if err != nil {
			fmt.Printf("[ERR] Can't check is there a new version: %s", err.Error())
		} else if data.IsNewVersion {
			changes := "+ " + strings.Join(data.Changes, "\n+ ")
			fmt.Printf("\n[INFO] New version (%s) is available.\nChanges:\n%s\nLink: %s\n\n", data.LastVersion, changes, data.Link)
		} else {
			fmt.Printf("[INFO] You use the last version of boltBrowser\n")
		}
	}

	// Init of channels
	stopSite := make(chan struct{})
	stop := make(chan os.Signal, 1)

	db.SetOffset(flags.offset)

	go web.Start(flags.port, flags.debug, stopSite)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	close(stopSite)
	time.Sleep(100 * time.Millisecond)
	dbs.CloseDBs()
	fmt.Println("[INFO] Program was stopped")
}
