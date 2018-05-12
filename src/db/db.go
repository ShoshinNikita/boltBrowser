package db

import (
	"path/filepath"
	"os"

	"github.com/boltdb/bolt"
)

const (
	bucket = "bucket"
	record = "record"
)

// DBApi is a warrep for *bolt.DB
type DBApi struct {
	db            	*bolt.DB
	currentBucket 	[]string
	Name 			string	`json:"name"`
	Path 			string `json:"path"`
	Size 			int64  `json:"size"`
}

// TODO func sort for Elements

// Element consists information about record in the db
type Element struct {
	T	  string `json:"type"`
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
	db.Path = path
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
// Also return bool, which shows can user move back (true – bucket, false – root of db)
func (db *DBApi) GetCMD() ([]Element, bool, error) {
	var elements []Element

	err := db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = string(k)
			} else {
				elem.T = record
				elem.Key = string(k)
				elem.Value = string(v)
			}
			elements = append(elements, elem)
		}

		return nil
	})

	return elements, false, err
}

// GetCurrent returns records from current bucket
// Also return bool, which shows can user move back (true – bucket, false – root of db)
func (db *DBApi) GetCurrent() ([]Element, bool, error) {
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
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = string(k)
			} else {
				elem.T = record
				elem.Key = string(k)
				elem.Value = string(v)
			}
			elements = append(elements, elem)
		}
		return nil
	})

	return elements, true, err
}

// Back return records from previous bucket
// Also return bool, which shows can user move back (true – bucket, false – root of db)
func (db *DBApi) Back() ([]Element, bool, error) {
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
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = string(k)
			} else {
				elem.T = record
				elem.Key = string(k)
				elem.Value = string(v)
			}
			elements = append(elements, elem)
		}

		return nil
	})

	return elements, true, err
}

// Next return records from next bucket with according name
func (db *DBApi) Next(name string) ([]Element, error) {
	var elements []Element

	db.currentBucket = append(db.currentBucket, name)
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem Element
			if v == nil {
				elem.T = bucket
				elem.Key = string(k)
			} else {
				elem.T = record
				elem.Key = string(k)
				elem.Value = string(v)
			}
			elements = append(elements, elem)
		}

		return nil
	})

	return elements, err
}