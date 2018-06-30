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

	"github.com/ShoshinNikita/boltBrowser/internal/db"
	"github.com/ShoshinNikita/boltBrowser/internal/dbs"
	"github.com/ShoshinNikita/boltBrowser/internal/flags"
	"github.com/ShoshinNikita/boltBrowser/internal/versioning"
	"github.com/ShoshinNikita/boltBrowser/internal/web"
)

const currentVersion = "v2.2"

func main() {
	flags.ParseFlags()

	fmt.Printf("boltBrowser %s\n", currentVersion)
	fmt.Print("[INFO] Start. flags:\n")
	showFlags()

	if flags.CheckVer {
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

	db.SetOffset(flags.Offset)

	go web.Start(flags.Port, stopSite)

	if flags.OpenBrowser {
		url := "http://localhost" + flags.Port
		if flags.NeatWindow {
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
	dbs.CloseDBs()

	// Wait just in case
	time.Sleep(100 * time.Millisecond)
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

func showFlags() {
	printSpaces := func(n int) {
		for i := 0; i < n; i++ {
			fmt.Print(" ")
		}
	}

	// flags should be printed under "flags:"
	const spaces = 14

	printSpaces(spaces)
	fmt.Printf("* port - %s\n", flags.Port)
	printSpaces(spaces)
	fmt.Printf("* should check version - %t\n", flags.CheckVer)
	printSpaces(spaces)
	fmt.Printf("* write mode - %t\n", flags.IsWriteMode)
	printSpaces(spaces)
	fmt.Printf("* offset - %d\n", flags.Offset)
	printSpaces(spaces)
	fmt.Printf("* should open a browser - %t\n", flags.OpenBrowser)
	printSpaces(spaces)
	fmt.Printf("* should open a neat window - %t\n", flags.NeatWindow)
	printSpaces(spaces)
	fmt.Printf("* debug - %t\n", flags.Debug)
}
