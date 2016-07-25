package main

import (
	"flag"
	"fmt"
	"github.com/Mitchell-Riley/dallas/compiler"
	"io/ioutil"
	"os"
)

var (
	disAsm   = flag.Bool("d", false, "set to true to disassemble .8xp files")
	progName = flag.String("p", "PROG", "set the program name")
	outName  = flag.String("o", *progName+".8xp", "set the name of the output file")
	archive  = flag.Bool("a", false, "set the archive bit; if false, ram is used to store the program")
	lock     = flag.Bool("e", false, "set the edit-lock bit")
	help     = flag.Bool("h", false, "display this help message")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "dallas is a TI-BASIC Compiler and Decompiler\n\nUsage:\n\tdallas [flags] filename\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	inName := flag.Arg(0)
	fmt.Println(inName)

	if inName == "" || *help {
		flag.Usage()
	} else {
		//this is all in the else block temporarily
		inFile, err := ioutil.ReadFile(inName)
		if err != nil {
			fmt.Println(err)
		}

		var output []byte
		if *disAsm == true {
			var b []byte
			output, b = compiler.Decompile(inFile)
			*progName = string(b)
			*outName = *progName + ".tib"
		} else {
			output = compiler.Compile(inFile, *progName, *archive)
		}

		outFile, err := os.Create(*outName)
		if err != nil {
			fmt.Println(err)
		}

		_, err = outFile.Write(output)
		if err != nil {
			fmt.Println(err)
		}
	}
}
