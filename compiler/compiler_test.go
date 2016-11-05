package compiler

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"testing"
)

type testData struct {
	inputFile      string
	expectedOutput string
}

func TestCompile(t *testing.T) {
	testFiles := []testData{
		{"../basic_tests/test.tib", "../basic_tests/test.8xp"},
		// {"../basic_tests/quad.tib", ""},
		// {"../basic_tests/radical.tib", ""},
	}

	for _, v := range testFiles {
		b, err := ioutil.ReadFile(v.inputFile)
		if err != nil {
			t.Skip("Test file not found: ", v.inputFile)
		}

		e, err := ioutil.ReadFile(v.expectedOutput)
		if err != nil {
			t.Skip("Test file not found: ", v.expectedOutput)
		}

		d := Compile(b, "TEST", false, false)

		if string(d) != string(e) {
			t.Fail()
		}
	}
}

func TestDecompile(t *testing.T) {
	testFiles := []testData{
		{"../basic_tests/test.8xp", "../basic_tests/test.tib"},
		{"../basic_tests/CHE.8xp", "../basic_tests/CHE.tib"},
		// {"../basic_tests/radical.8xp", ""},
	}
	// testRawData := [][]byte{[]byte("ClearEntries")}

	for _, v := range testFiles {
		b, err := ioutil.ReadFile(v.inputFile)
		if err != nil {
			t.Skip("Test file not found: ", v.inputFile)
		}

		e, err := ioutil.ReadFile(v.expectedOutput)
		if err != nil {
			t.Skip("Test file not found: ", v.expectedOutput)
		}

		d, _ := Decompile(b)

		if string(d) != string(e) {
			t.Fail()
		}
	}
}

// calculate the checksum from the given data
func TestChecksum(t *testing.T) {
	// s := []byte{0x0d, 0x00, 0x1a, 0x00, 0x05, 0xbb, 0xc4, 0xbb, 0xb4, 0xbb, 0xc3, 0xbb, 0xc4, 0x00, 0x00, 0x1a, 0x00, 0x18, 0x00, 0xde, 0x2a, 0x48, 0xbb, 0xb4, 0xbb, 0xbc, 0xbb, 0xbc, 0xbb, 0xbf, 0x29, 0x57, 0xbb, 0xbf, 0xbb, 0xc2, 0xbb, 0xbc, 0xbb, 0xb3, 0x2d, 0x2a, 0x3f}

	s := []byte{0x0d, 0x00, 0x12, 0x00, 0x05, 0x54, 0x45, 0x53, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x00, 0x10, 0x00, 0xde, 0x2a, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x29, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x2d, 0x2a, 0x3f}

	fmt.Println(len(s))
	fmt.Printf("%#v\n", s)

	var sum uint16
	for _, v := range s {
		sum += uint16(v)
	}

	low16 := splitUint16(sum, binary.LittleEndian)
	if fmt.Sprintf("%x", low16) != "4906" {
		t.Fail()
	}
}

func TestTitle(t *testing.T) {
	fmt.Println(title("WHERE"))
}
