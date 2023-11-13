/*
	A Simple Syntax-Directed Translator

	Translates infix to postfix using recursive descent parsing
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Parser struct {
	lookahead byte
	reader    *bufio.Reader
}

func NewParser() (p *Parser) {
	var err error
	p = new(Parser)
	p.reader = bufio.NewReader(os.Stdin)
	p.lookahead, err = p.reader.ReadByte()

	if err != nil {
		log.Fatal(err)
	}
	return
}

func (p *Parser) expr() {
	p.term()

	for {
		switch p.lookahead {
		case '+':
			p.match('+')
			p.term()
			fmt.Print("+")
		case '-':
			p.match('-')
			p.term()
			fmt.Print("-")
		default:
			return
		}
	}
}

func (p *Parser) term() {
	if p.lookahead >= '0' && p.lookahead <= '9' {
		fmt.Printf("%c", p.lookahead)
		p.match(p.lookahead)
	} else {
		log.Fatal("Syntax Error")
	}
}

func (p *Parser) match(t byte) {
	var err error

	if p.lookahead == t {
		p.lookahead, err = p.reader.ReadByte()
	}

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	parser := NewParser()
	parser.expr()
	fmt.Println("")
}
