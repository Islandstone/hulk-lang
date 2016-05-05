package tokenizer

import (
	"bufio"
	"io"
	"regexp"
	"unicode"
)

type Terminal int

const (
	UNKNOWN Terminal = iota
	EPSILON

	LPAR
	RPAR
	LBRACE
	RBRACE
	EQUALS
	DOUBLEEQUALS
	PLUS
	MINUS
	STAR
	DOUBLESTAR
	DIV
	DOT
	IDENTIFIER
	LEFTARROW
	FUNCTION
	EOF
)

func (t Terminal) String() string {
	switch t {
	case UNKNOWN:
		return "UNKNOWN"
	case EPSILON:
		return "EPSILON"
	case LPAR:
		return "LPAR"
	case RPAR:
		return "RPAR"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case EQUALS:
		return "EQUALS"
	case DOUBLEEQUALS:
		return "DOUBLEEQUALS"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case DOUBLESTAR:
		return "DOUBLESTAR"
	case DIV:
		return "DIV"
	case DOT:
		return "DOT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case LEFTARROW:
		return "LEFTARROW"
	case FUNCTION:
		return "FUNCTION"
	case EOF:
		return "EOF"
	default:
		return "UNNAMED TERMINAL"
	}
}

var typeTable map[string]Terminal = map[string]Terminal{
	"(":    LPAR,
	")":    RPAR,
	"{":    LBRACE,
	"}":    RBRACE,
	"=":    EQUALS,
	"==":   DOUBLEEQUALS,
	"+":    PLUS,
	"-":    MINUS,
	"*":    STAR,
	"**":   DOUBLESTAR,
	"/":    DIV,
	".":    DOT,
	"<-":   LEFTARROW,
	"func": FUNCTION,
}

var (
	identifiers = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")
)

type Token struct {
	Type Terminal
	Text string
}

type Tokenizer struct {
	reader *bufio.Reader
}

func NewTokenizer(input io.Reader) Tokenizer {
	return Tokenizer{bufio.NewReader(input)}
}

func isspace(r rune) bool {
	if r == ' ' {
		return true
	} else if r == '\n' {
		return true
	}

	return false
}

func (tokenizer *Tokenizer) GetNextToken() (t Token) {
	t.Type = UNKNOWN
	t.Text = ""

	for {
		c, length, err := tokenizer.reader.ReadRune()

		// fmt.Printf("Current: %#v - Read \"%c\", len: %v, err: %v\n", t, c, length, err)

		if length == 0 {
			// fmt.Printf("%+v\n", t)
			if t.Type != UNKNOWN {
				return
			} else {
				t = Token{EOF, "$"}
				break
			}
		}

		if err != nil {
			panic(err)
		}

		if unicode.IsSpace(c) {
			if t.Type == UNKNOWN {
				continue
			} else {
				break
			}
		}

		if identifiers.MatchString(t.Text + string(c)) {
			t.Type = IDENTIFIER
			t.Text += string(c)
		} else if t.Type == IDENTIFIER {
			tokenizer.reader.UnreadRune()
			break
		}

		if v, ok := typeTable[t.Text+string(c)]; ok {
			t.Type = v
			t.Text += string(c)
		} else if v, ok = typeTable[t.Text]; ok {
			t.Type = v

			// The last 'c' of "func" was added above by the identifier code
			// Do not unread in this case
			if v != FUNCTION {
				tokenizer.reader.UnreadRune()
			}

			break
		}
	}

	// fmt.Println()
	return
}
