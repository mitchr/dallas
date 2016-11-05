package compiler

import "strconv"

func Decompile(b []byte) ([]byte, []byte) {
	// check if b was built from dallas, or some other compiler
	// for some reason, other compilers spit out each byte as a pair of 2 ascii characters, where as dallas prints both characters together as an element of a []byte
	cleanSource := make([]uint16, len(b))
	// first byte is interpreted as "2"
	if b[0] == 50 {
		// split string into 2 hexNum (1 byte) string chunks
		// we want to convert to string here because adding two numbers will give us the sum, while adding two strings will concatenate them
		b = clean(b)

		hexSource := []byte{}
		for _, v := range string(b) {
			c, _ := strconv.ParseUint(string(v), 16, 8)
			hexSource = append(hexSource, byte(c))
		}

		// combine
		splitSource := []byte{}
		for i, j := 0, 1; j < len(hexSource); i, j = i+2, j+2 {
			splitSource = append(splitSource, byte(hexSource[i]<<4)|byte(hexSource[j]))
		}

		cleanSource = split(splitSource)

	} else if b[0] == 0x2a {
		// first byte is interpreted as "2a", or "*" (compiled by dallas)
		for i, v := range b {
			cleanSource[i] = uint16(v)
		}

		cleanSource = split(clean(b))
	}

	// now we subtract 0x02 from each offset for some reason
	// I legitimately do not know why, but I won't question it
	const dataOffset = 0x37 + 0x11 + 0x02 - 0x02
	const titleOffset = 0x37 + 0x05 - 0x02

	var data []byte
	// minus the checksum
	for _, v := range cleanSource[dataOffset : len(cleanSource)-2] {
		data = append(data, []byte(backwardsLex(v))...)
	}

	var title []byte
	for _, v := range cleanSource[titleOffset : titleOffset+7] {
		title = append(title, []byte(backwardsLex(v))...)
	}

	return data, title
}

func clean(b []byte) []byte {
	// remove spaces and line endings
	// for Windows compatibility, check for "\r\n"
	for i, v := range b {
		if v == ' ' || v == '\r' || v == '\n' {
			b = append(b[:i], b[i+1:]...)
		}
	}
	return b
}

func split(b []byte) []uint16 {
	// 2-byte elements exist as 2 separate elements right now
	// loop through and combine them so they exist as one element
	splitSource := make([]uint16, len(b))
	for i, v := range b {
		splitSource[i] = uint16(v)
	}

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
