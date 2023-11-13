package main

import (
	"bufio"
	"fmt"
	"gcc/pkg"
	"log"
	"os"
)

func main() {
	var filename string
	switch len(os.Args) {
	case 1:
		fmt.Println("Usage gcc [filename]")
	case 2:
		filename = os.Args[1]
	case 3:

	default:
		fmt.Println("Too many arguments to call, expected 1")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", wd, filename))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewReader(file)

	lex := pkg.NewLexer(scanner)
	parser := pkg.NewParser(lex)
	parser.Program()
}
