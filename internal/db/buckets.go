package db

import (
	"strings"

	"github.com/boltdb/bolt"
)

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
	data.RecordsAmount = db.recordsAmount

	return data, err
}

// GetCurrent returns records from current bucket
func (db *BoltAPI) GetCurrent() (data Data, err error) {
	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})
	data.PrevBucket = true
	data.PrevRecords = (db.pages.top() > 1)
	data.NextRecords = (db.recordsAmount > maxOffset*db.pages.top())
	data.Path = "/" + strings.Join(db.currentBucket, "/")
	data.RecordsAmount = db.recordsAmount

	return data, err
}

// Back return records from previous bucket
func (db *BoltAPI) Back() (data Data, err error) {
	db.currentBucket = db.currentBucket[:len(db.currentBucket)-1]
	db.pages.del()

	if len(db.currentBucket) == 0 {
		return db.GetRoot()
	}

	err = db.db.View(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})

	data.PrevBucket = true // if there is no previous bucket will be called GetRoot()
	data.PrevRecords = (db.pages.top() > 1)
	data.NextRecords = (db.recordsAmount > maxOffset*db.pages.top())
	data.Path = "/" + strings.Join(db.currentBucket, "/")
	data.RecordsAmount = db.recordsAmount

	return data, err
}

// Next return records from next bucket with according name
func (db *BoltAPI) Next(name string) (data Data, err error) {
	db.currentBucket = append(db.currentBucket, name)
	db.pages.add()

	err = db.db.View(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		c := b.Cursor()
		data.Records = db.getRecords(c)

		return nil
	})

	data.PrevBucket = true
	data.PrevRecords = false
	data.NextRecords = (db.recordsAmount > maxOffset)
	data.Path = "/" + strings.Join(db.currentBucket, "/")
	data.RecordsAmount = db.recordsAmount
	
	return data, err
}

func (db *BoltAPI) getCurrentBucket(tx *bolt.Tx) (b *bolt.Bucket) {
	if len(db.currentBucket) == 0 {
		b = tx.Root()
	} else {
		b = tx.Bucket([]byte(db.currentBucket[0]))
		for i := 1; i < len(db.currentBucket); i++ {
			b = b.Bucket([]byte(db.currentBucket[i]))
		}
	}

	return b
}
