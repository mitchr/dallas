package compiler

import (
	"strconv"
	"strings"
)

func Decompile(b []byte) ([]byte, []byte) {
	// check if b was built from dallas, or some other compiler
	// for some reason, other compilers spit out each byte as a pair of 2 ascii characters, where as dallas prints both characters together as an element of a []byte
	cleanSource := make([]uint16, len(b))
	if b[0] == 50 {
		// first byte is interpreted as "2"
		cleanSource = clean(b)
	} else if b[0] == 42 {
		// first byte is interpreted as "2a", or "*" (compiled by dallas)
		for i, v := range b {
			cleanSource[i] = uint16(v)
		}
	}

	// remove last 2 elements from slice, aka the checksum
	cleanSource = cleanSource[:len(cleanSource)-2]

	const dataOffset = 0x37 + 0x11 + 0x02 // = 4a
	const titleOffset = 0x37 + 0x05

	var data []byte
	for i := dataOffset; i < len(cleanSource); i++ {
		data = append(data, []byte(backwardsLex(cleanSource[i]))...)
	}

	var title []byte
	for i := titleOffset; i <= titleOffset+7; i++ {
		title = append(title, []byte(backwardsLex(cleanSource[i]))...)
	}

	return data, title
}

func clean(b []byte) []uint16 {
	// remove spaces and line endings
	// for Windows compatibility, check for "\r\n"
	newB := string(b)
	for _, v := range []string{" ", "\r\n", "\n"} {
		newB = strings.Replace(newB, string(v), "", -1)
	}

	// split string into 2 hexNum (1 byte) string chunks
	// we want to convert to string here because adding two numbers will give us the sum, while adding two strings will concatenate them
	splitB := []string{}
	for i, j := 0, 1; j < len(newB); i, j = i+2, j+2 {
		splitB = append(splitB, string(newB[i])+string(newB[j]))
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

	return cleanB
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
