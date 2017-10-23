package compiler

import (
	"fmt"
	"log"
)

// returns the decompiled program along with the title
func Decompile(b []byte) ([]byte, []byte) {
	const dataOffset = 0x37 + 0x11 + 0x02
	const titleOffset = 0x37 + 0x05

	var data []byte
	// dataOffset to the end minus the 2 bytes for the checksum
	for _, v := range combine2ByteElements(b[dataOffset : len(b)-2]) {
		k, err := backwardsLex(v)
		if err != nil {
			// should log.Fatal here once the token mapping is all figured out
			fmt.Println(err)
		}
		data = append(data, []byte(k)...)
	}

	var title []byte
	for _, v := range combine2ByteElements(b[titleOffset : titleOffset+7]) {
		// if we find a 0x00, then the title is less than 8 characters
		// so we can stop
		if v == 0 {
			break
		}
		k, err := backwardsLex(v)
		if err != nil {
			log.Fatal(err)
		}
		title = append(title, []byte(k)...)
	}

	return data, title
}

// 2-byte elements are currently expressed as 2 separate elements.
// (i.e. []uint16{0xbb01} is []byte{0xbb, 0x01})
// loop through and combine these separate bytes into a single uint16
func combine2ByteElements(b []byte) []uint16 {
	source := []uint16{}
	for i := 0; i < len(b); i++ {
		if is2ByteDelimeter(b[i]) {
			// combine the 2 elements
			source = append(source, (uint16(b[i])<<8)|uint16(b[i+1]))
			// advance one space so we skip the second half of the token
			// that was already consumed by the above line
			i = i + 1
			continue
		}
		source = append(source, uint16(b[i]))
	}
	return source
}

func backwardsLex(u uint16) (string, error) {
	if p, ok := revTwoBytes[u]; ok {
		return p, nil
	}
	if p, ok := revOneBytes[byte(u)]; ok {
		return p, nil
	}
	return "", fmt.Errorf("token not found for key %#x", u)
}

func is2ByteDelimeter(b byte) bool {
	return b == 0x5d || b == 0x5e || b == 0x63 || b == 0x7e || b == 0xaa || b == 0xbb
}
