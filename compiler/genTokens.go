// +build ignore

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
)

func main() {

	type Token struct {
		Name       string
		Hex        string
		DoubleByte bool
	}

	file, err := os.Open("tokens.txt")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	tokens := []Token{}
	myTokens := map[string]string{}

	for {
		token, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		double := false
		if len(token[1]) > 4 {
			double = true
		}

		tokens = append(tokens, Token{token[0], token[1], double})
		myTokens[token[0]] = token[1]
	}

	fileOut := `package compiler

	var Tokens = map[string]string { {{range $i, $v := .}}
		"{{.Name}}": "{{.Hex}}",{{end}}
	}

	var revTokens = func(m map[string]string) map[string]string {
		n := make(map[string]string)
		for k, v := range m {
			n[v] = k
		}
		return n
	}(Tokens)
	`

	templ, err := template.New("tokens").Parse(fileOut)
	if err != nil {
		log.Fatal(err)
	}

	err = templ.Execute(os.Stdout, tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(myTokens)
}
