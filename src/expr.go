package dragonbook

import (
	"fmt"
)

type ExprInterface interface {
	Gen() *Expr
	Reduce() *Expr
}

type Expr struct {
	Node
	op *Token
	t  *Type
}

func NewExpr(tok *Token, p *Type) *Expr {
	return &Expr{op: tok, t: p}
}

func (expr *Expr) Gen() *Expr {
	return expr
}

func (expr *Expr) Reduce() *Expr {
	return expr
}

func (expr *Expr) Jumping(t int, f int) {
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

type Id struct {
	Expr
	offset int
}

func NewId(id *Word, p *Type, b int) *Id {
	return &Id{Expr: *NewExpr(&id.Token, p), offset: b}
}

type Op struct {
	Expr
}

func (op *Op) reduce() *Temp {
	x := op.Gen()
	t := NewTemp(op.t)
	op.Emit(t.String() + " = " + x.String())
	return t
}

type Arith struct {
	Op
	expr1 Expr
	expr2 Expr
}

func NewArith(tok *Token, x1 *Expr, x2 *Expr) *Arith {
	arith := new(Arith)
	t := max(x1.t, x2.t)

	if t == nil {
		arith.Error("type error")
	}

	arith.op = tok
	arith.t = t
	arith.expr1 = *x1
	arith.expr2 = *x2
	return arith
}

func (arith *Arith) Gen() *Arith {
	return NewArith(arith.op, arith.expr1.Reduce(), arith.expr2.Reduce())
}

func (arith *Arith) String() string {
	return fmt.Sprintf("%s %s %s", arith.expr1.String(), arith.op.String(), arith.expr2.String())
}

type Temp struct {
	Expr
	count  int
	number int
}

func NewTemp(p *Type) *Temp {
	return &Temp{Expr: *NewExpr(&WordTemp.Token, p), count: 0, number: 1}
}

type Unary struct {
	Op
	expr Expr
}

func NewUnary(tok *Token, x *Expr) *Unary {
	t := max(Int, x.t)
	if t == nil {
		x.Error("type error")
	}

	return &Unary{Op: Op{Expr: *NewExpr(tok, t)}, expr: *x}
}

func (unary *Unary) Gen() *Unary {
	return NewUnary(unary.op, unary.expr.Reduce())
}

func (unary *Unary) String() string {
	return fmt.Sprintf("%s %s", unary.op.String(), unary.expr.String())
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

type Logical struct {
	Expr
	expr1 Expr
	expr2 Expr
}

func check(p1 *Type, p2 *Type) *Type {
	if *p1 == *Bool && *p2 == *Bool {
		return Bool
	}
	return nil
}

func NewLogical(tok *Token, x1 *Expr, x2 *Expr) *Logical {
	t := check(x1.t, x2.t)
	logical := &Logical{Expr: *NewExpr(tok, t), expr1: *x1, expr2: *x2}

	if t == nil {
		logical.Error("type error")
	}
	return logical
}

func (logical *Logical) Gen() *Temp {
	f := logical.NewLabel()
	a := logical.NewLabel()
	temp := NewTemp(logical.t)

	logical.Jumping(0, f)
	logical.Emit(fmt.Sprintf("%s = true", temp.String()))
	logical.Emit(fmt.Sprintf("goto L%d", a))
	logical.EmitLabel(f)
	logical.Emit(fmt.Sprintf("%s = false", temp.String()))
	logical.EmitLabel(a)
	return temp
}

func (logical *Logical) String() string {
	return fmt.Sprintf("%s %s %s", logical.expr1.String(), logical.op.String(), logical.expr2.String())
}
