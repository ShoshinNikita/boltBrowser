package db_test

import (
	"errors"
	"testing"

	. "db"
)

// Test db in testdata/test.db
//
// Structure (is such order!):
// root
// 	  |-> "anotherUsers"
// 						|-> "1"
//							| "age" - "99"
// 							| "name" - "Admin"
// 						|-> "2"
// 							| "name" – "Hi!!!!"
// 							| "prof" - "tester"
//						| "testData" - "15"
//	  |-> "user"
// 				| "age" - "15"
// 				| "name" - "TestUser"

// check are slices equal
func equal(want, got []Record) bool {
	if len(want) != len(got) {
		return false
	}

	for i := range want {
		if want[i] != got[i] {
			return false
		}
	}

	return true
}

// return Record{T: "bucket"}
func bckt(key string) Record {
	return Record{T: BucketTemplate, Key: key, Value: ""}
}

// return Record{T: "record"}
func rcrd(key, value string) Record {
	return Record{T: RecordTemplate, Key: key, Value: value}
}

func newErr(err string) error {
	if err == "" {
		return nil
	} else {
		return errors.New(err)
	}
}

func TestOpen(t *testing.T) {
	// Try to open correct db
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
	}
	// Check opened db
	if len(testDB.GetCurrentBucketsPath()) != 0 {
		t.Errorf("Wrong currentBucket Want: 0 Got: %d", len(testDB.GetCurrentBucketsPath()))
	}
	if testDB.Name != "test.db" {
		t.Errorf("Wrong Name Want: test.db Got: %s", testDB.Name)
	}
	// in Kv
	if testDB.Size/1024 != 32 {
		t.Errorf("Wrong Size Want: 32 Got: %d", testDB.Size/1024)
	}

	// Try to open wrong db
	_, err = Open("testdata/test123.db")
	if err == nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}
