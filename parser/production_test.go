package parse

import (
	"bytes"
	_ "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadProduction(t *testing.T) {
	a := assert.New(t)

	/*
		s := `
		[[A]]
			create = "Foo"
			[[NonTerm]]
			name = "A"
			[Term]
			token = "a"
			[[NonTerm]]
			name = "B"
		`
	*/

	s := `
# Please	
production:
- name: A
  right: ["n:A", "t:IDENTIFIER"]
`

	prods := ReadGrammar(bytes.NewBufferString(s))
	a.Len(prods, 1)

	// fmt.Println(prods[0])
}
