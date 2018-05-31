package db_test

import (
	"os"
	"testing"

	"params"
)

func TestMain(m *testing.M) {
	// Turn WriteMode on
	params.IsWriteMode = true

	m.Run()

	os.Exit(0)
}
