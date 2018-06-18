package db_test

import (
	"os"
	"testing"

	"github.com/ShoshinNikita/boltBrowser/internal/flags"
)

func TestMain(m *testing.M) {
	// Turn WriteMode on
	flags.IsWriteMode = true

	m.Run()

	os.Exit(0)
}
