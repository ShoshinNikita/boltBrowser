package db

import (
	"errors"

	"github.com/boltdb/bolt"
)

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
