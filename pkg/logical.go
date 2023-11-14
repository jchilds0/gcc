package pkg

import "fmt"

type Logical struct {
	Expr
	expr1 Exprer
	expr2 Exprer
}

func NewLogical(tok Tokener, x1 Exprer, x2 Exprer) *Logical {
	logical := &Logical{expr1: x1, expr2: x2}
	t := logical.check(x1.Type(), x2.Type())
	logical.Expr = *NewExpr(tok, t)

	if t == nil {
		logical.Error("type error")
	}
	return logical
}

func (_ *Logical) check(p1 Typer, p2 Typer) *Type {
	if p1.String() == Bool.String() && p2.String() == Bool.String() {
		return Bool
	}
	return nil
}

func (logical *Logical) Gen() Exprer {
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

type Or struct {
	Logical
}

func NewOr(tok Tokener, x1 Exprer, x2 Exprer) *Or {
	return &Or{Logical: *NewLogical(tok, x1, x2)}
}

func (or *Or) Jumping(t, f int) {
	var label int
	if t != 0 {
		label = t
	} else {
		label = or.NewLabel()
	}

	or.expr1.Jumping(label, 0)
	or.expr2.Jumping(t, f)
	if t == 0 {
		or.EmitLabel(label)
	}
}

type And struct {
	Logical
}

func NewAnd(tok Tokener, x1 Exprer, x2 Exprer) *And {
	return &And{Logical: *NewLogical(tok, x1, x2)}
}

func (and *And) Jumping(t, f int) {
	var label int
	if f != 0 {
		label = f
	} else {
		label = and.NewLabel()
	}

	and.expr1.Jumping(0, label)
	and.expr2.Jumping(t, f)
	if f == 0 {
		and.EmitLabel(label)
	}
}

type Not struct {
	Logical
}

func NewNot(tok Tokener, x1 Exprer) *Not {
	return &Not{Logical: *NewLogical(tok, x1, x1)}
}

func (not *Not) Jumping(t, f int) {
	not.expr2.Jumping(f, t)
}

func (not *Not) String() string {
	return fmt.Sprintf("%s %s", not.op.String(), not.expr2.String())
}

type Rel struct {
	Logical
}

func NewRel(tok Tokener, x1 Exprer, x2 Exprer) *Rel {
	rel := &Rel{Logical: Logical{expr1: x1, expr2: x2}}
	t := rel.check(x1.Type(), x2.Type())
	rel.Logical.Expr = *NewExpr(tok, t)

	if t == nil {
		rel.Error("type error")
	}
	return rel
}

func (rel *Rel) check(p1 Typer, p2 Typer) *Type {
	_, ok1 := p1.(Arrayer)
	_, ok2 := p2.(Arrayer)
	if ok1 || ok2 {
		return nil
	} else if p1.String() == p2.String() && p1.GetWidth() == p2.GetWidth() {
		return Bool
	} else {
		return nil
	}
}

func (rel *Rel) Jumping(t, f int) {
	a := rel.expr1.Reduce(rel.expr1)
	b := rel.expr2.Reduce(rel.expr2)
	test := fmt.Sprintf("%s %s %s", a.String(), rel.op.String(), b.String())

	rel.emitJumps(test, t, f)
}
