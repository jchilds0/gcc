package parser

import (
	"fmt"
	"io"
	"log"
)

var NodeLabels = 0
var NodeWrite io.Writer

type Node struct {
	lexline int
}

func NewNode() *Node {
	return &Node{lexline: 0}
}

func (node *Node) Error(s string) {
	log.Fatalf("Error near line %d: %s", node.lexline, s)
}

func (node *Node) NewLabel() int {
	NodeLabels++
	return NodeLabels
}

func (node *Node) EmitLabel(i int) {
	_, err := NodeWrite.Write([]byte(fmt.Sprintf("L%d:", i)))
	if err != nil {
		log.Fatal(err)
	}
}

func (node *Node) Emit(s string) {
	_, err := NodeWrite.Write([]byte(fmt.Sprintf("\t%s\n", s)))
	if err != nil {
		log.Fatal(err)
	}
}
