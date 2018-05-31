package db_test

import (
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

func TestNextRecords(t *testing.T) {
	SetOffset(100)

	tests := []struct {
		offset      int
		next        []string
		nextCounter int
		answer      []Record
	}{
		{1, []string{"user"}, 1, []Record{rcrd("name", "TestUser")}},
		{1, []string{"anotherUsers"}, 1, []Record{bckt("2")}},
		{1, []string{"anotherUsers", "1"}, 1, []Record{rcrd("name", "Admin")}},
		{2, []string{"anotherUsers"}, 1, []Record{rcrd("testData", "15")}},
	}

	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		SetOffset(test.offset)

		for _, n := range test.next {
			testDB.Next(n)
		}

		var data Data
		var err error
		for i := 0; i < test.nextCounter; i++ {
			data, err = testDB.NextRecords()
		}
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(data.Records, test.answer) {
			t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, test.answer, data.Records)
		}

		testDB.ClearPath()
	}
}

func TestPrevRecords(t *testing.T) {
	SetOffset(100)

	tests := []struct {
		offset      int
		next        []string
		nextCounter int
		backCounter int
		answer      []Record
	}{
		{1, []string{"user"}, 1, 1, []Record{rcrd("age", "15")}},
		{1, []string{"anotherUsers"}, 2, 1, []Record{bckt("2")}},
		{1, []string{"anotherUsers", "1"}, 1, 1, []Record{rcrd("age", "99")}},
		{3, []string{"anotherUsers"}, 1, 1, []Record{bckt("1"), bckt("2"), rcrd("testData", "15")}},
	}

	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		SetOffset(test.offset)

		for _, n := range test.next {
			testDB.Next(n)
		}

		for i := 0; i < test.nextCounter; i++ {
			testDB.NextRecords()
		}

		var data Data
		var err error
		for i := 0; i < test.backCounter; i++ {
			data, err = testDB.PrevRecords()
		}
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(data.Records, test.answer) {
			t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, test.answer, data.Records)
		}

		testDB.ClearPath()
	}
}

func TestSortRecords(t *testing.T) {
	SetOffset(100)

	tests := []struct {
		slice  []Record
		result []Record
	}{
		{
			[]Record{Record{Key: "a", T: RecordTemplate}, Record{Key: "b", T: BucketTemplate}},
			[]Record{Record{Key: "b", T: BucketTemplate}, Record{Key: "a", T: RecordTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: BucketTemplate}, Record{Key: "acd", T: BucketTemplate}},
			[]Record{Record{Key: "abc", T: BucketTemplate}, Record{Key: "acd", T: BucketTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: RecordTemplate}, Record{Key: "acd", T: BucketTemplate}, Record{Key: "hello", T: BucketTemplate}},
			[]Record{Record{Key: "acd", T: BucketTemplate}, Record{Key: "hello", T: BucketTemplate}, Record{Key: "abc", T: RecordTemplate}},
		},
		{
			[]Record{Record{Key: "abc", T: RecordTemplate}, Record{Key: "t", T: RecordTemplate}, Record{Key: "acd", T: BucketTemplate}, Record{Key: "hello", T: BucketTemplate}},
			[]Record{Record{Key: "acd", T: BucketTemplate}, Record{Key: "hello", T: BucketTemplate}, Record{Key: "abc", T: RecordTemplate}, Record{Key: "t", T: RecordTemplate}},
		},
	}

	for i, test := range tests {
		SortRecords(test.slice)
		for j := range test.slice {
			if test.slice[j] != test.result[j] {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.result, test.slice)
			}
		}
	}
}
