package main

import (
	"bufio"
	"fmt"
	"gcc/lexer"
	"gcc/parser"
	"log"
	"os"
)

func main() {
	var filename, output string
	switch len(os.Args) {
	case 1:
		fmt.Println("Usage gcc [filename]")
		return
	case 2:
		filename = os.Args[1]
		output = fmt.Sprintf("%s.out", filename[:len(filename)-2])
	case 4:
		if os.Args[2] != "-o" {
			fmt.Println("Usage gcc [filename] -o [output file]")
			return
		}
		output = os.Args[3]
	default:
		fmt.Println("Too many arguments to call, expected 1")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	inputFile, err := os.Open(fmt.Sprintf("%s/%s", wd, filename))
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()
	r := bufio.NewReader(inputFile)

	outputFile, err := os.Create(fmt.Sprintf("%s/%s", wd, output))
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	lex := lexer.NewLexer(r)
	p := parser.NewParser(lex, outputFile)
	p.Program()
}
