package db_test

import (
	"testing"

	. "db"
)

// benchmark
//   |-> B1ucket
//   |-> firstBucket
//   |-> secondBucket
//   |-> thirdBucket
// 	 | 1
// 	 | 151
// 	 | 1561
// 	 | 2
// 	 | 3
// 	 | 31
// 	 | 351
// 	 | 4
// 	 | 51
// 	 | 61
// 	 | 648
// 	 | 651
// 	 | 74
// 	 | 8
// 	 | 84
// 	 | 94
// 	 | 984

func TestSearch(t *testing.T) {
	SetOffset(100)

	tests := []struct {
		request string
		answer  []Record
	}{
		{"1", []Record{
			bckt("B1ucket"),
			rcrd("1", "hello"),
			rcrd("151", "hello"),
			rcrd("1561", "hello"),
			rcrd("31", "hello"),
			rcrd("351", "hello"),
			rcrd("51", "hello"),
			rcrd("61", "hello"),
			rcrd("651", "hello")}},
		{"12", []Record{}},
		{"51", []Record{
			rcrd("151", "hello"),
			rcrd("351", "hello"),
			rcrd("51", "hello"),
			rcrd("651", "hello")}},
		{"Bucket", []Record{
			bckt("firstBucket"),
			bckt("secondBucket"),
			bckt("thirdBucket")}},
		{"cket", []Record{
			bckt("B1ucket"),
			bckt("firstBucket"),
			bckt("secondBucket"),
			bckt("thirdBucket")}},
	}

	testDB, err := Open("testdata/search.db")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
	}

	testDB.Next("benchmark")
	for i, test := range tests {
		result, path, _, err := testDB.Search(test.request)
		if err != nil {
			t.Error(err)
			continue
		}

		if path != "/benchmark" {
			t.Errorf("Test #%d Bad path: %s", i, path)
		}
		if len(result) != len(test.answer) {
			t.Errorf("Test #%d Bad size Want: %v Got %v", i, test.answer, result)
			continue
		}

		if !equal(test.answer, result) {
			t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, test.answer, result)
			break
		}
	}
}

func TestSearchRegex(t *testing.T) {
	SetOffset(100)

	tests := []struct {
		request string
		err     string
		answer  []Record
	}{
		{"^1", "", []Record{
			rcrd("1", "hello"),
			rcrd("151", "hello"),
			rcrd("1561", "hello")}},
		{"12", "", []Record{}},
		{"51$", "", []Record{
			rcrd("151", "hello"),
			rcrd("351", "hello"),
			rcrd("51", "hello"),
			rcrd("651", "hello")}},
		{"^[seconthird]+Bucket", "", []Record{
			bckt("secondBucket"),
			bckt("thirdBucket")}},
		{"(?<=hello)print", "error parsing regexp: invalid or unsupported Perl syntax: `(?<`", []Record{}},
	}

	testDB, err := Open("testdata/search.db")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
	}

	testDB.Next("benchmark")
	for i, test := range tests {
		result, path, _, err := testDB.SearchRegexp(test.request)
		if err != nil {
			if err.Error() != test.err {
				t.Error(err)
			}
			continue
		}

		if path != "/benchmark" {
			t.Errorf("Test #%d Bad path: %s", i, path)
		}
		if len(result) != len(test.answer) {
			t.Errorf("Test #%d Bad size Want: %v Got %v", i, test.answer, result)
			continue
		}

		if !equal(test.answer, result) {
			t.Errorf("Test #%d Not equal. Want: %v Got: %v", i, test.answer, result)
			break
		}
	}
}
