package db

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var ErrNeedWriteMode = errors.New("WriteMode is required")

// Structures for saving data of a bucket in memory

type record struct {
	k []byte
	v []byte
}

type bucket struct {
	k         []byte
	nextLevel *data
}

type data struct {
	records []record
	buckets []bucket
}

func (d *data) addRecord(k, v []byte) {
	r := record{k, v}
	d.records = append(d.records, r)
}

func (d *data) addBucket(k []byte, pointer *data) {
	b := bucket{k, pointer}
	d.buckets = append(d.buckets, b)
}

// End of structures declaration

// AddBucket adds a new bucket.
// Function returns an error if:
// * the bucket already exists - "bucket already exists"
// * there's a record with same key - "it's a record"
func (db *BoltAPI) AddBucket(bucketName string) (err error) {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	err = db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		// Bucket already exists
		if b.Bucket([]byte(bucketName)) != nil {
			return errors.New("bucket already exists")
		}

		// There's a record with same key
		if b.Get([]byte(bucketName)) != nil {
			return errors.New("it's a record")
		}

		_, err := b.CreateBucket([]byte(bucketName))
		return err
	})

	if err == nil {
		db.recordsAmount++
	}

	return err
}

// DeleteBucket deletes a bucket.
// Function returns an error if:
// * it's a record, not bucket - "it's a record"
// * there's no such bucket - "there's no such bucket"
func (db *BoltAPI) DeleteBucket(key string) (err error) {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	err = db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		// It's a record, not bucket
		if b.Get([]byte(key)) != nil {
			return errors.New("it's a record")
		}

		// There's no such bucket
		if b.Bucket([]byte(key)) == nil {
			return errors.New("there's no such bucket")
		}

		return b.DeleteBucket([]byte(key))
	})

	if err == nil {
		db.recordsAmount--
	}

	return err
}

// EditBucketName renames buckets
// Process:
// * firstly, data of a bucket is copied to memory (it is a tree);
// * secondary, bucket with old name is deleted, bucket with new name is created;
// * thirdly, the data from memory is copied to the new bucket.
// Copying  works recursively.
func (db *BoltAPI) EditBucketName(oldKey, newKey string) (err error) {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	currentData := new(data)

	return db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)
		// Check is there a bucket with newKey
		if temp := b.Bucket([]byte(newKey)); temp != nil {
			return errors.New("Bucket " + newKey + " already exists")
		}

		oldBucket := b.Bucket([]byte(oldKey))

		// Copy data from a db to memory
		copyDataToMemory(currentData, oldBucket)

		newBucket, err := b.CreateBucket([]byte(newKey))
		if err != nil {
			return err
		}

		// Delete old bucket and create new with refreshed name
		b.DeleteBucket([]byte(oldKey))

		// Copy data to the db
		copyDataToDB(currentData, newBucket)

		return nil
	})
}

func copyDataToMemory(d *data, bucket *bolt.Bucket) {
	bucket.ForEach(func(k, v []byte) error {
		if v != nil {
			// record
			d.addRecord(k, v)
		} else {
			// bucket
			// Create a new element
			ptr := new(data)
			d.addBucket(k, ptr)
			// Go to the nested bucket
			nestedBucket := bucket.Bucket(k)
			copyDataToMemory(ptr, nestedBucket)
		}
		return nil
	})
}

func copyDataToDB(d *data, bucket *bolt.Bucket) {
	// Checking just in case
	if d == nil {
		return
	}

	// Add records
	for _, r := range d.records {
		bucket.Put(r.k, r.v)
	}

	// Add buckets
	for _, b := range d.buckets {
		newBucket, _ := bucket.CreateBucket(b.k)
		copyDataToDB(b.nextLevel, newBucket)
	}
}

// AddRecord adds a new record.
// Function returns an error if:
// * there's a bucket with same key - "it's a bucket"
// * the record already exists - "record already exists"
func (db *BoltAPI) AddRecord(key, value string) (err error) {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	err = db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		// If it is a bucket
		if b.Bucket([]byte(key)) != nil {
			return errors.New("it's a bucket")
		}

		// If record exist
		if b.Get([]byte(key)) != nil {
			return errors.New("record already exists")
		}

		return b.Put([]byte(key), []byte(value))
	})

	if err == nil {
		db.recordsAmount++
	}

	return err
}

// DeleteRecord deletes a record
// Function returns an error if:
// * it's a bucket, nor record - "it's a bucket"
// * there's no such record - "there's no such record"
func (db *BoltAPI) DeleteRecord(key string) (err error) {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	err = db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		// If it is bucket
		if b.Bucket([]byte(key)) != nil {
			return errors.New("it's a bucket")
		}

		// If there's no such record
		if b.Get([]byte(key)) == nil {
			return errors.New("there's no such record")
		}

		// If record exists
		return b.Delete([]byte(key))
	})

	if err == nil {
		db.recordsAmount--
	}

	return err
}

// EditRecord edits a record
// Function returns an error if:
// * there's a bucket with key == oldKey - "it's a bucket"
// * there's no record with key == oldKey - "there's no such record"
// * there's a bucket with key == newKey - "there's a bucket with key == newKey"
// * there's a record with key == newKey - "there's a record with key == newKey"
func (db *BoltAPI) EditRecord(oldKey, newKey, newValue string) error {
	if db.ReadOnly {
		return ErrNeedWriteMode
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := db.getCurrentBucket(tx)

		// If there's a bucket with key == oldKey
		if b.Bucket([]byte(oldKey)) != nil {
			return errors.New("it's a bucket")
		}

		// If there's no record with key == oldKey
		if b.Get([]byte(oldKey)) == nil {
			return errors.New("there's no such record")
		}

		// If there's a bucket with key == newKey
		if b.Bucket([]byte(newKey)) != nil {
			return errors.New("there's a bucket with key == newKey")
		}

		// If there's a record with key == newKey
		if oldKey != newKey && b.Get([]byte(newKey)) != nil {
			return errors.New("there's a record with key == newKey")
		}

		// We can overwrite the existing record
		if oldKey == newKey {
			return b.Put([]byte(oldKey), []byte(newValue))
		}

		// We have to delete old record and create new one
		b.Delete([]byte(oldKey))
		return b.Put([]byte(newKey), []byte(newValue))
	})
}
