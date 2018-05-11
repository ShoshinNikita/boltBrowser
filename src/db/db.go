package db

import (
	"github.com/boltdb/bolt"
)

const (
	bucket = "bucket"
	record = "record"
)

type DBApi struct {
	db            *bolt.DB
	currentBucket []string
}

type element struct {
	T     string
	Key   string
	Value string
}

func Open(path string) (*DBApi, error) {
	db := new(DBApi)
	var err error
	db.db, err = bolt.Open(path, 0644, nil)
	return db, err
}

func (db *DBApi) Close() error {
	return db.db.Close()
}

func (db *DBApi) GetCMD() ([]element, error) {
	var elements []element

	err := db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem element
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

func (db *DBApi) Back() ([]element, error) {
	var elements []element

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
			var elem element
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

func (db *DBApi) Next(name string) ([]element, error) {
	var elements []element

	db.currentBucket = append(db.currentBucket, name)
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var elem element
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