package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
	"github.com/ShoshinNikita/boltBrowser/internal/db"
	"github.com/ShoshinNikita/boltBrowser/internal/dbs"
	"github.com/ShoshinNikita/boltBrowser/internal/web"
)

const currentVersion = "v2.4"

func main() {
	err := config.ParseConfig()
	if err != nil {
		log.Errorf("Couldn't parse config: %s\n", err)
		os.Exit(2)
	}

	log.Printf("boltBrowser %s\n", currentVersion)
	log.Infoln("Start. flags:")
	showFlags()

	// Init of channels
	stopSite := make(chan struct{})
	stop := make(chan os.Signal, 1)

	db.SetOffset(config.Opts.Offset)

	go web.Start(config.Opts.Port, stopSite)

	if config.Opts.OpenBrowser {
		url := fmt.Sprintf("http://localhost:%d", config.Opts.Port)

		err := openBrowser(url)
		if err != nil {
			log.Errorf("%s\n", err.Error())
		}
	}

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	close(stopSite)
	dbs.CloseDBs()

	// Wait just in case
	time.Sleep(100 * time.Millisecond)
	log.Infoln("Program was stopped")
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
	fmt.Printf("* port - %d\n", config.Opts.Port)
	printSpaces(spaces)
	fmt.Printf("* offset - %d\n", config.Opts.Offset)
	printSpaces(spaces)
	fmt.Printf("* should open a browser - %t\n", config.Opts.OpenBrowser)
}
