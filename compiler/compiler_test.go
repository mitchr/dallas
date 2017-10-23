package compiler

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCompile(t *testing.T) {
	testFiles := []struct {
		inputFile      string
		expectedOutput string
	}{
		{"../basic_tests/test.tib", "../basic_tests/TEST (dallas).8xp"},
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

	testFiles := []struct {
		inputFile      string
		expectedOutput string
		title          string
	}{
		{"../basic_tests/TEST (dallas).8xp", "../basic_tests/test.tib", "TEST"},

		{"../basic_tests/CHE (sourcecoder).8xp", "../basic_tests/CHE.tib", "CHE"},
		{"../basic_tests/CHE (tokens).8xp", "../basic_tests/CHE.tib", "CHE"},
		{"../basic_tests/CHE (dallas).8xp", "../basic_tests/CHE.tib", "CHE"},

		{"../basic_tests/RAD (sourcecoder).8xp", "../basic_tests/radical.tib", "RAD"},
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

		d, title := Decompile(b)

		if string(d) != string(e) {
			fmt.Println(string(d))
			t.Fail()
		}

		if string(title) != v.title {
			t.Fail()
		}
	}
}




	}

	}
}

func TestTitle(t *testing.T) {
	tests := []string{
		"AVALANCHE",
		"BAMBOOZLE",
		"ZIGZAGGING",
		"BEDAZZLING",
		"LUMBERJACK",
		"MOZZARELLA",
		"JAYWALKING",
		"a",
		",",
		"~",
		"asdlhfiuwehqcwo8rqw89enryc9823nc792c29378c32789q2789y3q278932y7r82q3yr78932qyr98",
		"789067907",
	}

	for _, v := range tests {
		if len(title(v)) > 8 {
			t.Fail()
		}
	}
}
