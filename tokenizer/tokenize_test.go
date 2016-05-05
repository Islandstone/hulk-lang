package tokenizer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenizeSimpleExpr(t *testing.T) {
	assert := assert.New(t)

	in := bytes.NewBufferString("func shazbot ** a + b    - \t c*d / fooBarBaz")

	tokenizer := NewTokenizer(in)

	assert.Equal(Token{FUNCTION, "func"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "shazbot"},
		tokenizer.GetNextToken())

	assert.Equal(Token{DOUBLESTAR, "**"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "a"},
		tokenizer.GetNextToken())

	assert.Equal(Token{PLUS, "+"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "b"},
		tokenizer.GetNextToken())

	assert.Equal(Token{MINUS, "-"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "c"},
		tokenizer.GetNextToken())

	assert.Equal(Token{STAR, "*"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "d"},
		tokenizer.GetNextToken())

	assert.Equal(Token{DIV, "/"},
		tokenizer.GetNextToken())

	assert.Equal(Token{IDENTIFIER, "fooBarBaz"},
		tokenizer.GetNextToken())

	assert.Equal(Token{EOF, "$"},
		tokenizer.GetNextToken())
}
