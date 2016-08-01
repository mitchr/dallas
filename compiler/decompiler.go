package compiler

import (
	"strconv"
)

func Decompile(b []byte) ([]byte, []byte) {
	// remove newlines and spaces
	for i, v := range b {
		if v == byte('\n') || v == byte(' ') {
			// removes element from slice
			b = append(b[:i], b[i+1:]...)
		}
	}

	// split string into 2 hexNum (1 byte) string chunks
	// we want to convert to string here because adding two numbers will give us the sum, while adding two strings will concatenate them
	splitB := []string{}
	for i, j := 0, 1; j < len(b); i, j = i+2, j+2 {
		splitB = append(splitB, string(b[i])+string(b[j]))
	}

	for i := 0; i < len(splitB); i++ {
		// first bytes of 2-byte tokens
		if v := splitB[i]; v == "5d" || v == "5e" || v == "7e" || v == "aa" || v == "bb" {
			splitB[i] += splitB[i+1]
			splitB = append(splitB[:i+1], splitB[i+2:]...)
		}
	}

	cleanB := []uint16{}
	for _, v := range splitB {
		var c uint64
		if len(v) == 2 {
			c, _ = strconv.ParseUint(v, 16, 8)
		} else if len(v) == 4 {
			c, _ = strconv.ParseUint(v, 16, 16)
		}
		cleanB = append(cleanB, uint16(c))
	}

	// remove last 2 elements from slice, aka the checksum
	cleanB = cleanB[:len(cleanB)-2]

	const dataOffset = 0x37 + 0x11 + 0x02 // = 4a
	const titleOffset = 0x37 + 0x05

	var data []byte
	for i := dataOffset; i < len(cleanB); i++ {
		data = append(data, []byte(backwardsLex(cleanB[i]))...)
	}

	var title []byte
	for i := titleOffset; i <= titleOffset+7; i++ {
		title = append(title, []byte(backwardsLex(cleanB[i]))...)
	}

	return data, title
}

func backwardsLex(u uint16) string {
	if p, ok := revTwoBytes[u]; ok {
		return p
	}
	if p, ok := revOneBytes[byte(u)]; ok {
		return p
	}
	return ""
}
