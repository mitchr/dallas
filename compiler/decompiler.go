package compiler

import "strconv"

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
	for i, v := range b {
		if v == ' ' || v == '\r' || v == '\n' {
			b = append(b[:i], b[i+1:]...)
		}
	}

	splitSource := split(b)
	return splitSource
}

func split(b []byte) []uint16 {
	// split string into 2 hexNum (1 byte) string chunks
	// we want to convert to string here because adding two numbers will give us the sum, while adding two strings will concatenate them
	hexSource := []byte{}
	for _, v := range string(b) {
		c, _ := strconv.ParseUint(string(v), 16, 8)
		hexSource = append(hexSource, byte(c))
	}

	splitSource := []uint16{}
	for i, j := 0, 1; j < len(hexSource); i, j = i+2, j+2 {
		splitSource = append(splitSource, uint16(hexSource[i]<<4)|uint16(hexSource[j]))
	}

	// 2-byte elements exist as 2 separate elements right now
	// loop through and combine them so they exist as one element
	for i := 0; i < len(splitSource); i++ {
		if is2ByteDelimeter(byte(splitSource[i])) {
			splitSource[i] = (uint16(splitSource[i] << 8)) | uint16(splitSource[i+1])
			splitSource = append(splitSource[:i+1], splitSource[i+2:]...)
		}
	}

	return splitSource
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

func is2ByteDelimeter(b byte) bool {
	return b == 0x5d || b == 0x5e || b == 0x7e || b == 0xaa || b == 0xbb
}
