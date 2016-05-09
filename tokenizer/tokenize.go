package tokenizer

import (
	"bufio"
	// "fmt"
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
	LESS

	PLUS
	MINUS
	STAR
	DOUBLESTAR
	DIV
	DOT
	SEMICOLON
	COMMA
	VAR
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
	case SEMICOLON:
		return "SEMICOLON"
	case EOF:
		return "EOF"
	case VAR:
		return "VAR"
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
	"<":    LESS,
	"<-":   LEFTARROW,
	";":    SEMICOLON,
	"func": FUNCTION,
	"var":  VAR,
	",":    COMMA,
}

var TerminalReflection map[string]Terminal = map[string]Terminal{
	"EPSILON":      EPSILON,
	"LPAR":         LPAR,
	"RPAR":         RPAR,
	"LBRACE":       LBRACE,
	"RBRACE":       RBRACE,
	"EQUALS":       EQUALS,
	"DOUBLEEQUALS": DOUBLEEQUALS,
	"PLUS":         PLUS,
	"MINUS":        MINUS,
	"STAR":         STAR,
	"DOUBLESTAR":   DOUBLESTAR,
	"DIV":          DIV,
	"DOT":          DOT,
	"SEMICOLON":    SEMICOLON,
	"VAR":          VAR,
	"IDENTIFIER":   IDENTIFIER,
	"LEFTARROW":    LEFTARROW,
	"FUNCTION":     FUNCTION,
	"EOF":          EOF,
	"COMMA":        COMMA,
}

var (
	identifiers = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
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
	switch r {
	case ' ':
		return true
	case '\n':
		return true
	default:
		return false
	}
}

func (tokenizer *Tokenizer) GetNextToken() (t Token) {
	t.Type = UNKNOWN
	t.Text = ""

	// TODO: /* */ comments
	// TODO: Nested /* */ comments
	// commentLevel := 0

	insideComment := false

	// fmt.Println("Restart")
	for {
		c, length, err := tokenizer.reader.ReadRune()

		// fmt.Printf("Current: %#v - Read \"%c\", len: %v, err: %v\n", t, c, length, err)

		if length == 0 {
			// fmt.Printf("%+v\n", t)
			if t.Type != UNKNOWN {
				// fmt.Println("b")
				return
			} else {
				t = Token{EOF, "$"}
				// fmt.Println("e")
				break
			}
		}

		if err != nil {
			panic(err)
		}

		if insideComment {
			if c == '\n' {
				insideComment = false

				t.Type = UNKNOWN
				t.Text = ""
			} else {

			}
			continue
		}

		if unicode.IsSpace(c) {
			if t.Type == UNKNOWN {
				// fmt.Println("a")
				continue
			} else {
				// fmt.Println("b")
				break
			}
		}

		if t.Text == "/" && c == '/' {
			insideComment = true
			continue
		}

		if identifiers.MatchString(t.Text + string(c)) {
			t.Type = IDENTIFIER
			t.Text += string(c)
			// fmt.Println("e")
		} else if t.Type == IDENTIFIER {
			tokenizer.reader.UnreadRune()
			// fmt.Println("c")
			break
		}

		if v, ok := typeTable[t.Text+string(c)]; ok {
			t.Type = v
			t.Text += string(c)
			// fmt.Println("w")
		} else if v, ok = typeTable[t.Text]; ok {
			t.Type = v

			// The last 'c' of "func" was added above by the identifier code
			// Do not unread in this case
			if v != FUNCTION && v != VAR {
				// fmt.Println("c")
				tokenizer.reader.UnreadRune()
			}

			// fmt.Println("d")
			break
		}
	}

	// fmt.Println()
	return
}
