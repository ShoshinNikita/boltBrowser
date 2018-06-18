package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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
	const currentVersion = "v2.1"

	params.ParseFlags()

	fmt.Printf("boltBrowser %s\n", currentVersion)
	fmt.Print("[INFO] Start. Params:\n")
	showParams()

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

	if params.OpenBrowser {
		url := "http://localhost" + params.Port
		if params.NeatWindow {
			url += "/wrapper"
		}
		err := openBrowser(url)
		if err != nil {
			fmt.Printf("[ERR] %s\n", err.Error())
		}
	}

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	close(stopSite)
	time.Sleep(100 * time.Millisecond)
	dbs.CloseDBs()
	fmt.Println("[INFO] Program was stopped")
}

func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "windows":
		{
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		}
	case "linux":
		{
			err = exec.Command("xdg-open", url).Start()
		}
	case "darwin":
		{
			err = exec.Command("open", url).Start()
		}
	default:
		{
			err = fmt.Errorf("unsupported platform")
		}
	}

	return err
}

func showParams() {
	printSpaces := func(n int) {
		for i := 0; i < n; i++ {
			fmt.Print(" ")
		}
	}

	// params should be printed under "Params:"
	const spaces = 14

	printSpaces(spaces)
	fmt.Printf("* port - %s\n", params.Port)
	printSpaces(spaces)
	fmt.Printf("* should check version - %t\n", params.CheckVer)
	printSpaces(spaces)
	fmt.Printf("* write mode - %t\n", params.IsWriteMode)
	printSpaces(spaces)
	fmt.Printf("* offset - %d\n", params.Offset)
	printSpaces(spaces)
	fmt.Printf("* should open a browser - %t\n", params.OpenBrowser)
	printSpaces(spaces)
	fmt.Printf("* should open a neat window - %t\n", params.NeatWindow)
	printSpaces(spaces)
	fmt.Printf("* debug - %t\n", params.Debug)
}
