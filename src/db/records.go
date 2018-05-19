package db

import (
	"sort"

	"github.com/boltdb/bolt"

	"converters"
)

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

// return records from current page (db.page.top())
// also updates db.recordsAmount
func (db *BoltAPI) getRecords(c *bolt.Cursor) (records []Record) {
	var (
		i       int
		counter int
	)

	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if maxOffset*(db.pages.top()-1) <= i && i < maxOffset*db.pages.top() {
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
	for k, v := c.First(); k != nil && i < maxOffset*(db.pages.top()+1); k, v = c.Next() {
		if maxOffset*db.pages.top() <= i {
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
	db.pages.inc()

	canMoveNext = (db.pages.top()*maxOffset < db.recordsAmount)

	sortRecords(records)
	return records, canMoveNext
}

func (db *BoltAPI) getPrevRecords(c *bolt.Cursor) (records []Record, canMoveBack bool) {
	db.pages.dec()

	var i = 0
	// [ maxOffset * (db.offset - 1); maxOffset * db.offset )
	for k, v := c.First(); k != nil && i < maxOffset*db.pages.top(); k, v = c.Next() {
		if maxOffset*(db.pages.top()-1) <= i {
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

	canMoveBack = (db.pages.top() > 1)

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
