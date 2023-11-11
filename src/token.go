package dragonbook

import "strconv"

type TokenInterface interface {
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

var And = NewWord(AND, "&&")
var Or = NewWord(OR, "&&")
var Eq = NewWord(EQ, "&&")
var Ne = NewWord(NE, "&&")
var Le = NewWord(LE, "&&")
var Ge = NewWord(GE, "&&")

// var Minus = NewWord(MINUS, "&&")
var WordTrue = NewWord(TRUE, "&&")
var WordFalse = NewWord(FALSE, "&&")
var WordTemp = NewWord(TEMP, "&&")

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

func (t *Type) numeric() bool {
	return t == Int || t == Float || t == Char
}

func max(p1 *Type, p2 *Type) *Type {
	if !p1.numeric() || !p2.numeric() {
		return nil
	} else if p1 == Float || p2 == Float {
		return Float
	} else if p1 == Int || p2 == Int {
		return Int
	} else {
		return Char
	}
}

type Array struct {
	Type
	size int
	of   *Type
}

func NewArray(sz int, p *Type) *Array {
	return &Array{Type: *NewType(INDEX, "[]", sz*p.width), size: sz, of: p}
}

func (array *Array) String() string {
	return "[" + string(array.size) + "] " + array.of.String()
}
