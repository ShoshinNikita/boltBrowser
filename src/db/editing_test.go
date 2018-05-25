package db_test

import (
	"testing"

	. "db"
)

// Test db in testdata/edit.db
//
// Structure:
// root
// 	  |-> "hello"
// 	  |-> "xyz"
// 		  | "byte" - "15"
// 		  | "hi" - "yeah"
//
// The db has same structure after all tests

var allBuckets = []struct {
	path    []string
	records []Record
}{
	{[]string{}, []Record{bckt("hello"), bckt("xyz")}},
	{[]string{"hello"}, []Record{}},
	{[]string{"xyz"}, []Record{rcrd("byte", "15"), rcrd("hi", "yeah")}},
}

// AddBucket and DeleteBucket()
func TestBucketsEditing(t *testing.T) {
	SetOffset(100)

	addingTests := []struct {
		path    []string
		name    string
		err     error
		records []Record
	}{
		{[]string{}, "123", newErr(""),
			[]Record{bckt("123"), bckt("hello"), bckt("xyz")}},
		{[]string{"hello"}, "546", newErr(""),
			[]Record{bckt("546")}},
		{[]string{"hello", "546"}, "1", newErr(""),
			[]Record{bckt("1")}},
		{[]string{"xyz"}, "byte", newErr("\"byte\" is a record"),
			[]Record{rcrd("byte", "15"), rcrd("hi", "yeah")}},
		{[]string{}, "hello", newErr("bucket already exist"),
			[]Record{bckt("123"), bckt("hello"), bckt("xyz")}},
	}

	deletingTests := []struct {
		path []string
		name string
		err  error
	}{
		{[]string{}, "123", newErr("")},
		{[]string{"hello", "546"}, "1", newErr("")},
		{[]string{"hello"}, "546", newErr("")},
		{[]string{"xyz"}, "byte", newErr("\"byte\" is a record")},
	}

	testDB, err := Open("testdata/edit.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Testing of adding buckets
	t.Log("Testing of adding buckets")
	for i, test := range addingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}
		err := testDB.AddBucket(test.name)

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
		} else if err != nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
			}
		}

		testDB.ClearPath()
	}

	// Testing are there new buckets
	t.Log("Testing are there new buckets")
	for i, test := range addingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}

	// Testing of deleting buckets
	t.Log("Testing of deleting buckets")
	for i, test := range deletingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}
		err := testDB.DeleteBucket(test.name)

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
		} else if err != nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
			}
		}

		testDB.ClearPath()
	}

	// Testing is a new file equal to the old file
	t.Log("Checking buckets and records")
	for i, test := range allBuckets {
		for _, s := range test.path {
			testDB.Next(s)
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}
}

func TestRecordEditing(t *testing.T) {
	SetOffset(100)

	addingTests := []struct {
		path    []string
		key     string
		value   string
		err     error
		records []Record
	}{
		{[]string{}, "123", "15", newErr(""),
			[]Record{bckt("hello"), bckt("xyz"), rcrd("123", "15")}},
		{[]string{}, "123", "16", newErr("\"123\" already exists"),
			[]Record{bckt("hello"), bckt("xyz"), rcrd("123", "15")}},
		{[]string{}, "hello", "5", newErr("\"hello\" is a bucket"),
			[]Record{bckt("hello"), bckt("xyz"), rcrd("123", "15")}},
		{[]string{"xyz"}, "hello", "1", newErr(""),
			[]Record{rcrd("byte", "15"), rcrd("hello", "1"), rcrd("hi", "yeah")}},
	}

	deletingTests := []struct {
		path []string
		key  string
		err  error
	}{
		{[]string{}, "123", newErr("")},
		{[]string{"xyz"}, "hello", newErr("")},
		{[]string{}, "hello", newErr("\"hello\" is a bucket")},
		{[]string{"hello"}, "123", newErr("\"123\" doesn't exist")},
	}

	testDB, err := Open("testdata/edit.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Testing adding of records
	t.Log("Testing adding of records")
	for i, test := range addingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}
		err := testDB.AddRecord(test.key, test.value)

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
		} else if err != nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
			}
		}

		testDB.ClearPath()
	}

	// Testing are there new records
	t.Log("Testing are there new records")
	for i, test := range addingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}

	// Testing of deleting records
	t.Log("Testing of deleting records")
	for i, test := range deletingTests {
		for _, s := range test.path {
			testDB.Next(s)
		}
		err := testDB.DeleteRecord(test.key)

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
		} else if err != nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
			}
		}

		testDB.ClearPath()
	}

	// Testing is a new file equal to the old file
	t.Log("Checking buckets and records")
	for i, test := range allBuckets {
		for _, s := range test.path {
			testDB.Next(s)
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}
}

func TestEditRecord(t *testing.T) {
	tests := []struct {
		path     []string
		oldKey   string
		newKey   string
		newValue string
		err      error
		records  []Record
	}{
		// 	{byte 15} -> {byte 35}
		// 	[hi yeah] -> {hello 88}
		{[]string{"xyz"}, "byte", "byte", "35", newErr(""),
			[]Record{rcrd("byte", "35"), rcrd("hi", "yeah")}},
		{[]string{"xyz"}, "hi", "hello", "88", newErr(""),
			[]Record{rcrd("byte", "35"), rcrd("hello", "88")}},
		{[]string{}, "hello", "hi", "35", newErr("\"hello\" is a bucket"),
			[]Record{bckt("hello"), bckt("xyz")}},
		{[]string{}, "test", "test1", "15", newErr("\"test\" doesn't exist"),
			[]Record{bckt("hello"), bckt("xyz")}},
		// return default values
		{[]string{"xyz"}, "byte", "byte", "15", newErr(""),
			[]Record{rcrd("byte", "15"), rcrd("hello", "88")}},
		{[]string{"xyz"}, "hello", "hi", "yeah", newErr(""),
			[]Record{rcrd("byte", "15"), rcrd("hi", "yeah")}},
	}

	testDB, err := Open("testdata/edit.db")
	defer testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Testing of editing recrods
	t.Log("Testing of editing recrods")
	for i, test := range tests {
		for _, s := range test.path {
			testDB.Next(s)
		}
		err := testDB.EditRecord(test.oldKey, test.newKey, test.newValue)

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
		} else if err != nil && test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("Test #%d Want: %v Got: %v", i, test.err, err)
			}
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}

	// Checking buckets and records
	t.Log("Checking buckets and records")
	for i, test := range allBuckets {
		for _, s := range test.path {
			testDB.Next(s)
		}

		res, err := testDB.GetCurrent()
		if err != nil {
			t.Error(err)
			continue
		}

		if !equal(test.records, res.Records) {
			t.Errorf("Test #%d Want: %v Got: %v", i, test.records, res.Records)
		}

		testDB.ClearPath()
	}
}
