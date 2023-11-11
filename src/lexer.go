package dragonbook

import (
	"bufio"
	"log"
	"os"
	"unicode"
)

const (
	AND = iota + 256
	BASIC
	BREAK
	DO
	ELSE
	EQ
	FALSE
	GE
	ID
	IF
	INDEX
	LE
	MINUS
	NE
	NUM
	OR
	REAL
	TEMP
	TRUE
	WHILE
)

type Lexer struct {
	line   int
	peek   rune
	reader *bufio.Reader
	words  map[string]*Word
}

func NewLexer() (lexer *Lexer) {
	lexer = new(Lexer)
	lexer.line = 1
	lexer.peek = ' '
	lexer.reader = bufio.NewReader(os.Stdin)

	// reserve words in the hash table
	lexer.reserve(NewWord(IF, "if"))
	lexer.reserve(NewWord(ELSE, "else"))
	lexer.reserve(NewWord(WHILE, "while"))
	lexer.reserve(NewWord(DO, "do"))
	lexer.reserve(NewWord(BREAK, "break"))

	lexer.reserve(True)
	lexer.reserve(False)

	lexer.reserve(&Int.Word)
	lexer.reserve(&Float.Word)
	lexer.reserve(&Char.Word)
	lexer.reserve(&Bool.Word)
	return
}

func (lexer *Lexer) reserve(t *Word) {
	lexer.words[t.lexeme] = t
}

func (lexer *Lexer) readch(b ...rune) bool {
	var err error
	if len(b) == 0 {
		lexer.peek, _, err = lexer.reader.ReadRune()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		_ = lexer.readch()
		if lexer.peek != b[0] {
			lexer.peek = ' '
			return false
		}
	}

	return true
}

func (lexer *Lexer) Scan() TokenInterface {
WS:
	for {
		lexer.readch()

		switch lexer.peek {
		case ' ':
		case '\t':
		case '\n':
			lexer.line++
		default:
			break WS
		}
	}

	switch lexer.peek {
	case '&':
		if lexer.readch('&') {
			return And
		} else {
			return NewToken('&')
		}
	case '|':
		if lexer.readch('|') {
			return Or
		} else {
			return NewToken('|')
		}
	case '=':
		if lexer.readch('=') {
			return Eq
		} else {
			return NewToken('=')
		}
	case '!':
		if lexer.readch('=') {
			return Ne
		} else {
			return NewToken('!')
		}
	case '<':
		if lexer.readch('=') {
			return Le
		} else {
			return NewToken('<')
		}
	case '>':
		if lexer.readch('=') {
			return Ge
		} else {
			return NewToken('>')
		}
	}

	//if lexer.peek == '/' {
	//	lexer.peek, _, err = reader.ReadRune()
	//
	//	if lexer.peek == '/' {
	//		// continue to end of line
	//		for {
	//			lexer.peek, _, err = reader.ReadRune()
	//
	//			if lexer.peek == '\n' {
	//				lexer.peek, _, err = reader.ReadRune()
	//				break
	//			}
	//		}
	//	} else if lexer.peek == '*' {
	//		// continue until */
	//		star := false
	//		for {
	//			lexer.peek, _, err = reader.ReadRune()
	//
	//			if lexer.peek == '*' {
	//				star = true
	//			} else if lexer.peek == '/' && star {
	//				lexer.peek, _, err = reader.ReadRune()
	//				break
	//			} else {
	//				star = false
	//			}
	//		}
	//	} else {
	//		// didnt find a comment
	//		t := NewToken('/')
	//		return t, err
	//	}
	//}

	if unicode.IsDigit(lexer.peek) {
		v := int(lexer.peek - '0')
		lexer.readch()

		for unicode.IsDigit(lexer.peek) {
			v = 10*v + int(lexer.peek-'0')
			lexer.readch()
		}

		if lexer.peek != '.' {
			return NewNum(v)
		}
		d := 10.0
		x := float64(v)

		for unicode.IsDigit(lexer.peek) {
			if !unicode.IsDigit(lexer.peek) {
				break
			}

			x += float64(lexer.peek-'0') / d
			d *= 10
		}

		return NewReal(x)
	}

	if unicode.IsLetter(lexer.peek) {
		b := []rune{lexer.peek}
		lexer.readch()

		for unicode.IsLetter(lexer.peek) || unicode.IsDigit(lexer.peek) {
			b = append(b, lexer.peek)
			lexer.readch()
		}

		s := string(b)
		w := lexer.words[s]

		if w == nil {
			w = NewWord(ID, s)
			lexer.words[s] = w
		}

		return w
	}

	t := NewToken(int(lexer.peek))
	lexer.peek = ' '
	return t
}
