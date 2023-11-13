package pkg

import "strconv"

type Tokener interface {
	GetTokenTag() int
	String() string
}

type Token struct {
	tag int
}

func NewToken(t int) (token *Token) {
	return &Token{tag: t}
}

func (token *Token) GetTokenTag() int {
	return token.tag
}

func (token *Token) String() string {
	return string(token.tag)
}

type Num struct {
	Token
	value int
}

func NewNum(v int) *Num {
	return &Num{Token: Token{tag: NUM}, value: v}
}

func (num *Num) String() string {
	return string(num.value)
}

type Word struct {
	Token
	lexeme string
}

func NewWord(tag int, s string) (word *Word) {
	return &Word{Token: Token{tag: tag}, lexeme: s}
}

func (word *Word) String() string {
	return word.lexeme
}

var WordAnd = NewWord(AND, "&&")
var WordOr = NewWord(OR, "||")
var WordEq = NewWord(EQ, "==")
var WordNe = NewWord(NE, "!=")
var WordLe = NewWord(LE, "<=")
var WordGe = NewWord(GE, ">=")

var WordMinus = NewWord(MINUS, "minus")
var WordTrue = NewWord(TRUE, "true")
var WordFalse = NewWord(FALSE, "false")
var WordTemp = NewWord(TEMP, "t")

type Real struct {
	Token
	value float64
}

func NewReal(v float64) *Real {
	return &Real{Token: Token{tag: REAL}, value: v}
}

func (r *Real) String() string {
	return strconv.FormatFloat(r.value, 'f', -1, 64)
}

type Typer interface {
	Tokener
	GetWidth() int
	Numeric() bool
}

type Type struct {
	Word
	width int
}

func NewType(tag int, s string, w int) *Type {
	return &Type{Word: *NewWord(tag, s), width: w}
}

var Int = NewType(BASIC, "int", 4)
var Float = NewType(BASIC, "float", 8)
var Char = NewType(BASIC, "char", 1)
var Bool = NewType(BASIC, "bool", 1)

func (t *Type) Numeric() bool {
	Tint := t.lexeme == Int.lexeme && t.width == Int.width
	Tfloat := t.lexeme == Float.lexeme && t.width == Float.width
	Tchar := t.lexeme == Float.lexeme && t.width == Char.width

	return Tint || Tfloat || Tchar
}

func (t *Type) GetWidth() int {
	return t.width
}

func max(p1 Typer, p2 Typer) *Type {
	if !p1.Numeric() || !p2.Numeric() {
		return nil
	} else if p1 == Float || p2 == Float {
		return Float
	} else if p1 == Int || p2 == Int {
		return Int
	} else {
		return Char
	}
}

type Arrayer interface {
	Typer
	GetType() Typer
}

type Array struct {
	Type
	size int
	of   Typer
}

func NewArray(sz int, p Typer) *Array {
	return &Array{Type: *NewType(INDEX, "[]", sz*p.GetWidth()), size: sz, of: p}
}

func (array *Array) String() string {
	return "[" + string(array.size) + "] " + array.of.String()
}

func (array *Array) GetType() Typer {
	return array.of
}
