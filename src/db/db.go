package db

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/boltdb/bolt"

	"converters"
)

const (
	bucketTemplate = "bucket"
	recordTemplate = "record"
	maxOffset      = 100
)

// BoltAPI is a warrep for *bolt.DB
type BoltAPI struct {
	db            *bolt.DB
	currentBucket []string
	offset        int    // not in records, but in pages (1 page == maxOffset). Starts from 0
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

// Open returns info about the file of db, wrapper for *bolt.DB
func Open(path string) (*BoltAPI, error) {
	db := new(BoltAPI)
	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	db.db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

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

// GetRoot returns records from root of db
func (db *BoltAPI) GetRoot() (records []Record, bucketsPath []string, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		records = db.getFirstRecords(c)

		return nil
	})

	return records, []string{}, err
}

// GetCurrent returns records from current bucket
func (db *BoltAPI) GetCurrent() (records []Record, bucketsPath []string, err error) {
	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		records = db.getFirstRecords(c)

		return nil
	})

	return records, db.currentBucket, err
}

// Back return records from previous bucket
func (db *BoltAPI) Back() (records []Record, bucketsPath []string, err error) {
	db.currentBucket = db.currentBucket[:len(db.currentBucket)-1]
	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		records = db.getFirstRecords(c)

		return nil
	})

	return records, db.currentBucket, err
}

// Next return records from next bucket with according name
func (db *BoltAPI) Next(name string) (records []Record, bucketsPath []string, err error) {
	db.currentBucket = append(db.currentBucket, name)
	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		records = db.getFirstRecords(c)

		return nil
	})

	return records, db.currentBucket, err
}

// NextRecords return next part of records and bool, which shows is there next records
func (db *BoltAPI) NextRecords() (records []Record, canMoveNext bool, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		var c *bolt.Cursor
		if len(db.currentBucket) == 0 {
			c = tx.Cursor()
		} else {
			b := tx.Bucket([]byte(db.currentBucket[0]))
			for i := 1; i < len(db.currentBucket); i++ {
				b = b.Bucket([]byte(db.currentBucket[i]))
			}
			c = b.Cursor()
		}

		records, canMoveNext = db.getNextRecords(c)
		return nil
	})

	return records, canMoveNext, err
}

// PrevRecords return prev part of records and bool, which shows is there previous records
func (db *BoltAPI) PrevRecords() (records []Record, canMoveBack bool, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		var c *bolt.Cursor
		if len(db.currentBucket) == 0 {
			c = tx.Cursor()
		} else {
			b := tx.Bucket([]byte(db.currentBucket[0]))
			for i := 1; i < len(db.currentBucket); i++ {
				b = b.Bucket([]byte(db.currentBucket[i]))
			}
			c = b.Cursor()
		}

		records, canMoveBack = db.getPrevRecords(c)
		return nil
	})

	return records, canMoveBack, err
}

func (db *BoltAPI) getFirstRecords(c *bolt.Cursor) []Record {
	db.offset = 0
	var (
		records []Record
		counter = 0
		i       = 0
	)

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if i < maxOffset {
			var r Record
			if v == nil {
				r.T = bucketTemplate
				r.Key = converters.ConvertKey(k)
			} else {
				r.T = recordTemplate
				r.Key = converters.ConvertKey(k)
				r.Value = converters.ConvertValue(v)
			}
			records = append(records, r)
		}
		counter++
		i++
	}
	// Updating number of records
	db.recordsAmount = counter

	sortRecords(records)
	return records
}

func (db *BoltAPI) getNextRecords(c *bolt.Cursor) (records []Record, canMoveNext bool) {
	var i = 0
	// [ maxOffset * db.offset; maxOffset * (db.offset + 1) )
	for k, v := c.First(); k != nil && i < maxOffset*(db.offset+1); k, v = c.Next() {
		if maxOffset*db.offset <= i {
			var r Record
			if v == nil {
				r.T = bucketTemplate
				r.Key = converters.ConvertKey(k)
			} else {
				r.T = recordTemplate
				r.Key = converters.ConvertKey(k)
				r.Value = converters.ConvertValue(v)
			}
			records = append(records, r)
		}
		i++
	}
	db.offset++

	canMoveNext = (db.offset*maxOffset < db.recordsAmount)

	sortRecords(records)
	return records, canMoveNext
}

func (db *BoltAPI) getPrevRecords(c *bolt.Cursor) (records []Record, canMoveBack bool) {
	var i = 0
	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil && i < maxOffset*db.offset; k, v = c.Next() {
		if maxOffset*(db.offset-1) <= i {
			var r Record
			if v == nil {
				r.T = bucketTemplate
				r.Key = converters.ConvertKey(k)
			} else {
				r.T = recordTemplate
				r.Key = converters.ConvertKey(k)
				r.Value = converters.ConvertValue(v)
			}
			records = append(records, r)
		}
		i++
	}
	db.offset--

	canMoveBack = (db.offset > 0)

	sortRecords(records)
	return records, canMoveBack
}

func sortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].T == records[j].T {
			// compare keys
			return records[i].Key < records[j].Key
		}
		// compare type ("bucket" and "record")
		return records[i].T < records[j].T
	})
}
