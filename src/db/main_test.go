package db_test

import (
	"os"
	"testing"

	"flags"
)

func TestMain(m *testing.M) {
	// Turn WriteMode on
	flags.IsWriteMode = true

	m.Run()

	os.Exit(0)
}
