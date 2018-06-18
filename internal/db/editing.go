package db

import (
	"errors"

	"github.com/boltdb/bolt"
)

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

func (db *BoltAPI) AddBucket(bucketName string) (err error) {
	err = db.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := db.getCurrentBucket(tx)
		if b.Bucket([]byte(bucketName)) != nil {
			// Bucket already exists
			err = errors.New("bucket already exist")
		} else if b.Get([]byte(bucketName)) != nil {
			// There's a record with same key
			err = errors.New("\"" + bucketName + "\" is a record")
		} else {
			_, err = b.CreateBucket([]byte(bucketName))
		}

		return err
	})

	if err == nil {
		db.recordsAmount++
	}

	return err
}

func (db *BoltAPI) DeleteBucket(key string) (err error) {
	err = db.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := db.getCurrentBucket(tx)

		if b.Bucket([]byte(key)) != nil {
			// Bucket can be deleted
			err = b.DeleteBucket([]byte(key))
		} else if b.Get([]byte(key)) != nil {
			// It is a record
			err = errors.New("\"" + key + "\" is a record")
		}

		return err
	})

	if err == nil {
		db.recordsAmount--
	}

	return err
}

// EditBucketName renames buckets
// Proccess:
// * firstly, data of a bucket is copied to memory (it is a tree);
// * secondary, bucket with old name is deleted, bucket with new name is created;
// * thirdly, the data from memory is copied to the new bucket.
// Copying  works recursively.
func (db *BoltAPI) EditBucketName(oldKey, newKey string) (err error) {
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

func (db *BoltAPI) AddRecord(key, value string) (err error) {
	err = db.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := db.getCurrentBucket(tx)
		if b.Bucket([]byte(key)) != nil {
			// If it is bucket
			err = errors.New("\"" + key + "\" is a bucket")
		} else if b.Get([]byte(key)) == nil {
			// If record doesn't exist
			err = b.Put([]byte(key), []byte(value))
		} else {
			// If record exist
			err = errors.New("\"" + key + "\" already exists")
		}

		return err
	})

	if err == nil {
		db.recordsAmount++
	}

	return err
}

func (db *BoltAPI) DeleteRecord(key string) (err error) {
	err = db.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := db.getCurrentBucket(tx)
		if b.Get([]byte(key)) != nil {
			// If record exists
			b.Delete([]byte(key))
		} else if b.Bucket([]byte(key)) != nil {
			// If it is bucket
			err = errors.New("\"" + key + "\" is a bucket")
		} else {
			// If record doesn't exist
			err = errors.New("\"" + key + "\" doesn't exist")
		}

		return err
	})

	if err == nil {
		db.recordsAmount--
	}

	return err
}

func (db *BoltAPI) EditRecord(oldKey, newKey, newValue string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		var err error
		b := db.getCurrentBucket(tx)

		if b.Bucket([]byte(oldKey)) != nil {
			// If it is bucket
			err = errors.New("\"" + oldKey + "\" is a bucket")
		} else if b.Get([]byte(oldKey)) != nil {
			if oldKey == newKey {
				err = b.Put([]byte(oldKey), []byte(newValue))
			} else {
				if b.Get([]byte(newKey)) != nil {
					// If newKey already exists
					err = errors.New("\"" + newKey + "\" already exists")
				} else {
					b.Delete([]byte(oldKey))
					err = b.Put([]byte(newKey), []byte(newValue))
				}
			}
		} else {
			// If record doesn't exist
			err = errors.New("\"" + oldKey + "\" doesn't exist")
		}

		return err
	})
}
