package internal

import (
	"fmt"
	"log"
)

var NodeLabels = 0

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
	fmt.Printf("L%d:", i)
}

func (node *Node) Emit(s string) {
	fmt.Printf("\t%s\n", s)
}
