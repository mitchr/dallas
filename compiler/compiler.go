package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// return err if no token could be found for the given input
func tokenMatch(b []byte) ([]byte, error) {
	if e, ok := Tokens[Token{string(b), false}]; ok {
		return []byte{e.(byte)}, nil
	}
	if e, ok := Tokens[Token{string(b), true}]; ok {
		return splitUint16(e.(uint16), binary.BigEndian), nil
	}
	return nil, fmt.Errorf("token not found for string: %s", string(b))
}

// tokenizes a []byte b and appends it to *[]byte t
// takes a pointer so we are able to mutate the slice without having to return it
func tokenize(b []byte) []byte {
	if e, _ := tokenMatch(b); e != nil {
		return e
	} else {
		// range over every rune and match it with a token
		var tokBuf []byte
		for _, v := range b {
			if e, err := tokenMatch([]byte{v}); err == nil {
				tokBuf = append(tokBuf, e...)
			}
		}
		return tokBuf
	}
}

func lex(b []byte) []byte {
	var tokBuf bytes.Buffer
	var curTok bytes.Buffer
	for i, v := range b {
		switch v {
		case '-':
			s := parseNegOrMinus(b[i-1 : i])
			tokBuf.Write(tokenize(curTok.Bytes()))
			tokBuf.WriteByte(s)
			curTok.Reset()
		case '\r', '\n', '+', '*', '/', '^', '=': // '-'
			// on Windows, carriage return contains 2 characters, but we only care about '\n'
			if v == '\r' {
				break
			}

			tokBuf.Write(tokenize(curTok.Bytes()))
			// depending on what v is, append it so it doesn't get lost
			// v is only 1 element, so it can only belong to oneBytes
			tokBuf.WriteByte(Tokens[Token{string(v), false}].(byte))
			curTok.Reset()
		case ' ': //'(':
			curTok.WriteByte(v)
			tokBuf.Write(tokenize(curTok.Bytes()))
			curTok.Reset()
		default:
			curTok.WriteByte(v)
			// last element of slice; AKA eof
			if i == len(b)-1 {
				tokBuf.Write(tokenize(curTok.Bytes()))
			}
		}
	}
	return tokBuf.Bytes()
}

// Check if the '-' is a negative, or if it is performing subtraction
// and return the correct respective byte for negation or subtraction
// taken from http://math.stackexchange.com/a/217316:
// - If you have numbers or variables on both sides of symbol −− then it means substraction.
// - If you have no number or variables before the symbol −− then it means negation.
//   Beware: parenthesis aren't variables.
func parseNegOrMinus(b []byte) byte {
	const NEG = 0xB0 // 176
	const SUB = 0x71 // 45

	// keep track of the previous rune in the sequence
	var previous byte
	for i, v := range b {
		if v == '-' {
			// check if the first index of b is a '-';
			if i == 0 || previous == '(' || previous == '\n' || previous == '\r' {
				return NEG
			}
		}
		previous = v
	}
	return SUB
}

// returns the lower 16 bits of the sum of all bytes in b
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
		// this will fail if len(p) > 8, so do some bounds checking
		if i > 7 {
			break
		}
		title[i] = r
	}
	return title
}

// f is the slice containing the source code
// p is the program title
// a sets the archive byte
// t sets the output device
// if f is too large, Compile might take an ungodly amount of time to finish
// streaming would be a more elegant solution to this
func Compile(f []byte, p string, a, t bool) []byte {
	identifier := "**TI83F*"
	if t {
		identifier = "**TI83**"
	}

	signature := append([]byte(identifier), 0x1a, 0x0a, 0x00)
	comment := make([]byte, 42) // []byte("COMPILED BY DALLAS")
	comment = append(comment, make([]byte, 42-len(comment))...)

	tokens := lex(f)

	// len(tokens) = number of tokens present
	varData := append(splitUint16(uint16(len(tokens)), binary.LittleEndian), tokens...)

	archive := 0x00
	if a {
		archive = 0x80
	}

	// len (null terminated)
	length := splitUint16(uint16(len(varData)), binary.LittleEndian)
	varEntry := concatBytes(0x0d, 0x00, length, 0x05, title(p), 0x00, archive, length, varData)

	return concatBytes(signature, comment, length, varEntry, checksum(varEntry))
}

// splitUint16 takes a uint16 and returns a []byte containing the hi and
// low 8-bits of u ordered depending on the byte order argument
func splitUint16(u uint16, b binary.ByteOrder) []byte {
	if b == binary.LittleEndian {
		// LSB first
		return []byte{byte(u & 0xff), byte(u >> 8)}
	}
	// MSB first
	return []byte{byte(u >> 8), byte(u & 0xff)}
}

func concatBytes(e ...interface{}) []byte {
	var temp []byte
	for _, v := range e {
		switch v.(type) {
		case byte:
			temp = append(temp, []byte{v.(byte)}...)
		case int:
			temp = append(temp, []byte{byte(v.(int))}...)
		case []byte:
			temp = append(temp, v.([]byte)...)
		}
	}
	return temp
}
