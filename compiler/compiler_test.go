package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCompile(t *testing.T) {
	f, err := ioutil.ReadFile("../basic_tests/test.tib")
	if err != nil {
		fmt.Println(err)
	}

	p, err := ioutil.ReadFile("../basic_tests/quad.tib")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", Compile(f, "TEST", false, false))
	fmt.Printf("%#v\n", Compile(p, "QUAD", false, false))
}

func TestDecompile(t *testing.T) {
	b, err := ioutil.ReadFile("../basic_tests/test.8xp")
	if err != nil {
		fmt.Println(err)
	}

	c, err := ioutil.ReadFile("../basic_tests/CHE.8xp")
	if err != nil {
		fmt.Println(err)
	}

	d, title := Decompile(b)
	fmt.Printf("data: %#v\ntitle: %#v\n", d, title)

	d, title = Decompile(c)
	fmt.Printf("data: %#v\ntitle: %#v\n", d, title)
}

// try to calculate the checksum from the given data
func TestChecksum(t *testing.T) {
	// s := []byte{0x0d, 0x00, 0x1a, 0x00, 0x05, 0xbb, 0xc4, 0xbb, 0xb4, 0xbb, 0xc3, 0xbb, 0xc4, 0x00, 0x00, 0x1a, 0x00, 0x18, 0x00, 0xde, 0x2a, 0x48, 0xbb, 0xb4, 0xbb, 0xbc, 0xbb, 0xbc, 0xbb, 0xbf, 0x29, 0x57, 0xbb, 0xbf, 0xbb, 0xc2, 0xbb, 0xbc, 0xbb, 0xb3, 0x2d, 0x2a, 0x3f}

	s := []byte{0x0d, 0x00, 0x12, 0x00, 0x05, 0x54, 0x45, 0x53, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x00, 0x10, 0x00, 0xde, 0x2a, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x29, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x2d, 0x2a, 0x3f}

	fmt.Println(len(s))
	fmt.Printf("%#v\n", s)

	var sum uint32
	for _, v := range s {
		sum += uint32(v)
	}

	low16 := uint16(sum & 0x0000ffff)
	// hi8 := byte(low16 & 0xff)
	// low8 := byte((low16 >> 8)) //& 0xff)

	// fmt.Printf("%x, %x\n", sum, low16)
	// fmt.Printf("%x, %x\n", hi8, low8)

	// this gives the same results
	c := new(bytes.Buffer)
	binary.Write(c, binary.LittleEndian, low16)

	if fmt.Sprintf("%x", c.Bytes()) != "4906" {
		t.Fail()
	}
}
