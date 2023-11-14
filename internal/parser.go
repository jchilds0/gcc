package internal

import (
	"bufio"
	"fmt"
	"log"
)

type Parser struct {
	lex  *Lexer
	look Tokener
	top  *Env
	used int
	w    *bufio.Writer
}

func NewParser(lexer *Lexer, write *bufio.Writer) *Parser {
	parser := &Parser{lex: lexer}
	parser.move()
	parser.top = NewEnv(nil)

	return parser
}

func (parser *Parser) move() {
	parser.look = parser.lex.Scan()
}

func (parser *Parser) error(s string) {
	log.Fatalf("near line %d: %s", parser.lex.line, s)
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
	for parser.look.GetTokenTag() == BASIC {
		p := parser.types()
		tok := parser.look
		parser.match(ID)
		parser.match(';')
		id := NewId(NewWord(tok.GetTokenTag(), tok.String()), p, parser.used)
		parser.top.Put(tok.String(), id)
		parser.used += p.GetWidth()
	}
}

func (parser *Parser) types() Typer {
	p := NewType(parser.look.GetTokenTag(), parser.look.String(), parser.look.(Typer).GetWidth())
	parser.match(BASIC)
	if parser.look.GetTokenTag() != '[' {
		return p
	} else {
		return parser.dims(p)
	}
}

func (parser *Parser) dims(p Typer) Typer {
	parser.match('[')
	parser.match(NUM)
	parser.match(']')

	if parser.look.GetTokenTag() == '[' {
		p = parser.dims(p)
	}

	return NewArray(0, p)
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
	case IF:
		parser.match(IF)
		parser.match('(')
		x = parser.bool()
		parser.match(')')
		s1 = parser.stmt()
		if parser.look.GetTokenTag() != ELSE {
			return NewIf(x, s1)
		} else {
			parser.match(ELSE)
			s2 = parser.stmt()
			return NewElse(x, s1, s2)
		}
	case WHILE:
		whilenode := new(While)
		savedStmt := StmtEnclosing
		StmtEnclosing = Stmter(whilenode)
		parser.match(WHILE)
		parser.match('(')
		x = parser.bool()
		parser.match(')')
		s1 = parser.stmt()
		whilenode.Init(x, s1)
		StmtEnclosing = savedStmt
		return whilenode
	case DO:
		donode := new(Do)
		savedStmt := StmtEnclosing
		StmtEnclosing = donode.stmt
		parser.match(DO)
		s1 := parser.stmt()
		parser.match(WHILE)
		parser.match('(')
		x := parser.bool()
		parser.match(')')
		parser.match(';')
		donode.Init(s1, x)
		StmtEnclosing = savedStmt
		return donode
	case BREAK:
		parser.match(BREAK)
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
	parser.match(ID)
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
	for parser.look.GetTokenTag() == OR {
		tok := parser.look
		parser.move()
		x = NewOr(tok, x, parser.join())
	}
	return x
}

func (parser *Parser) join() Exprer {
	x := parser.equality()
	for parser.look.GetTokenTag() == AND {
		tok := parser.look
		parser.move()
		x = NewAnd(tok, x, parser.equality())
	}
	return x
}

func (parser *Parser) equality() Exprer {
	x := parser.rel()
	for parser.look.GetTokenTag() == EQ || parser.look.GetTokenTag() == NE {
		tok := parser.look
		parser.move()
		x = NewRel(tok, x, parser.rel())
	}

	return x
}

func (parser *Parser) rel() Exprer {
	x := parser.expr()
	switch parser.look.GetTokenTag() {
	case '<', LE, GE, '>':
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
		return NewUnary(WordMinus, parser.unary())
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
	case NUM:
		x := NewConstant(parser.look, Int)
		parser.move()
		return x
	case REAL:
		x := NewConstant(parser.look, Float)
		parser.move()
		return x
	case TRUE:
		parser.move()
		return ConstantTrue
	case FALSE:
		parser.move()
		return ConstantFalse
	case ID:
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
	typed = typed.(Arrayer).GetType()
	w = NewConstantInt(typed.GetWidth())
	t1 = NewArith(NewToken('*'), i, w)
	loc = t1

	for parser.look.GetTokenTag() == '[' {
		parser.match('[')
		i = parser.bool()
		parser.match(']')
		typed = typed.(Arrayer).GetType()
		w = NewConstantInt(typed.GetWidth())
		t1 = NewArith(NewToken('*'), i, w)
		t2 = NewArith(NewToken('+'), loc, t1)
		loc = t2
	}
	return NewAccess(id, loc, typed)
}
