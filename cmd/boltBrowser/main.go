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

	"github.com/ShoshinNikita/boltBrowser/internal/config"
	"github.com/ShoshinNikita/boltBrowser/internal/db"
	"github.com/ShoshinNikita/boltBrowser/internal/dbs"
	"github.com/ShoshinNikita/boltBrowser/internal/versioning"
	"github.com/ShoshinNikita/boltBrowser/internal/web"
)

const currentVersion = "v2.2"

func main() {
	err := config.ParseConfig()
	if err != nil {
		fmt.Printf("[ERR] Couldn't parse config: %s\n", err.Error())
		os.Exit(2)
	}

	fmt.Printf("boltBrowser %s\n", currentVersion)
	fmt.Print("[INFO] Start. flags:\n")
	showFlags()

	if config.Opts.CheckVer {
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

	db.SetOffset(config.Opts.Offset)

	go web.Start(config.Opts.Port, stopSite)

	if config.Opts.OpenBrowser {
		url := "http://localhost" + config.Opts.Port
		if config.Opts.NeatWindow {
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
	fmt.Printf("* port - %s\n", config.Opts.Port)
	printSpaces(spaces)
	fmt.Printf("* should check version - %t\n", config.Opts.CheckVer)
	printSpaces(spaces)
	fmt.Printf("* write mode - %t\n", config.Opts.IsWriteMode)
	printSpaces(spaces)
	fmt.Printf("* offset - %d\n", config.Opts.Offset)
	printSpaces(spaces)
	fmt.Printf("* should open a browser - %t\n", config.Opts.OpenBrowser)
	printSpaces(spaces)
	fmt.Printf("* should open a neat window - %t\n", config.Opts.NeatWindow)
	printSpaces(spaces)
	fmt.Printf("* debug - %t\n", config.Opts.Debug)
}
