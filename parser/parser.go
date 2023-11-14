package parser

import (
	"fmt"
	"gcc/lexer"
	"io"
	"log"
)

type Parser struct {
	lex  *lexer.Lexer
	look lexer.Tokener
	top  *Env
	used int
}

func NewParser(lexer *lexer.Lexer, write io.Writer) *Parser {
	parser := &Parser{lex: lexer}
	parser.move()
	parser.top = NewEnv(nil)
	NodeWrite = write

	return parser
}

func (parser *Parser) move() {
	parser.look = parser.lex.Scan()
}

func (parser *Parser) error(s string) {
	log.Fatalf("near line %d: %s", parser.lex.Line, s)
}

func (parser *Parser) match(t int) {
	if parser.look.GetTokenTag() == t {
		parser.move()
	} else {
		parser.error("syntax error")
	}
}

func (parser *Parser) Program() {
	StmtNull = &Stmt{}
	StmtEnclosing = &Stmt{}

	s := parser.block()
	begin := s.GetNode().NewLabel()
	after := s.GetNode().NewLabel()
	s.GetNode().EmitLabel(begin)
	s.Gen(begin, after)
	s.GetNode().EmitLabel(after)
}

func (parser *Parser) block() Stmter {
	parser.match('{')
	savedEnv := parser.top
	parser.top = NewEnv(parser.top)
	parser.decls()
	s := parser.stmts()
	parser.match('}')
	parser.top = savedEnv
	return s
}

func (parser *Parser) decls() {
	for parser.look.GetTokenTag() == lexer.BASIC {
		p := parser.types()
		tok := parser.look
		parser.match(lexer.ID)
		parser.match(';')
		id := NewId(lexer.NewWord(tok.GetTokenTag(), tok.String()), p, parser.used)
		parser.top.Put(tok.String(), id)
		parser.used += p.GetWidth()
	}
}

func (parser *Parser) types() lexer.Typer {
	p := lexer.NewType(parser.look.GetTokenTag(), parser.look.String(), parser.look.(lexer.Typer).GetWidth())
	parser.match(lexer.BASIC)
	if parser.look.GetTokenTag() != '[' {
		return p
	} else {
		return parser.dims(p)
	}
}

func (parser *Parser) dims(p lexer.Typer) lexer.Typer {
	parser.match('[')
	parser.match(lexer.NUM)
	parser.match(']')

	if parser.look.GetTokenTag() == '[' {
		p = parser.dims(p)
	}

	return lexer.NewArray(0, p)
}

func (parser *Parser) stmts() Stmter {
	if parser.look.GetTokenTag() == '}' {
		return StmtNull
	} else {
		return NewSeq(parser.stmt(), parser.stmts())
	}
}

func (parser *Parser) stmt() Stmter {
	var x Exprer
	var s1, s2 Stmter

	switch parser.look.GetTokenTag() {
	case ';':
		parser.move()
		return StmtNull
	case lexer.IF:
		parser.match(lexer.IF)
		parser.match('(')
		x = parser.bool()
		parser.match(')')
		s1 = parser.stmt()
		if parser.look.GetTokenTag() != lexer.ELSE {
			return NewIf(x, s1)
		} else {
			parser.match(lexer.ELSE)
			s2 = parser.stmt()
			return NewElse(x, s1, s2)
		}
	case lexer.WHILE:
		whilenode := new(While)
		savedStmt := StmtEnclosing
		StmtEnclosing = Stmter(whilenode)
		parser.match(lexer.WHILE)
		parser.match('(')
		x = parser.bool()
		parser.match(')')
		s1 = parser.stmt()
		whilenode.Init(x, s1)
		StmtEnclosing = savedStmt
		return whilenode
	case lexer.DO:
		donode := new(Do)
		savedStmt := StmtEnclosing
		StmtEnclosing = donode.stmt
		parser.match(lexer.DO)
		s1 := parser.stmt()
		parser.match(lexer.WHILE)
		parser.match('(')
		x := parser.bool()
		parser.match(')')
		parser.match(';')
		donode.Init(s1, x)
		StmtEnclosing = savedStmt
		return donode
	case lexer.BREAK:
		parser.match(lexer.BREAK)
		parser.match(';')
		return NewBreak()
	case '{':
		return parser.block()
	default:
		return parser.assign()
	}
}

func (parser *Parser) assign() Stmter {
	var stmt Stmter
	t := parser.look
	parser.match(lexer.ID)
	id, err := parser.top.Get(t.String())
	if err != nil {
		parser.error(fmt.Sprintf("%s undeclared", t.String()))
	}

	if parser.look.GetTokenTag() == '=' {
		parser.move()
		stmt = NewSet(id, parser.bool())
	} else {
		x := parser.offset(id)
		parser.match('=')
		stmt = NewSetElem(x, parser.bool())
	}
	parser.match(';')
	return stmt
}

func (parser *Parser) bool() Exprer {
	x := parser.join()
	for parser.look.GetTokenTag() == lexer.OR {
		tok := parser.look
		parser.move()
		x = NewOr(tok, x, parser.join())
	}
	return x
}

func (parser *Parser) join() Exprer {
	x := parser.equality()
	for parser.look.GetTokenTag() == lexer.AND {
		tok := parser.look
		parser.move()
		x = NewAnd(tok, x, parser.equality())
	}
	return x
}

func (parser *Parser) equality() Exprer {
	x := parser.rel()
	for parser.look.GetTokenTag() == lexer.EQ || parser.look.GetTokenTag() == lexer.NE {
		tok := parser.look
		parser.move()
		x = NewRel(tok, x, parser.rel())
	}

	return x
}

func (parser *Parser) rel() Exprer {
	x := parser.expr()
	switch parser.look.GetTokenTag() {
	case '<', lexer.LE, lexer.GE, '>':
		tok := parser.look
		parser.move()
		return NewRel(tok, x, parser.expr())
	default:
		return x
	}
}

func (parser *Parser) expr() Exprer {
	x := parser.term()
	for parser.look.GetTokenTag() == '+' || parser.look.GetTokenTag() == '-' {
		tok := parser.look
		parser.move()
		x = NewArith(tok, x, parser.term())
	}
	return x
}

func (parser *Parser) term() Exprer {
	x := parser.unary()
	for parser.look.GetTokenTag() == '*' || parser.look.GetTokenTag() == '/' {
		tok := parser.look
		parser.move()
		x = NewArith(tok, x, parser.unary())
	}
	return x
}

func (parser *Parser) unary() Exprer {
	if parser.look.GetTokenTag() == '-' {
		parser.move()
		return NewUnary(lexer.WordMinus, parser.unary())
	} else if parser.look.GetTokenTag() == '!' {
		tok := parser.look
		parser.move()
		return NewNot(tok, parser.unary())
	} else {
		return parser.factor()
	}
}

func (parser *Parser) factor() Exprer {
	switch parser.look.GetTokenTag() {
	case '(':
		parser.move()
		x := parser.bool()
		parser.match(')')
		return x
	case lexer.NUM:
		x := NewConstant(parser.look, lexer.Int)
		parser.move()
		return x
	case lexer.REAL:
		x := NewConstant(parser.look, lexer.Float)
		parser.move()
		return x
	case lexer.TRUE:
		parser.move()
		return ConstantTrue
	case lexer.FALSE:
		parser.move()
		return ConstantFalse
	case lexer.ID:
		s := parser.look.String()
		id, err := parser.top.Get(s)
		if err != nil {
			parser.error(fmt.Sprintf("%s undeclared", s))
		}
		parser.move()
		if parser.look.GetTokenTag() != '[' {
			return id
		} else {
			return parser.offset(id)
		}
	default:
		parser.error("syntax error")
		return nil
	}
}

func (parser *Parser) offset(id *Id) Accesser {
	var w, t1, t2, loc Exprer
	typed := id.t
	parser.match('[')
	i := parser.bool()
	parser.match(']')
	typed = typed.(lexer.Arrayer).GetType()
	w = NewConstantInt(typed.GetWidth())
	t1 = NewArith(lexer.NewToken('*'), i, w)
	loc = t1

	for parser.look.GetTokenTag() == '[' {
		parser.match('[')
		i = parser.bool()
		parser.match(']')
		typed = typed.(lexer.Arrayer).GetType()
		w = NewConstantInt(typed.GetWidth())
		t1 = NewArith(lexer.NewToken('*'), i, w)
		t2 = NewArith(lexer.NewToken('+'), loc, t1)
		loc = t2
	}
	return NewAccess(id, loc, typed)
}
