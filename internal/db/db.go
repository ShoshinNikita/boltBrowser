package db

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
)

const (
	bucketTemplate = "bucket"
	recordTemplate = "record"
)

var maxOffset = 100

// BoltAPI is a warrep for *bolt.DB
//
// pages is a number of pages (1 page = maxOffset)
// 1 – [0, maxOffset)
// 2 – [maxOffset, 2*maxOffset)
// 3 – [2*maxOffset, 3*maxOffset)
// etc.
type BoltAPI struct {
	db            *bolt.DB
	currentBucket []string
	pages         pagesStack
	recordsAmount int    // number of records in current bucket
	Name          string `json:"name"`
	DBPath        string `json:"dbPath"`
	Size          int64  `json:"size"`
}

// Record consists information about record in the db
type Record struct {
	T     string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Data serves for returning
type Data struct {
	Records       []Record
	PrevBucket    bool
	PrevRecords   bool
	NextRecords   bool
	RecordsAmount int
	Path          string
}

// SetOffset change value of maxOffset (default – 100)
func SetOffset(offset int) {
	maxOffset = offset
}

// Open returns info about the file of db, wrapper for *bolt.DB
func Open(path string) (*BoltAPI, error) {
	db := new(BoltAPI)
	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	options := &bolt.Options{Timeout: time.Second}

	// Check is ReadOnly mode
	if !config.Opts.IsWriteMode {
		options.ReadOnly = true
	}

	db.db, err = bolt.Open(path, 0600, options)
	if err != nil {
		if err == bolt.ErrTimeout {
			return nil, fmt.Errorf("The database is already opened")
		}

		return nil, err
	}

	// For root
	db.pages.add()

	// Getting info about the file
	db.DBPath = path
	db.Name = filepath.Base(path)
	file, _ := os.Stat(path)
	db.Size = file.Size()

	return db, nil
}

// Create a new db. If path consists only a name, the db will be created on the Desktop
// The function calls Open() and returns result of this calling
func Create(path string) (*BoltAPI, error) {
	// Add ".db" if path hasn't it
	if !strings.HasSuffix(path, ".db") {
		path += ".db"
	}

	nameRegex := regexp.MustCompile(`^[\w_-]*\.db$`)
	if nameRegex.Match([]byte(path)) {
		// Path consists only a name, so we have to add the path to the Desktop
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		if runtime.GOOS == "windows" {
			path = home + "\\Desktop\\" + path
		} else {
			path = home + "/Desktop/" + path
		}
	}

	// Normalizing of path (replacing '\\' and '\')
	reg := regexp.MustCompile(`\\\\|\\`)
	path = reg.ReplaceAllString(path, "/")

	// Check if a db already exists
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("Database with path\"%s\" already exists", path)
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Close()

	return Open(path)
}

// Close closes db
func (db *BoltAPI) Close() error {
	return db.db.Close()
}
