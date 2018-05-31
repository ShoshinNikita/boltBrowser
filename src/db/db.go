package db

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"

	"params"
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

	var options *bolt.Options
	// Check is ReadOnly mode
	if !params.IsWriteMode {
		options = &bolt.Options{ReadOnly: true}
	} else {
		options = nil
	}

	db.db, err = bolt.Open(path, 0600, options)
	if err != nil {
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

// Close closes db
func (db *BoltAPI) Close() error {
	return db.db.Close()
}
