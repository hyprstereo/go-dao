package dao

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// #nosec G103
// GetString returns a string pointer without allocation
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// #nosec G103
// GetBytes returns a byte pointer without allocation
func UnsafeBytes(s string) (bs []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return
}

// SafeString copies a string to make it immutable
func SafeString(s string) string {
	return string(UnsafeBytes(s))
}

// SafeBytes copies a slice to make it immutable
func SafeBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

const (
	uByte = 1 << (10 * iota)
	uKilobyte
	uMegabyte
	uGigabyte
	uTerabyte
	uPetabyte
	uExabyte
)

func ByteSize(b any) string {
	var bytes uint64
	switch val := b.(type) {
	case uint64:
		bytes = val
	case []byte:
		bytes = uint64(len(val))
	case string:
		bytes = uint64(len(UnsafeBytes(val)))
	default:
		return "ERROR:on supports uint64, []bytes or string"
	}
	unit := ""
	value := float64(bytes)
	switch {
	case bytes >= uExabyte:
		unit = "EB"
		value = value / uExabyte
	case bytes >= uPetabyte:
		unit = "PB"
		value = value / uPetabyte
	case bytes >= uTerabyte:
		unit = "TB"
		value = value / uTerabyte
	case bytes >= uGigabyte:
		unit = "GB"
		value = value / uGigabyte
	case bytes >= uMegabyte:
		unit = "MB"
		value = value / uMegabyte
	case bytes >= uKilobyte:
		unit = "KB"
		value = value / uKilobyte
	case bytes >= uByte:
		unit = "B"
	default:
		return "0B"
	}
	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// Deprecated fn's

// #nosec G103
// GetString returns a string pointer without allocation
func GetString(b []byte) string {
	return UnsafeString(b)
}

// #nosec G103
// GetBytes returns a byte pointer without allocation
func GetBytes(s string) []byte {
	return UnsafeBytes(s)
}

// ImmutableString copies a string to make it immutable
func ImmutableString(s string) string {
	return SafeString(s)
}
