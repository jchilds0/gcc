package pkg

import (
	"fmt"
	"log"
)

type Node struct {
	lexline int
	labels  int
}

func NewNode() *Node {
	return &Node{lexline: 0, labels: 0}
}

func (node *Node) Error(s string) {
	log.Fatalf("Error near line %d: %s", node.lexline, s)
}

func (node *Node) NewLabel() int {
	node.labels++
	return node.labels
}

func (node *Node) EmitLabel(i int) {
	fmt.Printf("L%d:", i)
}

func (node *Node) Emit(s string) {
	fmt.Printf("\t%s", s)
}
