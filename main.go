package main

import (
	"bufio"
	"fmt"
	"gcc/internal"
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

	file, err := os.Open(fmt.Sprintf("%s/%s", wd, filename))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)

	file, err = os.Open(fmt.Sprintf("%s/%s", wd, output))
	if err != nil {
		file, err = os.Create(fmt.Sprintf("%s/%s", wd, output))
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	lex := internal.NewLexer(r)
	parser := internal.NewParser(lex, w)
	parser.Program()
}
