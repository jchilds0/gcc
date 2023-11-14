package internal

import "fmt"

type Op struct {
	Expr
}

func (op *Op) Reduce(expr Exprer) Exprer {
	x := expr.Gen()
	t := NewTemp(expr.Type())
	op.Emit(t.String() + " = " + x.String())
	return t
}

type Arith struct {
	Op
	expr1 Exprer
	expr2 Exprer
}

func NewArith(tok Tokener, x1 Exprer, x2 Exprer) *Arith {
	arith := new(Arith)
	t := max(x1.Type(), x2.Type())

	if t == nil {
		arith.Error("type error")
	}

	arith.op = tok
	arith.t = t
	arith.expr1 = x1
	arith.expr2 = x2
	return arith
}

func (arith *Arith) Gen() Exprer {
	return NewArith(arith.op, arith.expr1.Reduce(arith.expr1), arith.expr2.Reduce(arith.expr2))
}

func (arith *Arith) String() string {
	return fmt.Sprintf("%s %s %s", arith.expr1.String(), arith.op.String(), arith.expr2.String())
}

type Unary struct {
	Op
	expr Exprer
}

func NewUnary(tok Tokener, x Exprer) (un *Unary) {
	t := max(Int, x.Type())
	if t == nil {
		un.Error("type error")
	}
	un.Op = Op{Expr: *NewExpr(tok, t)}
	un.expr = x

	return
}

func (unary *Unary) Gen() Exprer {
	return NewUnary(unary.op, unary.expr.Reduce(unary.expr))
}

func (unary *Unary) String() string {
	return fmt.Sprintf("%s %s", unary.op.String(), unary.expr.String())
}

type Accesser interface {
	Exprer
	GetId() *Id
	GetIndex() Exprer
}

type Access struct {
	Op
	array Id
	index Exprer
}

func NewAccess(a *Id, i Exprer, p Typer) *Access {
	return &Access{
		Op:    Op{Expr: *NewExpr(NewWord(INDEX, "[]"), p)},
		array: *a,
		index: i,
	}
}

func (access *Access) Gen() Exprer {
	return NewAccess(&access.array, access.index.Reduce(access.index), access.t)
}

func (access *Access) Jumping(t, f int) {
	access.emitJumps(access.Reduce(access).String(), t, f)
}

func (access *Access) String() string {
	return fmt.Sprintf("%s [ %s ]", access.array.String(), access.index.String())
}

func (access *Access) GetId() *Id {
	return &access.array
}

func (access *Access) GetIndex() Exprer {
	return access.index
}
