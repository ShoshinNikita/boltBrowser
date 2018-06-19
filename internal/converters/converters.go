package converters

import (
	"encoding/binary"
	"strconv"
)

// ConvertKey is a wrapper for converting key
func ConvertKey(b []byte) string {
	return fromString(b)
}

// ConvertValue is a wrapper for converting value
func ConvertValue(b []byte) string {
	return fromString(b)
}

func fromString(b []byte) string {
	return string(b)
}

func fromBigEndianUint64(b []byte) string {
	n := binary.BigEndian.Uint64(b)
	result := strconv.FormatUint(n, 10)
	return result
}

func fromLittleEndianUint64(b []byte) string {
	n := binary.LittleEndian.Uint64(b)
	result := strconv.FormatUint(n, 10)
	return result
}
