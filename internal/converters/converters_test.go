package converters

import (
	"encoding/binary"
	"testing"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		b      []byte
		result string
	}{
		{[]byte("Hello"), "Hello"},
		{[]byte("30.03.2001"), "30.03.2001"},
		{[]byte("Привет, мир!"), "Привет, мир!"},
	}

	for i, test := range tests {
		res := fromString(test.b)
		if res != test.result {
			t.Errorf("Test #%d Want: %s Got: %s", i, test.result, res)
		}
	}
}

func TestFromBigEndianUint64(t *testing.T) {
	type testStruct struct {
		b      []byte
		result string
	}
	var tests []testStruct

	number := []uint64{15, 841351, 513, 151, 484, 84, 153, 15}
	results := []string{"15", "841351", "513", "151", "484", "84", "153", "15"}

	for i := range number {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, number[i])
		tests = append(tests, testStruct{b, results[i]})
	}

	for i, test := range tests {
		res := fromBigEndianUint64(test.b)
		if res != test.result {
			t.Errorf("Test #%d Want: %s Got: %s", i, test.result, res)
		}
	}
}

func TestFromLittleEndianUint64(t *testing.T) {
	type testStruct struct {
		b      []byte
		result string
	}
	var tests []testStruct

	number := []uint64{15, 841351, 513, 151, 484, 84, 153, 15}
	results := []string{"15", "841351", "513", "151", "484", "84", "153", "15"}

	for i := range number {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, number[i])
		tests = append(tests, testStruct{b, results[i]})
	}

	for i, test := range tests {
		res := fromLittleEndianUint64(test.b)
		if res != test.result {
			t.Errorf("Test #%d Want: %s Got: %s", i, test.result, res)
		}
	}
}
