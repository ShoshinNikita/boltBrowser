package db

import (
	"testing"
)

// Test db in testdata/test.db
//
// Structure (is such order!):
// cmd
// 	  |-> "anotherUsers"
// 						|-> "1"
//							| "age" - "99"
// 							| "name" - "Admin"
// 						|-> "2"
// 							| "name" – "Hi!!!!"
// 							| "prof" - "tester"
//						| "testData" - "15"
//    |-> "user"
// 				| "age" - "15"
// 				| "name" - "TestUser"

var cmd = []Element{Element{T: bucket, Key: "anotherUsers", Value: ""}, Element{T: bucket, Key: "user", Value: ""}}
var anotherUsers = []Element{Element{T: bucket, Key: "1", Value: ""}, Element{T: bucket, Key: "2", Value: ""}, 
Element{T: record, Key: "testData", Value: "15"}}
var bucket1 = []Element{Element{T: record, Key: "age", Value: "99"}, Element{T: record, Key: "name", Value: "Admin"}}
var bucket2 = []Element{Element{T: record, Key: "name", Value: "hi!!!!"}, Element{T: record, Key: "prof", Value: "tester"}}
var user = []Element{Element{T: record, Key: "age", Value: "15"}, Element{T: record, Key: "name", Value: "TestUser"}}


func TestSortElements(t *testing.T) {
	tests := []struct{
		slice 	[]Element
		result 	[]Element
	}{
		{
			[]Element{Element{Key: "a", T: record}, Element{Key: "b", T: bucket}},
			[]Element{Element{Key: "b", T: bucket}, Element{Key: "a", T: record}},
		},
		{
			[]Element{Element{Key: "abc", T: bucket}, Element{Key: "acd", T: bucket}},
			[]Element{Element{Key: "abc", T: bucket}, Element{Key: "acd", T: bucket}},
		},
		{
			[]Element{Element{Key: "abc", T: record}, Element{Key: "acd", T: bucket}, Element{Key: "hello", T: bucket}},
			[]Element{Element{Key: "acd", T: bucket}, Element{Key: "hello", T: bucket}, Element{Key: "abc", T: record}},
		},
		{
			[]Element{Element{Key: "abc", T: record}, Element{Key: "t", T: record}, Element{Key: "acd", T: bucket}, Element{Key: "hello", T: bucket}},
			[]Element{Element{Key: "acd", T: bucket}, Element{Key: "hello", T: bucket}, Element{Key: "abc", T: record}, Element{Key: "t", T: record}},
		},
	}


	for i, test := range tests {
		sortElements(test.slice)
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
	if testDB.Size / 1024 != 32 {
		t.Errorf("Wrong Size Want: 32 Got: %d", testDB.Size / 1024)
	}

	// Try to open wrong db
	_, err = Open("testdata/test123.db")
	if err == nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

func TestGetCMD(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}
	elements, path, err := testDB.GetCMD()
	checkCMD(t, elements, path, err)
}

func TestNext(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Bucket "anotherUsers"
	elements, path, err := testDB.Next("anotherUsers")
	checkAnotherUsers(t, elements, path, err)

	// Next bucket "1"
	elements, path, err = testDB.Next("1")
	checkBucket1(t, elements, path, err)
}

func TestBack(t *testing.T) {
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	testDB.Next("user")
	// Back to cmd
	elements, path, err := testDB.Back()
	checkCMD(t, elements, path, err)

	// Next to anotherUsers
	testDB.Next("anotherUsers")
	// Next to 1
	testDB.Next("1")
	// Back to anotherUsers
	elements, path, err = testDB.Back()
	checkAnotherUsers(t, elements, path, err)
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
	elements, path, err := testDB.GetCurrent()
	checkBucket1(t, elements, path, err)

	testDB.Back()
	testDB.Back()

	// Get cmd
	elements, path, err = testDB.GetCurrent()
	checkCMD(t, elements, path, err)
}


// Functions for testring buckets

func checkCMD(t *testing.T, elements []Element, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 0 {
		t.Errorf("Wrong currentBucket Want: [] Got: %v", path)
	}
	
	if len(elements) != 2 {
		t.Fatalf("Wrong elements Want: %v Got: %v", cmd, elements)
	}
	for i, e := range elements {
		if e != cmd[i] {
			t.Errorf("Wrong element Want: %v Got: %v", e, cmd[i])
		}
	}
}

func checkUser(t *testing.T, elements []Element, path []string, err error) {
	// Nothing
}

func checkBucket1(t *testing.T, elements []Element, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 2 || path[0] != "anotherUsers" || path[1] != "1" {
		t.Errorf("Wrong currentBucket Want: [anotherUsers 1] Got: %v", path)
	}
	
	if len(elements) != 2 {
		t.Fatalf("Wrong elements Want: %v Got: %v", bucket1, elements)
	}
	for i, e := range elements {
		if e != bucket1[i] {
			t.Errorf("Wrong element Want: %v Got: %v", e, bucket1[i])
		}
	}
}

func checkAnotherUsers(t *testing.T, elements []Element, path []string, err error) {
	if err != nil {
		t.Error(err)
	}
	if len(path) != 1 || path[0] != "anotherUsers" {
		t.Errorf("Wrong currentBucket Want: [anotherUsers] Got: %v", path)
	}
	if len(elements) != 3 {
		t.Fatalf("Wrong elements Want: %v Got: %v", anotherUsers, elements)
	}
	for i, e := range elements {
		if e != anotherUsers[i] {
			t.Errorf("Wrong element Want: %v Got: %v", e, anotherUsers[i])
		}
	}
}