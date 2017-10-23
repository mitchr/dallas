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

// decompiling the compilation should give back the original source code
func TestQuine(t *testing.T) {
	testFiles := []string{
		"../basic_tests/test.tib",
		"../basic_tests/CHE.tib",
		"../basic_tests/radical.tib",
	}

	for _, v := range testFiles {
		b, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fail()
		}

		out, _ := Decompile(Compile(b, "TEST", false, false))
		if string(out) != string(b) {
			t.Fail()
		}
	}
}

func TestParseNegOrMinus(t *testing.T) {
	tests := []struct {
		data []byte
		neg  bool
	}{
		{[]byte(`-3`), true},
		{[]byte(`(-3)`), true},
		{[]byte(`- 3`), true},
		{[]byte(`-(-3)`), true},

		{[]byte(`1-3`), false},
		{[]byte(`1 - 3`), false},
		{[]byte(`(1-3)`), false},
	}

	for _, v := range tests {
		neg := parseNegOrMinus(v.data)
		if neg != v.neg {
			t.Fail()
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
