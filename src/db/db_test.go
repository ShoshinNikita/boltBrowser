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

func TestSortRecords(t *testing.T) {
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

func TestOpen(t *testing.T) {
	// Try to open correct db
	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
	}
	// Check opened db
	if len(testDB.GetCurrentBucket()) != 0 {
		t.Errorf("Wrong currentBucket Want: 0 Got: %d", len(testDB.GetCurrentBucket()))
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
	tests := []struct {
		offset int
		answer []Record
	}{
		{100, []Record{Record{T: BucketTemplate, Key: "anotherUsers", Value: ""}, Record{T: BucketTemplate, Key: "user", Value: ""}}},
		{2, []Record{Record{T: BucketTemplate, Key: "anotherUsers", Value: ""}, Record{T: BucketTemplate, Key: "user", Value: ""}}},
		{1, []Record{Record{T: BucketTemplate, Key: "anotherUsers", Value: ""}}},
	}

	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		SetOffset(test.offset)
		data, err := testDB.GetRoot()
		if err != nil {
			t.Error(err)
			continue
		}
		if !equal(test.answer, data.Records) {
			t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, test.answer, data.Records)
		}
	}
}

func TestNext(t *testing.T) {
	type T struct {
		next   string
		answer []Record
	}
	tests := []struct {
		offset int
		data   []T
	}{
		{100, []T{
			T{"anotherUsers", []Record{bckt("1"), bckt("2"), rcrd("testData", "15")}},
			T{"1", []Record{rcrd("age", "99"), rcrd("name", "Admin")}}}},
		{1, []T{
			T{"anotherUsers", []Record{bckt("1")}},
			T{"1", []Record{rcrd("age", "99")}}}},
		{2, []T{
			T{"anotherUsers", []Record{bckt("1"), bckt("2")}},
			T{"2", []Record{rcrd("name", "hi!!!!"), rcrd("prof", "tester")}}}},
		{1, []T{
			T{"user", []Record{rcrd("age", "15")}}}},
	}

	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		SetOffset(test.offset)
		for _, d := range test.data {
			data, err := testDB.Next(d.next)
			if err != nil {
				t.Error(err)
				break
			}

			if !equal(d.answer, data.Records) {
				t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, d.answer, data.Records)
				break
			}
		}

		testDB.ClearPath()
	}

}

func TestBack(t *testing.T) {
	tests := []struct {
		offset int
		next   []string
		answer [][]Record
	}{
		{100, []string{"anotherUsers", "1"}, [][]Record{
			[]Record{bckt("1"), bckt("2"), rcrd("testData", "15")},
			[]Record{bckt("anotherUsers"), bckt("user")}}},
		{1, []string{"anotherUsers", "1"}, [][]Record{
			[]Record{bckt("1")},
			[]Record{bckt("anotherUsers")}}},
		{2, []string{"anotherUsers", "2"}, [][]Record{
			[]Record{bckt("1"), bckt("2")},
			[]Record{bckt("anotherUsers"), bckt("user")}}},
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

		for _, d := range test.answer {
			data, err := testDB.Back()
			if err != nil {
				t.Error(err)
				continue
			}
			if !equal(d, data.Records) {
				t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, d, data.Records)
				break
			}
		}

		testDB.ClearPath()
	}

}

func TestGetCurrent(t *testing.T) {
	tests := []struct {
		offset int
		next   []string
		answer []Record
	}{
		{100, []string{"user"}, []Record{rcrd("age", "15"), rcrd("name", "TestUser")}},
		{1, []string{"user"}, []Record{rcrd("age", "15")}},
		{100, []string{"anotherUsers"}, []Record{bckt("1"), bckt("2"), rcrd("testData", "15")}},
		{100, []string{"anotherUsers", "1"}, []Record{rcrd("age", "99"), rcrd("name", "Admin")}},
		{1, []string{"anotherUsers", "1"}, []Record{rcrd("age", "99")}},
		{2, []string{"anotherUsers", "1"}, []Record{rcrd("age", "99"), rcrd("name", "Admin")}},
	}

	testDB, err := Open("testdata/test.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		SetOffset(test.offset)
		for _, s := range test.next {
			testDB.Next(s)
		}

		data, err := testDB.GetCurrent()
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

func TestNextRecords(t *testing.T) {
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