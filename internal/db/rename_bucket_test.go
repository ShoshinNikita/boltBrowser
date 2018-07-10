package db_test

import (
	"testing"

	db "github.com/ShoshinNikita/boltBrowser/internal/db"
)

// Test db in testdata/rename.db
//
// Structure:
// root
// 	  |-> "hello"
//				|-> firstBucket
//							| "help" - "132"
//				|-> secondBucket
//							| "789" - "987"
//							| "world" - "555"
//				| "test" - "123"
//				| "test2" - "1234"
//
// The db has same structure after all tests

var (
	helloBucket = []db.Record{
		bckt("firstBucket"),
		bckt("secondBucket"),
		rcrd("test", "123"),
		rcrd("test2", "1234"),
	}
	firstBucket = []db.Record{
		rcrd("help", "132"),
	}
	secondBucket = []db.Record{
		rcrd("789", "987"),
		rcrd("world", "555"),
	}
)

func TestBucketRenaming(t *testing.T) {
	db.SetOffset(100)
	const (
		oldName = "hello"
		newName = "123"
	)

	testDB, err := db.Open("testdata/rename.db")
	if err != nil {
		t.Fatal(err)
	}

	// hello -> 123
	t.Log("hello -> 123")
	err = testDB.EditBucketName(oldName, newName)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	testDB.ClearPath()

	// Check
	data, err := testDB.Next(newName)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if !equal(data.Records, helloBucket) {
		t.Errorf("Want: %v Got: %v", helloBucket, data.Records)
	}
	// firstBucket
	data, err = testDB.Next("firstBucket")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if !equal(data.Records, firstBucket) {
		t.Errorf("Want: %v Got: %v", firstBucket, data.Records)
	}
	testDB.Back()
	// secondBucket
	data, err = testDB.Next("secondBucket")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if !equal(data.Records, secondBucket) {
		t.Errorf("Want: %v Got: %v", secondBucket, data.Records)
	}
	testDB.ClearPath()

	// 123 -> hello
	t.Log("123 -> hello")
	err = testDB.EditBucketName(newName, oldName)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}

	// Check
	data, err = testDB.Next(oldName)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if len(data.Records) != len(helloBucket) {
		t.Errorf("Want: %v Got: %v", helloBucket, data.Records)
		return
	}
	for i := range data.Records {
		if data.Records[i] != helloBucket[i] {
			t.Errorf("Want: %v Got: %v", helloBucket, data.Records)
			return
		}
	}
	data, err = testDB.Next("firstBucket")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if len(data.Records) != len(firstBucket) {
		t.Errorf("Want: %v Got: %v", firstBucket, data.Records)
		return
	}
	for i := range data.Records {
		if data.Records[i] != firstBucket[i] {
			t.Errorf("Want: %v Got: %v", firstBucket, data.Records)
			return
		}
	}
	testDB.Back()

	data, err = testDB.Next("secondBucket")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}
	if len(data.Records) != len(secondBucket) {
		t.Errorf("Want: %v Got: %v", secondBucket, data.Records)
		return
	}
	for i := range data.Records {
		if data.Records[i] != secondBucket[i] {
			t.Errorf("Want: %v Got: %v", secondBucket, data.Records)
			return
		}
	}

	err = testDB.Close()
	if err != nil {
		t.Error(err)
	}
}
