package parse

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseExpectedUint32le(reader io.Reader, param1 string, param2 string) (uint32, error) {

	if param1 != "uint32le" {
		return 0, fmt.Errorf("wrong type")
	}
	var b uint32
	err := binary.Read(reader, binary.LittleEndian, &b)
	return b, err
}

func parseExpectedUint16le(reader io.Reader, param1 string, param2 string) (uint16, error) {

	if param1 != "uint16le" {
		return 0, fmt.Errorf("wrong type")
	}
	var b uint16
	err := binary.Read(reader, binary.LittleEndian, &b)
	return b, err
}

func parseExpectedByte(reader io.Reader, param1 string, param2 string) (byte, error) {

	if param1 != "uint8" && param1 != "byte" {
		return 0, fmt.Errorf("wrong type")
	}
	// XXX "byte", params[2] describes a bit field
	var b byte
	err := binary.Read(reader, binary.LittleEndian, &b)
	return b, err
}

func parseExpectedBytes(layout *Layout, reader io.Reader, param1 string, param2 string) ([]byte, error) {

	p1 := strings.Split(param1, ":")

	if p1[0] != "byte" || len(p1) != 2 {
		return nil, fmt.Errorf("wrong type")
	}

	expectedLen, err := parseExpectedLen(p1[1])
	if err != nil {
		return nil, err
	}

	// "byte:3", params[2] holds the bytes
	buf, err := layout.parseByteN(reader, expectedLen)
	if err != nil {
		return nil, err
	}

	// split expected forms on comma
	expectedForms := strings.Split(param2, ",")
	for _, expectedForm := range expectedForms {

		expectedBytes := []byte(expectedForm)
		if int64(len(expectedForm)) == 2*expectedLen {
			// hex string?
			bytes, err := hex.DecodeString(expectedForm)
			if err == nil && byteSliceEquals(buf, bytes) {
				return expectedBytes, nil
			}
		}
		if string(buf) == string(expectedBytes) {
			return expectedBytes, nil
		}
	}

	return nil, fmt.Errorf("didnt find expected bytes %s", param2)
}

func parseExpectedLen(s string) (int64, error) {
	expectedLen, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if expectedLen > 255 {
		return 0, fmt.Errorf("len too big (max 255)")
	}
	if expectedLen <= 0 {
		return 0, fmt.Errorf("len too small (min 1)")
	}
	return expectedLen, nil
}

func byteSliceEquals(a []byte, b []byte) bool {

	if len(a) != len(b) {
		fmt.Println("error: a has len", len(a), " and b has len ", len(b))
		return false
	}

	for i, c1 := range a {
		c2 := b[i]
		if c1 != c2 {
			return false
		}
	}
	return true
}
