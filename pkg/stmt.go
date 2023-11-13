package pkg

import "fmt"

type Stmt struct {
	Node
	after int
}

type Stmter interface {
	Gen(b, a int)
	GetNode() *Node
}

var StmtNull Stmter
var StmtEnclosing Stmter

func (stmt *Stmt) Gen(b, a int) {}

func (stmt *Stmt) GetNode() *Node {
	return &stmt.Node
}

type If struct {
	Stmt
	expr Exprer
	stmt Stmter
}

func NewIf(x Exprer, s Stmter) *If {
	iff := &If{expr: x, stmt: s}
	if x.Type() != Bool {
		iff.Error("boolean is required in if")
	}

	return iff
}

func (iff *If) Gen(b, a int) {
	label := iff.NewLabel()
	iff.expr.Jumping(0, a)
	iff.EmitLabel(label)
	iff.stmt.Gen(label, a)
}

type Else struct {
	Stmt
	expr  Exprer
	stmt1 Stmter
	stmt2 Stmter
}

func NewElse(x Exprer, s1 Stmter, s2 Stmter) *Else {
	el := &Else{expr: x, stmt1: s1, stmt2: s2}
	if x.Type() == Bool {
		el.Error("boolean required in if")
	}

	return el
}

func (el *Else) Gen(b, a int) {
	label1 := el.NewLabel()
	label2 := el.NewLabel()

	el.expr.Jumping(0, label2)
	el.EmitLabel(label1)
	el.stmt1.Gen(label1, a)
	el.Emit(fmt.Sprintf("goto L%d", a))

	el.EmitLabel(label2)
	el.stmt2.Gen(label2, a)
}

type While struct {
	Stmt
	expr Exprer
	stmt Stmter
}

func (wh *While) Init(x Exprer, s Stmter) {
	if x.Type() != Bool {
		wh.Error("boolean required in while")
	}
	wh.expr = x
	wh.stmt = s
}

func (wh *While) Gen(b, a int) {
	wh.after = a
	wh.expr.Jumping(0, a)
	label := wh.NewLabel()
	wh.EmitLabel(label)
	wh.stmt.Gen(label, b)
	wh.Emit(fmt.Sprintf("goto L%d", b))
}

type Do struct {
	Stmt
	expr Exprer
	stmt Stmter
}

func (do *Do) Init(s Stmter, x Exprer) {
	if x.Type() != Bool {
		do.Error("boolean required in do")
	}

	do.expr = x
	do.stmt = s
}

func (do *Do) Gen(b, a int) {
	do.after = a
	label := do.NewLabel()
	do.stmt.Gen(b, label)
	do.EmitLabel(label)
	do.expr.Jumping(b, 0)
}

type Set struct {
	Stmt
	id   Id
	expr Exprer
}

func NewSet(i *Id, x Exprer) *Set {
	set := &Set{id: *i, expr: x}
	if set.check(i.Type(), x.Type()) == nil {
		set.Error("type error")
	}

	return set
}

func (_ *Set) check(p1 Typer, p2 Typer) Typer {
	if p1.Numeric() && p2.Numeric() {
		return p2
	} else if p1.String() == Bool.String() && p2.String() == Bool.String() {
		return p2
	} else {
		return nil
	}
}

func (set *Set) Gen(b, a int) {
	set.Emit(fmt.Sprintf("%s = %s", set.id.String(), set.expr.String()))
}

type SetElem struct {
	Stmt
	array *Id
	index Exprer
	expr  Exprer
}

func NewSetElem(x Accesser, y Exprer) *SetElem {
	setElem := &SetElem{array: x.GetId(), index: x.GetIndex(), expr: y}
	if setElem.check(x.Type(), y.Type()) == nil {
		setElem.Error("type error")
	}

	return setElem
}

func (_ *SetElem) check(p1 Typer, p2 Typer) Typer {
	_, ok1 := p1.(Arrayer)
	_, ok2 := p2.(Arrayer)
	if ok1 || ok2 {
		return nil
	} else if p1.String() == p2.String() {
		return p2
	} else if p1.Numeric() && p2.Numeric() {
		return p2
	} else {
		return nil
	}
}

func (setElem *SetElem) Gen(b, a int) {
	s1 := setElem.index.Reduce().String()
	s2 := setElem.expr.Reduce().String()

	setElem.Emit(fmt.Sprintf("%s [ %s ] = %s", setElem.array.String(), s1, s2))
}

type Seq struct {
	Stmt
	stmt1 Stmter
	stmt2 Stmter
}

func NewSeq(s1 Stmter, s2 Stmter) *Seq {
	return &Seq{stmt1: s1, stmt2: s2}
}

func (seq *Seq) Gen(b, a int) {
	if seq.stmt1 == nil {
		seq.stmt2.Gen(b, a)
	} else if seq.stmt2 == nil {
		seq.stmt1.Gen(b, a)
	} else {
		label := seq.NewLabel()
		seq.stmt1.Gen(b, label)
		seq.EmitLabel(label)
		seq.stmt2.Gen(label, a)
	}
}

type Break struct {
	Stmt
	stmt Stmter
}

func NewBreak() *Break {
	br := &Break{stmt: StmtEnclosing}
	if StmtEnclosing == nil {
		br.Error("unenclosed break")
	}
	return br
}

func (br *Break) Gen(b, a int) {
	br.Emit(fmt.Sprintf("goto L%d", br.after))
}
