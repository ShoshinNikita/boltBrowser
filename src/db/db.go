package db

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/boltdb/bolt"

	"converters"
)

const (
	bucketTemplate = "bucket"
	recordTemplate = "record"
)

var maxOffset int

// BoltAPI is a warrep for *bolt.DB
//
// offset is a number of pages (1 page = maxOffset). offset points at end of page, for example,
// 1 – [0, 100)
// 2 – [100, 200)
// 3 – [200, 300)
// etc.
type BoltAPI struct {
	db            *bolt.DB
	currentBucket []string
	offset        offsetStack // not in records, but in pages (1 page == maxOffset). Page points on the n * maxOffset
	recordsAmount int         // number of records in current bucket
	Name          string      `json:"name"`
	DBPath        string      `json:"dbPath"`
	Size          int64       `json:"size"`
}

// Record consists information about record in the db
type Record struct {
	T     string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Data serves for returning
type Data struct {
	Records     []Record
	PrevBucket  bool
	PrevRecords bool
	NextRecords bool
	Path        string
}

// SetOffset change value of maxOffset
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

	db.db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// For root
	db.offset.add()

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
func (db *BoltAPI) GetRoot() (data Data, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})
	data.PrevBucket = false
	data.PrevRecords = false
	data.NextRecords = (db.recordsAmount > maxOffset)
	data.Path = "/"

	return data, err
}

// GetCurrent returns records from current bucket
func (db *BoltAPI) GetCurrent() (data Data, err error) {
	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})
	data.PrevBucket = true
	data.PrevRecords = (db.offset.top() > 1)
	data.NextRecords = (db.recordsAmount > maxOffset*db.offset.top())
	data.Path = "/" + strings.Join(db.currentBucket, "/")

	return data, err
}

// Back return records from previous bucket
func (db *BoltAPI) Back() (data Data, err error) {
	db.currentBucket = db.currentBucket[:len(db.currentBucket)-1]
	db.offset.del()

	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})

	data.PrevBucket = true // if there is no previous bucket will be called GetRoot()
	data.PrevRecords = (db.offset.top() > 1)
	data.NextRecords = (db.recordsAmount > maxOffset*db.offset.top())
	data.Path = "/" + strings.Join(db.currentBucket, "/")

	return data, err
}

// Next return records from next bucket with according name
func (db *BoltAPI) Next(name string) (data Data, err error) {
	db.currentBucket = append(db.currentBucket, name)
	db.offset.add()

	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})

	data.PrevBucket = true
	data.PrevRecords = false
	data.NextRecords = (db.recordsAmount > maxOffset)
	data.Path = "/" + strings.Join(db.currentBucket, "/")

	return data, err
}

// NextRecords return next part of records and bool, which shows is there next records
func (db *BoltAPI) NextRecords() (data Data, err error) {
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

		data.Records, data.NextRecords = db.getNextRecords(c)
		return nil
	})
	data.PrevBucket = (len(db.currentBucket) != 0)
	data.PrevRecords = true

	return data, err
}

// PrevRecords return prev part of records and bool, which shows is there previous records
func (db *BoltAPI) PrevRecords() (data Data, err error) {
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

		data.Records, data.PrevRecords = db.getPrevRecords(c)
		return nil
	})
	data.PrevBucket = (len(db.currentBucket) != 0)
	data.NextRecords = true

	return data, err
}

// return records with current offset (db.offset.top())
// also update db.recordsAmount
func (db *BoltAPI) getRecords(c *bolt.Cursor) (records []Record) {
	var (
		i       int
		counter int
	)

	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if maxOffset*(db.offset.top()-1) <= i && i < maxOffset*db.offset.top() {
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
		counter++
	}

	// Updating number of records
	db.recordsAmount = counter

	sortRecords(records)
	return records
}

func (db *BoltAPI) getNextRecords(c *bolt.Cursor) (records []Record, canMoveNext bool) {
	var i = 0
	// [ maxOffset * db.offset; maxOffset * (db.offset + 1) )
	for k, v := c.First(); k != nil && i < maxOffset*(db.offset.top()+1); k, v = c.Next() {
		if maxOffset*db.offset.top() <= i {
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
	db.offset.inc()

	canMoveNext = (db.offset.top()*maxOffset < db.recordsAmount)

	sortRecords(records)
	return records, canMoveNext
}

func (db *BoltAPI) getPrevRecords(c *bolt.Cursor) (records []Record, canMoveBack bool) {
	db.offset.dec()

	var i = 0
	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil && i < maxOffset*db.offset.top(); k, v = c.Next() {
		if maxOffset*(db.offset.top()-1) <= i {
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

	canMoveBack = (db.offset.top() > 1)

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
