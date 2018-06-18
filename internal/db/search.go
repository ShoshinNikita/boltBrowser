package db

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/boltdb/bolt"
)

// Search use bytes.Contains()
func (db *BoltAPI) Search(needle string) (records []Record, path string, recordsAmount int, err error) {
	bNeedle := []byte(needle)

	err = db.db.View(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)
		c := b.Cursor()

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
	recordsAmount = len(records)

	sortRecords(records)
	return records, path, recordsAmount, err
}

// SearchRegexp use regexp.Match()
func (db *BoltAPI) SearchRegexp(expr string) (records []Record, path string, recordsAmount int, err error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return []Record{}, "", 0, err
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)
		c := b.Cursor()

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
	recordsAmount = len(records)

	sortRecords(records)
	return records, path, recordsAmount, err
}
