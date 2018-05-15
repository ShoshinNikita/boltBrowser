package db

import (
	"testing"
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

var root = []Record{Record{T: bucketTemplate, Key: "anotherUsers", Value: ""}, Record{T: bucketTemplate, Key: "user", Value: ""}}
var anotherUsers = []Record{Record{T: bucketTemplate, Key: "1", Value: ""}, Record{T: bucketTemplate, Key: "2", Value: ""},
	Record{T: recordTemplate, Key: "testData", Value: "15"}}
var bucket1 = []Record{Record{T: recordTemplate, Key: "age", Value: "99"}, Record{T: recordTemplate, Key: "name", Value: "Admin"}}
var bucket2 = []Record{Record{T: recordTemplate, Key: "name", Value: "hi!!!!"}, Record{T: recordTemplate, Key: "prof", Value: "tester"}}
var user = []Record{Record{T: recordTemplate, Key: "age", Value: "15"}, Record{T: recordTemplate, Key: "name", Value: "TestUser"}}

func TestSortRecords(t *testing.T) {
	tests := []struct {
		slice  []Record
		result []Record
	}{
		{
			[]Record{Record{Key: "a", T: recordTemplate}, Record{Key: "b", T: bucketTemplate}},
			[]Record{Record{Key: "b", T: bucketTemplate}, Record{Key: "a", T: recordTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: bucketTemplate}, Record{Key: "acd", T: bucketTemplate}},
			[]Record{Record{Key: "abc", T: bucketTemplate}, Record{Key: "acd", T: bucketTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: recordTemplate}, Record{Key: "acd", T: bucketTemplate}, Record{Key: "hello", T: bucketTemplate}},
			[]Record{Record{Key: "acd", T: bucketTemplate}, Record{Key: "hello", T: bucketTemplate}, Record{Key: "abc", T: recordTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: recordTemplate}, Record{Key: "t", T: recordTemplate}, Record{Key: "acd", T: bucketTemplate}, Record{Key: "hello", T: bucketTemplate}},
			[]Record{Record{Key: "acd", T: bucketTemplate}, Record{Key: "hello", T: bucketTemplate}, Record{Key: "abc", T: recordTemplate}, Record{Key: "t", T: recordTemplate}},
		},
	}

	for i, test := range tests {
		sortRecords(test.slice)
		for j := range test.slice {
			if test.slice[j] != test.result[j] {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.result, test.slice)
			}
		}
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
	if len(testDB.currentBucket) != 0 {
		t.Errorf("Wrong currentBucket Want: 0 Got: %d", len(testDB.currentBucket))
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

func TestGetRoot(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}
	records, path, err := testDB.GetRoot()
	checkRoot(t, records, path, err)
}

func TestNext(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Bucket "anotherUsers"
	records, path, err := testDB.Next("anotherUsers")
	checkAnotherUsers(t, records, path, err)

	// Next bucket "1"
	records, path, err = testDB.Next("1")
	checkBucket1(t, records, path, err)
}

func TestBack(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	testDB.Next("user")
	// Back to root
	records, path, err := testDB.Back()
	checkRoot(t, records, path, err)

	// Next to anotherUsers
	testDB.Next("anotherUsers")
	// Next to 1
	testDB.Next("1")
	// Back to anotherUsers
	records, path, err = testDB.Back()
	checkAnotherUsers(t, records, path, err)
}

func TestGetCurrent(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	testDB.Next("user")
	testDB.Back()
	testDB.Next("anotherUsers")
	testDB.Next("1")

	// Get bucket "1"
	records, path, err := testDB.GetCurrent()
	checkBucket1(t, records, path, err)

	testDB.Back()
	testDB.Back()

	// Get root
	records, path, err = testDB.GetCurrent()
	checkRoot(t, records, path, err)
}

// Functions for testring buckets

func checkRoot(t *testing.T, records []Record, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 0 {
		t.Errorf("Wrong currentBucket Want: [] Got: %v", path)
	}

	if len(records) != 2 {
		t.Fatalf("Wrong records Want: %v Got: %v", root, records)
	}
	for i, e := range records {
		if e != root[i] {
			t.Errorf("Wrong record Want: %v Got: %v", e, root[i])
		}
	}
}

func checkUser(t *testing.T, records []Record, path []string, err error) {
	// Nothing
}

func checkBucket1(t *testing.T, records []Record, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 2 || path[0] != "anotherUsers" || path[1] != "1" {
		t.Errorf("Wrong currentBucket Want: [anotherUsers 1] Got: %v", path)
	}

	if len(records) != 2 {
		t.Fatalf("Wrong records Want: %v Got: %v", bucket1, records)
	}
	for i, e := range records {
		if e != bucket1[i] {
			t.Errorf("Wrong record Want: %v Got: %v", e, bucket1[i])
		}
	}
}

func checkAnotherUsers(t *testing.T, records []Record, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 1 || path[0] != "anotherUsers" {
		t.Errorf("Wrong currentBucket Want: [anotherUsers] Got: %v", path)
	}
	if len(records) != 3 {
		t.Fatalf("Wrong records Want: %v Got: %v", anotherUsers, records)
	}
	for i, e := range records {
		if e != anotherUsers[i] {
			t.Errorf("Wrong record Want: %v Got: %v", e, anotherUsers[i])
		}
	}
}
