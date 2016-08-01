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

	// split string into 2 byte-size (kek) string chunks
	// we want to convert to string here becaus adding two numbers will give us the sum, while adding two strings wqill concatenate them
	splitB := []string{}
	for i, j := 0, 1; j < len(b); i, j = i+2, j+2 {
		splitB = append(splitB, string(b[i])+string(b[j]))
	}

	b = []byte{}
	for _, v := range splitB {
		c, _ := strconv.ParseUint(v, 16, 8)
		b = append(b, byte(c))
	}
	// remove last 2 elements from slice, aka the checksum
	b = b[:len(b)-1]
	b = b[:len(b)-1]

	// offset is the beginning of the data section
	const dataOffset = 0x37 + 0x11 + 0x02 // = 4a
	const titleOffset = 0x37 + 0x05

	var data []byte
	for i := dataOffset; i < len(b); i++ {
		if p, ok := revOneBytes[b[i]]; ok {
			data = append(data, []byte(p)...)
		}
	}

	var title []byte
	for i := titleOffset; i <= titleOffset+7; i++ {
		if p, ok := revOneBytes[b[i]]; ok {
			title = append(title, []byte(p)...)
		}
	}

	return data, title
}
