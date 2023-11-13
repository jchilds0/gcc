package pkg

import (
	"fmt"
)

type Exprer interface {
	Gen() Exprer
	Reduce() Exprer
	String() string
	Type() Typer
	Jumping(t, f int)
}

type Expr struct {
	Node
	op *Token
	t  Typer
}

func NewExpr(tok *Token, p Typer) *Expr {
	return &Expr{Node: *NewNode(), op: tok, t: p}
}

func (expr *Expr) Gen() Exprer {
	return expr
}

func (expr *Expr) Reduce() Exprer {
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
		expr.Emit(fmt.Sprintf("iffalse %s goto L%d", test, t))
	}
}

func (expr *Expr) String() string {
	return expr.op.String()
}

func (expr *Expr) Type() Typer {
	return expr.t
}

type Id struct {
	Expr
	offset int
}

func NewId(id *Word, p Typer, b int) *Id {
	return &Id{Expr: *NewExpr(&id.Token, p), offset: b}
}

type Temp struct {
	Expr
	count  int
	number int
}

func NewTemp(p Typer) *Temp {
	return &Temp{Expr: *NewExpr(&WordTemp.Token, p), count: 0, number: 1}
}

type Constant struct {
	Expr
}

func NewConstant(tok *Token, p *Type) *Constant {
	return &Constant{Expr: *NewExpr(tok, p)}
}

func NewConstantInt(i int) *Constant {
	return &Constant{Expr: *NewExpr(&NewNum(i).Token, Int)}
}

var ConstantTrue = NewConstant(&WordTrue.Token, Bool)
var ConstantFalse = NewConstant(&WordFalse.Token, Bool)

func (cons *Constant) Jumping(t int, f int) {
	if cons == ConstantTrue && t != 0 {
		cons.Emit(fmt.Sprintf("goto L%d", t))
	} else if cons == ConstantFalse && f != 0 {
		cons.Emit(fmt.Sprintf("goto L%d", f))
	}
}
