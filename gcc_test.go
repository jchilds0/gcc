package main

import (
	"bufio"
	"fmt"
	"gcc/lexer"
	"gcc/parser"
	"log"
	"os"
	"testing"
)

func TestGcc(t *testing.T) {
	inputs := []string{"test_program"}
	outputs := []string{"test_correct"}

	for i := 0; i < len(inputs); i++ {
		testSingleFile(inputs[i], outputs[i], t)
	}
}

func testSingleFile(input string, test string, t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	output := compileFile(input)

	testFile, err := os.Open(fmt.Sprintf("%s/testing/%s.out", wd, test))
	if err != nil {
		log.Fatalln(err)
	}
	defer testFile.Close()

	outputFile, err := os.Open(fmt.Sprintf("%s/testing/%s", wd, output))
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	scanOutput := bufio.NewScanner(outputFile)
	scanTest := bufio.NewScanner(testFile)

	i := 1
	for scanOutput.Scan() {
		if scanTest.Scan() == false {
			t.Error("Compiled file has too many lines")
		}

		if scanTest.Text() != scanOutput.Text() {
			t.Errorf("Compiled file and test differ on line %d\n Compiled file: %s\n Test: %s\n", i, scanOutput.Text(), scanTest.Text())
		}
		i++
	}

	if scanTest.Scan() != false {
		t.Error("Compiled file has too few lines")
	}
}

func compileFile(input string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	inputFile, err := os.Open(fmt.Sprintf("%s/testing/%s.c", wd, input))
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()
	r := bufio.NewReader(inputFile)

	outputFile, err := os.Create(fmt.Sprintf("%s/testing/test.out", wd))
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	lex := lexer.NewLexer(r)
	p := parser.NewParser(lex, outputFile)
	p.Program()

	return "test.out"
}
