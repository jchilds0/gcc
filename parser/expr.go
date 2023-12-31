package parser

import (
	"fmt"
	"gcc/lexer"
)

type Exprer interface {
	Gen() Exprer
	Reduce(Exprer) Exprer
	String() string
	Type() lexer.Typer
	Jumping(t, f int)
}

type Expr struct {
	Node
	op lexer.Tokener
	t  lexer.Typer
}

func NewExpr(tok lexer.Tokener, p lexer.Typer) *Expr {
	return &Expr{Node: *NewNode(), op: tok, t: p}
}

func (expr *Expr) Gen() Exprer {
	return expr
}

func (_ *Expr) Reduce(expr Exprer) Exprer {
	return expr
}

func (expr *Expr) Jumping(t, f int) {
	expr.emitJumps("", t, f)
}

func (expr *Expr) emitJumps(test string, t int, f int) {
	if t != 0 && f != 0 {
		expr.Emit(fmt.Sprintf("if %s goto L%d", test, t))
		expr.Emit(fmt.Sprintf("goto L%d", f))
	} else if t != 0 {
		expr.Emit(fmt.Sprintf("if %s goto L%d", test, t))
	} else if f != 0 {
		expr.Emit(fmt.Sprintf("iffalse %s goto L%d", test, f))
	}
}

func (expr *Expr) String() string {
	return expr.op.String()
}

func (expr *Expr) Type() lexer.Typer {
	return expr.t
}

type Id struct {
	Expr
	offset int
}

func NewId(id *lexer.Word, p lexer.Typer, b int) *Id {
	return &Id{Expr: *NewExpr(id, p), offset: b}
}

var TempCount = 0

type Temp struct {
	Expr
	number int
}

func NewTemp(p lexer.Typer) *Temp {
	temp := &Temp{Expr: *NewExpr(lexer.WordTemp, p)}
	TempCount++
	temp.number = TempCount
	return temp
}

func (temp *Temp) String() string {
	return fmt.Sprintf("t%d", temp.number)
}

type Constant struct {
	Expr
}

func NewConstant(tok lexer.Tokener, p *lexer.Type) *Constant {
	return &Constant{Expr: *NewExpr(tok, p)}
}

func NewConstantInt(i int) *Constant {
	return &Constant{Expr: *NewExpr(lexer.NewNum(i), lexer.Int)}
}

var ConstantTrue = NewConstant(lexer.WordTrue, lexer.Bool)
var ConstantFalse = NewConstant(lexer.WordFalse, lexer.Bool)

func (cons *Constant) Jumping(t int, f int) {
	if cons == ConstantTrue && t != 0 {
		cons.Emit(fmt.Sprintf("goto L%d", t))
	} else if cons == ConstantFalse && f != 0 {
		cons.Emit(fmt.Sprintf("goto L%d", f))
	}
}
