package parse

import (
	tok "../tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowStart(t *testing.T) {
	follow := Follow([]Production{
		// S' -> S
		NewProduction("S'", []ProductionElement{
			ProductionElement{false, "S", tok.UNKNOWN},
		}, nil),
		// S -> a
		NewProduction("S", []ProductionElement{
			ProductionElement{true, "", tok.IDENTIFIER},
		}, nil),
	}, map[string][]tok.Terminal{
		"S'": []tok.Terminal{tok.IDENTIFIER},
		"S":  []tok.Terminal{tok.IDENTIFIER},
	})

	contains := func(set interface{}, elem tok.Terminal) {
		assert.Contains(t, set, elem)
	}

	// assert.Contains(t, follow["S"], "$", "Missing $")
	contains(follow["S'"], tok.EOF)
	contains(follow["S"], tok.EOF)
}

/*
func TestFollow(t *testing.T) {
	prods := []Production{
		NewProduction("S' -> E"),
		NewProduction("E -> E A T"),
		NewProduction("E -> T"),
		NewProduction("A -> +"),
		NewProduction("A -> -"),
		NewProduction("T -> T M F"),
		NewProduction("T -> F"),
		NewProduction("M -> *"),
		NewProduction("F -> ( E )"),
		NewProduction("F -> num"),
	}

	Follow(prods, First(prods))
}
*/

func TestFollowNullableSuffix(t *testing.T) {
	prods := []Production{
		// S' -> E
		NewProduction("S'", []ProductionElement{
			ProductionElement{false, "E", tok.UNKNOWN},
		}, nil),
		// E -> E A T
		NewProduction("E", []ProductionElement{
			ProductionElement{false, "E", tok.UNKNOWN},
			ProductionElement{false, "A", tok.UNKNOWN},
			ProductionElement{false, "T", tok.UNKNOWN},
		}, nil),
		// E -> y
		NewProduction("E", []ProductionElement{
			ProductionElement{true, "", tok.EQUALS},
		}, nil),
		// T -> epsilon
		NewProduction("T", []ProductionElement{
			ProductionElement{true, "", tok.EPSILON},
		}, nil),
		// A -> x
		NewProduction("A", []ProductionElement{
			ProductionElement{true, "", tok.STAR},
		}, nil),
	}

	follow := Follow(prods, First(prods))
	contains := func(set interface{}, elem tok.Terminal) {
		assert.Contains(t, set, elem)
	}

	contains(follow["E"], tok.STAR)
}
