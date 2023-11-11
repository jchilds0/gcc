package dragonbook

func main() {
	lex := NewLexer()
	parser := NewParser(lex)
	parser.Program()
}
