package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"db"
	"dbs"
	"params"
	"versioning"
	"web"
)

func main() {
	const currentVersion = "v2.0"

	params.ParseFlags()

	fmt.Printf("boltBrowser %s\n", currentVersion)
	fmt.Printf("[INFO] Start, port - %s, debug mode - %t, offset - %d, check version - %t, read-only: %t\n", params.Port, params.Debug, params.Offset, params.CheckVer, !params.IsWriteMode)

	if params.CheckVer {
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

	db.SetOffset(params.Offset)

	go web.Start(params.Port, stopSite)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	close(stopSite)
	time.Sleep(100 * time.Millisecond)
	dbs.CloseDBs()
	fmt.Println("[INFO] Program was stopped")
}
