package compiler

import (
	"bytes"
	"encoding/binary"
)

func isKnownToken(b []byte) []byte {
	if e, ok := oneBytes[string(b)]; ok {
		return []byte{e}
	}
	if e, ok := twoBytes[string(b)]; ok {
		// not sure why this needs to be in BidEndian, but otherwise it's fucked
		return split(e, binary.BigEndian)
	}
	return nil
}

func tokenize(b []byte, t *[]byte) {
	if e := isKnownToken(b); e != nil {
		*t = append(*t, e...)
	} else {
		// range over every rune and match it with a token
		for _, v := range b {
			if e := isKnownToken([]byte{v}); e != nil {
				*t = append(*t, e...)
			}
		}
	}
}

func lex(b []byte) []byte {
	var tokBuf []byte
	var curTok []byte
	for i, v := range b {
		switch v {
		case '\r', '\n', '+', '-', '*', '/', '^', '=':
			// on Windows, carriage return contains 2 characters, but we only care about '\n'
			if v == '\r' {
				break
			}

			tokenize(curTok, &tokBuf)
			// depending on what v is, append it so it doesn't get lost
			tokBuf = append(tokBuf, oneBytes[string(v)])
			curTok = []byte{}
		case ' ', '(':
			curTok = append(curTok, v)
			tokenize(curTok, &tokBuf)
			curTok = []byte{}
		default:
			curTok = append(curTok, v)
			// last element of slice; AKA eof
			if i == len(b)-1 {
				tokenize(curTok, &tokBuf)
			}
		}
	}
	return tokBuf
}

func checksum(b []byte) []byte {
	var sum uint32
	for _, v := range b {
		sum += uint32(v)
	}

	// mask the upper 16 bits
	return split(uint16(sum), binary.LittleEndian)
}

func Compile(f []byte, p string, a bool, t bool) []byte {
	var identifier string
	if t {
		identifier = "**TI83**"
	} else {
		identifier = "**TI83F*"
	}

	signature := append([]byte(identifier), 0x1a, 0x0a, 0x00)
	comment := make([]byte, 42) // []byte("TESTCOMPILE")
	comment = append(comment, make([]byte, 42-len(comment))...)

	u := lex(f)
	// test data so i can work on other parts of this monstrosity
	// u := []byte{0xde, 0x2a, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x29, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x2d, 0x2a, 0x3f}

	// len(u) = number of tokens present
	varData := append(split(uint16(len(u)), binary.LittleEndian), u...)

	var archive byte = 0x00
	if a {
		archive = 0x80
	}
	// split(archive, binary.BigEndian)

	// if title is less than 8, 0 pad the right
	title := make([]byte, 8)
	for i, r := range lex([]byte(p)) {
		// this will fail if p > 8, so do some bounds checking
		if i > 7 {
			break
		}
		title[i] = r
	}

	// len (null terminated)
	length := split(uint16(len(varData)), binary.LittleEndian)
	varEntry := concatSlices([]byte{0x0d, 0x00}, length, []byte{0x05}, title, []byte{0x00, archive}, length, varData)
	checksum := checksum(varEntry)

	return concatSlices(signature, comment, length, varEntry, checksum)
}

// split takes a uint16 and splits it into a []byte containing the highest 8 bits and lowest 8 bits
func split(u uint16, b binary.ByteOrder) []byte {
	g := new(bytes.Buffer)
	binary.Write(g, b, u)
	return g.Bytes()
}

func concatSlices(slices ...[]byte) []byte {
	var newSlice []byte
	for _, v := range slices {
		newSlice = append(newSlice, v...)
	}
	return newSlice
}
