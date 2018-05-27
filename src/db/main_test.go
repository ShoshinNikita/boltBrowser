package db_test

import (
	"testing"

	"params"
)

func TestMain(m *testing.M) {
	// Turn WriteMode on
	params.IsWriteMode = true

	m.Run()
}
