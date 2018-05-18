package db

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/boltdb/bolt"
)

// Search use bytes.Contains()
func (db *BoltAPI) Search(needle string) (records []Record, path string, err error) {
	bNeedle := []byte(needle)

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

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Contains(k, bNeedle) {
				var r Record
				if v == nil {
					r = Record{T: bucketTemplate, Key: string(k), Value: ""}
				} else {
					r = Record{T: recordTemplate, Key: string(k), Value: string(v)}
				}

				records = append(records, r)
			}
		}

		return nil
	})

	path = "/" + strings.Join(db.currentBucket, "/")
	return records, path, err
}

// SearchRegexp use regexp.Match()
func (db *BoltAPI) SearchRegexp(expr string) (records []Record, path string, err error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return []Record{}, "", err
	}

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

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if reg.Match(k) {
				var r Record
				if v == nil {
					r = Record{T: bucketTemplate, Key: string(k), Value: ""}
				} else {
					r = Record{T: recordTemplate, Key: string(k), Value: string(v)}
				}

				records = append(records, r)
			}
		}

		return nil
	})

	path = "/" + strings.Join(db.currentBucket, "/")
	return records, path, err
}
