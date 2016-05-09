package parse

import (
	tok "../tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PE ProductionElement

func TestNewProduction(t *testing.T) {
	// A -> a B

	p := NewProduction("A", []ProductionElement{
		ProductionElement{true, "", tok.IDENTIFIER},
		ProductionElement{false, "B", tok.UNKNOWN},
	}, nil)

	assert.Equal(t, p.Left, ProductionElement{false, "A", tok.UNKNOWN})

	// assert.Equal(t, p.Right, []string{"a", "B"})
	assert.Equal(t, p.Right, []ProductionElement{
		ProductionElement{true, "", tok.IDENTIFIER},
		ProductionElement{false, "B", tok.UNKNOWN},
	})
}

func TestRuleOne(t *testing.T) {
	p := NewProduction("A", []ProductionElement{
		ProductionElement{true, "a", tok.IDENTIFIER},
		ProductionElement{false, "B", tok.UNKNOWN},
	}, nil)

	first := First([]Production{p})
	assert.Equal(t, []tok.Terminal{tok.IDENTIFIER}, first["A"])
}

func TestRuleTwo(t *testing.T) {
	first := First([]Production{
		// NewProduction("A -> B a", nil, nil),
		NewProduction("A", []ProductionElement{
			ProductionElement{false, "B", tok.UNKNOWN},
			ProductionElement{true, "a", tok.IDENTIFIER},
		}, nil),
		// NewProduction("B -> b", nil, nil),
		NewProduction("B", []ProductionElement{
			ProductionElement{true, "b", tok.EQUALS},
		}, nil),
	})

	assert.Equal(t, []tok.Terminal{tok.EQUALS}, first["A"])
	assert.Equal(t, []tok.Terminal{tok.EQUALS}, first["B"])
}

func TestRuleThree(t *testing.T) {
	first := First([]Production{
		// NewProduction("A -> a", nil, nil),
		NewProduction("A", []ProductionElement{
			ProductionElement{true, "a", tok.IDENTIFIER},
		}, nil),
		// NewProduction("B -> C A", nil, nil),
		NewProduction("B", []ProductionElement{
			ProductionElement{false, "C", tok.UNKNOWN},
			ProductionElement{false, "A", tok.UNKNOWN},
		}, nil),
		// NewProduction("C -> ε", nil, nil),
		NewProduction("C", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
	})

	// assert.Equal(t, first["A"], []string{"a"})
	assert.Equal(t, []tok.Terminal{tok.IDENTIFIER}, first["A"])
	// assert.Equal(t, first["B"], []string{"a"})
	assert.Equal(t, []tok.Terminal{tok.IDENTIFIER}, first["B"])
	// assert.Equal(t, first["C"], []string{EPSILON})
	assert.Equal(t, []tok.Terminal{tok.EPSILON}, first["C"])
}

func TestNullability(t *testing.T) {
	first := First([]Production{
		// A -> ε
		NewProduction("A", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
		// B -> C A
		NewProduction("B", []ProductionElement{
			ProductionElement{false, "C", tok.UNKNOWN},
			ProductionElement{false, "A", tok.UNKNOWN},
		}, nil),
		// C -> ε
		NewProduction("C", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
	})

	assert.Equal(t, []tok.Terminal{tok.EPSILON}, first["A"])
	assert.Equal(t, []tok.Terminal{tok.EPSILON}, first["B"])
	assert.Equal(t, []tok.Terminal{tok.EPSILON}, first["C"])
}

func Test1(t *testing.T) {
	first := First([]Production{
		// A -> a
		NewProduction("A", []ProductionElement{
			ProductionElement{true, "", tok.IDENTIFIER},
		}, nil),
		// A -> ε
		NewProduction("A", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
		// B -> C A
		NewProduction("B", []ProductionElement{
			ProductionElement{false, "C", tok.UNKNOWN},
			ProductionElement{false, "A", tok.UNKNOWN},
		}, nil),
		// C -> ε
		NewProduction("C", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
	})

	assert.Len(t, first["A"], 2)
	assert.Contains(t, first["A"], tok.IDENTIFIER, "Missing <identifier> in First(A)")
	assert.Contains(t, first["A"], tok.EPSILON, "Missing "+EPSILON+" in First(A)")

	assert.Len(t, first["B"], 2)
	assert.Contains(t, first["B"], tok.IDENTIFIER, "Missing <identifier> in First(B)")
	assert.Contains(t, first["B"], tok.EPSILON, "Missing "+EPSILON+" in First(B)")

	assert.Equal(t, []tok.Terminal{tok.EPSILON}, first["C"])
}

/*
func TestPage172(t *testing.T) {
	// first := First([]Production{
	// S -> I
	// S -> o [other]
	// I -> i ( E ) S L [if-statement]
	// L -> ε
	// L -> e S
	// E -> 0
	// E -> 1
	// })

	// assert.Contains(t, first["S"], "i", "Missing i")
}
*/

/*
func TestPage170(t *testing.T) {
	first := First([]Production{
		NewProduction("E -> E A T", nil, nil),
		NewProduction("E -> T", nil, nil),
		NewProduction("A -> +", nil, nil),
		NewProduction("A -> -", nil, nil),
		NewProduction("T -> T M F", nil, nil),
		NewProduction("T -> F", nil, nil),
		NewProduction("M -> *", nil, nil),
		NewProduction("F -> ( E )", nil, nil),
		NewProduction("F -> num", nil, nil),
	})

	assert.Contains(t, first["E"], "(", "Missing (")
	assert.Contains(t, first["E"], "num", "Missing num")
}
*/
