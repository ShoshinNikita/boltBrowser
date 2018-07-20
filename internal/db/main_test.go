package db_test

import (
	"os"
	"testing"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
)

func TestMain(m *testing.M) {
	m.Run()

	os.Exit(0)
}
