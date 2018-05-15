package db

import (
	"sort"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"

	"converters"
)

const (
	bucket = "bucket"
	record = "record"
	maxOffset = 100
)

// DBApi is a warrep for *bolt.DB
type DBApi struct {
	db            *bolt.DB
	currentBucket []string
	offset        int // not in records, but in pages (1 page == maxOffset). Starts from 0
	recordsAmount int // number of records in current bucket
	Name          string `json:"name"`
	FilePath      string `json:"filePath"`
	Size          int64  `json:"size"`
}

// Element consists information about record in the db
type Element struct {
	T     string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Open return info about the file of db, wrapper for *bolt.DB
func Open(path string) (*DBApi, error) {
	db := new(DBApi)
	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	db.db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Getting info about the file
	db.FilePath = path
	db.Name = filepath.Base(path)
	file, _ := os.Stat(path)
	db.Size = file.Size()

	return db, nil
}

// Close closes db
func (db *DBApi) Close() error {
	return db.db.Close()
}

// GetCMD returns records from root of db
func (db *DBApi) GetCMD() ([]Element, []string, error) {
	var elements []Element
	err := db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		elements = db.getFirstElements(c)

		return nil
	})

	return elements, []string{}, err
}

// GetCurrent returns records from current bucket
func (db *DBApi) GetCurrent() ([]Element, []string, error) {
	var elements []Element
	if len(db.currentBucket) == 0 {
		return db.GetCMD()
	}

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		elements = db.getFirstElements(c)

		return nil
	})

	return elements, db.currentBucket, err
}

// Back return records from previous bucket
func (db *DBApi) Back() ([]Element, []string, error) {
	var elements []Element
	db.currentBucket = db.currentBucket[:len(db.currentBucket)-1]
	if len(db.currentBucket) == 0 {
		return db.GetCMD()
	}

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		elements = db.getFirstElements(c)

		return nil
	})

	return elements, db.currentBucket, err
}

// Next return records from next bucket with according name
func (db *DBApi) Next(name string) ([]Element, []string, error) {
	var elements []Element

	db.currentBucket = append(db.currentBucket, name)
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		elements = db.getFirstElements(c)

		return nil
	})

	return elements, db.currentBucket, err
}

// NextRecords return next part of records and bool, which shows is there next elements
func (db *DBApi) NextRecords() (elements []Element, canMoveNext bool, err error) {
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

		elements, canMoveNext = db.getNextElements(c)
		return nil
	})

	return elements, canMoveNext, err
}

// PrevRecords return prev part of records and bool, which shows is there previous elements
func (db *DBApi) PrevRecords() (elements []Element, canMoveBack bool, err error) {
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

		elements, canMoveBack = db.getPrevElements(c)
		return nil
	})

	return elements, canMoveBack, err
}


func (db *DBApi) getFirstElements(c *bolt.Cursor) []Element {
	db.offset = 0
	var (
		elements []Element
		counter = 0
		i = 0
	) 
	
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if i < maxOffset {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = converters.ConvertKey(k)
			} else {
				elem.T = record
				elem.Key = converters.ConvertKey(k)
				elem.Value= converters.ConvertValue(v)
			}
			elements = append(elements, elem)
		}
		counter++
		i++
	}
	// Updating number of records
	db.recordsAmount = counter

	sortElements(elements)
	return elements
}

func (db *DBApi) getNextElements(c *bolt.Cursor) (elements []Element, canMoveNext bool) {
	var i = 0
	// [ maxOffset * db.offset; maxOffset * (db.offset + 1) )
	for k, v := c.First(); k != nil && i < maxOffset * (db.offset + 1); k, v = c.Next() {
		if maxOffset * db.offset <= i {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = converters.ConvertKey(k)
			} else {
				elem.T = record
				elem.Key = converters.ConvertKey(k)
				elem.Value= converters.ConvertValue(v)
			}
			elements = append(elements, elem)
		}
		i++
	}
	db.offset++
	
	canMoveNext = (db.offset * maxOffset < db.recordsAmount)

	sortElements(elements)
	return elements, canMoveNext
}

func (db *DBApi) getPrevElements(c *bolt.Cursor) (elements []Element, canMoveBack bool) {
	var i = 0
	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil && i < maxOffset * db.offset; k, v = c.Next() {
		if  maxOffset * (db.offset - 1) <= i {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = converters.ConvertKey(k)
			} else {
				elem.T = record
				elem.Key = converters.ConvertKey(k)
				elem.Value= converters.ConvertValue(v)
			}
			elements = append(elements, elem)
		}
		i++
	}
	db.offset--

	canMoveBack = (db.offset > 0)

	sortElements(elements)
	return elements, canMoveBack
}

func sortElements(elements []Element) {
	sort.Slice(elements, func (i, j int) bool {
		if elements[i].T == elements[j].T {
			// compare keys
			return elements[i].Key < elements[j].Key
		}
		// compare type ("bucket" and "record")
		return elements[i].T < elements[j].T
	})
}