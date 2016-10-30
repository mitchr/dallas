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
		return splitUint16(e, binary.BigEndian)
	}
	return nil
}

// takes a pointer so we are able to mutate the slice without having to return it
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
	var sum uint16
	for _, v := range b {
		sum += uint16(v)
	}

	return splitUint16(sum, binary.LittleEndian)
}

func title(p string) []byte {
	// if title is less than 8, 0 pad the right
	title := make([]byte, 8)
	for i, r := range lex([]byte(p)) {
		// this will fail if p > 8, so do some bounds checking
		if i > 7 {
			break
		}
		title[i] = r
	}
	return title
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

	// len(u) = number of tokens present
	varData := append(splitUint16(uint16(len(u)), binary.LittleEndian), u...)

	var archive byte = 0x00
	if a {
		archive = 0x80
	}

	// len (null terminated)
	length := splitUint16(uint16(len(varData)), binary.LittleEndian)
	varEntry := concatBytes(0x0d, 0x00, length, 0x05, title(p), 0x00, archive, length, varData)

	return concatBytes(signature, comment, length, varEntry, checksum(varEntry))
}

// splitUint16 takes a uint16 and returns a []byte containing 2 8-bit elements ordered depending by the byte order argument
func splitUint16(u uint16, b binary.ByteOrder) []byte {
	g := new(bytes.Buffer)
	binary.Write(g, b, u)
	return g.Bytes()
}

func concatBytes(e ...interface{}) []byte {
	var slice []byte
	for _, v := range e {
		switch v.(type) {
		case byte:
			slice = append(slice, []byte{v.(byte)}...)
		case int:
			slice = append(slice, []byte{byte(v.(int))}...)
		case []byte:
			slice = append(slice, v.([]byte)...)
		}
	}
	return slice
}
