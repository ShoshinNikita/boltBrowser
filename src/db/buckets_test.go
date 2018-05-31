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

func TestGetRoot(t *testing.T) {
	SetOffset(100)

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
	SetOffset(100)

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
	SetOffset(100)

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
	SetOffset(100)

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
