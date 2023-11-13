package pkg

import (
	"bufio"
	"log"
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
	words  map[string]Tokener
}

func NewLexer(reader *bufio.Reader) (lexer *Lexer) {
	lexer = &Lexer{line: 1, peek: ' ', reader: reader, words: map[string]Tokener{}}

	// reserve words in the hash table
	lexer.reserve(WordTrue)
	lexer.reserve(WordFalse)

	lexer.reserve(Int)
	lexer.reserve(Float)
	lexer.reserve(Char)
	lexer.reserve(Bool)

	lexer.reserve(NewWord(IF, "if"))
	lexer.reserve(NewWord(ELSE, "else"))
	lexer.reserve(NewWord(WHILE, "while"))
	lexer.reserve(NewWord(DO, "do"))
	lexer.reserve(NewWord(BREAK, "break"))
	return
}

func (lexer *Lexer) reserve(t Tokener) {
	lexer.words[t.String()] = t
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
			return false
		}
		lexer.peek = ' '
	}

	return true
}

func (lexer *Lexer) Scan() Tokener {
WS:
	for {
		switch lexer.peek {
		case ' ', '\t', '\r':
		case '\n':
			lexer.line++
		default:
			break WS
		}

		lexer.readch()
	}

	switch lexer.peek {
	case '&':
		if lexer.readch('&') {
			return WordAnd
		} else {
			return NewToken('&')
		}
	case '|':
		if lexer.readch('|') {
			return WordOr
		} else {
			return NewToken('|')
		}
	case '=':
		if lexer.readch('=') {
			return WordEq
		} else {
			return NewToken('=')
		}
	case '!':
		if lexer.readch('=') {
			return WordNe
		} else {
			return NewToken('!')
		}
	case '<':
		if lexer.readch('=') {
			return WordLe
		} else {
			return NewToken('<')
		}
	case '>':
		if lexer.readch('=') {
			return WordGe
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
		w, ok := lexer.words[s]

		if ok == false {
			w = NewWord(ID, s)
			lexer.words[s] = w
		}

		return w
	}

	t := NewToken(int(lexer.peek))
	lexer.peek = ' '
	return t
}
