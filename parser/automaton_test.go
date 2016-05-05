package parse

import (
	tok "../tokenizer"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClosure(t *testing.T) {
	/*
		a := assert.New(t)
		state := State{[]Item{
			NewItem(NewProduction("S -> E")),
		},
			nil,
		}

		prods := []Production{
			NewProduction("S -> E"),
			NewProduction("E -> A"),
			NewProduction("A -> a"),
		}

		state.ApplyClosure(prods)

		a.Contains(state.items, Item{NewProduction("E -> A"), 0})
		a.Contains(state.items, Item{NewProduction("A -> a"), 0})
	*/
}

func TestGenerateItems(t *testing.T) {
	prods := []Production{
		// S' -> E
		NewProduction("S'", []ProductionElement{
			{false, "E", tok.UNKNOWN},
		}, nil),
		// E -> E A T
		NewProduction("E", []ProductionElement{
			{false, "E", tok.UNKNOWN},
			{false, "A", tok.UNKNOWN},
			{false, "T", tok.UNKNOWN},
		}, nil),
		// E -> T
		NewProduction("E", []ProductionElement{
			{false, "T", tok.UNKNOWN},
		}, nil),
		// A -> +
		NewProduction("A", []ProductionElement{
			{true, "", tok.PLUS},
		}, nil),
		// A -> -
		NewProduction("A", []ProductionElement{
			{true, "", tok.MINUS},
		}, nil),
		// T -> T M F
		NewProduction("T", []ProductionElement{
			{false, "T", tok.UNKNOWN},
			{false, "M", tok.UNKNOWN},
			{false, "F", tok.UNKNOWN},
		}, nil),
		// T -> F
		NewProduction("T", []ProductionElement{
			{false, "F", tok.UNKNOWN},
		}, nil),
		// M -> *
		NewProduction("M", []ProductionElement{
			{true, "", tok.STAR},
		}, nil),
		// F -> ( E )
		NewProduction("F", []ProductionElement{
			{true, "", tok.LPAR},
			{false, "", tok.UNKNOWN},
			{true, "", tok.RPAR},
		}, nil),
		// F -> num
		NewProduction("F", []ProductionElement{
			{true, "", tok.IDENTIFIER},
		}, nil),
	}

	a := &Automaton{}

	a.init(prods)

	assert.Len(t, a.items, 26)
}

func TestMakeStartState(t *testing.T) {
	prods := []Production{
		// S' -> E
		NewProduction("S'", []ProductionElement{
			{false, "E", tok.UNKNOWN},
		}, nil),
		// E -> E A T
		NewProduction("E", []ProductionElement{
			{false, "E", tok.UNKNOWN},
			{false, "A", tok.UNKNOWN},
			{false, "T", tok.UNKNOWN},
		}, nil),
		// E -> T
		NewProduction("E", []ProductionElement{
			{false, "T", tok.UNKNOWN},
		}, nil),
		// A -> +
		NewProduction("A", []ProductionElement{
			{true, "", tok.PLUS},
		}, nil),
		// A -> -
		NewProduction("A", []ProductionElement{
			{true, "", tok.MINUS},
		}, nil),
		// T -> T M F
		NewProduction("T", []ProductionElement{
			{false, "T", tok.UNKNOWN},
			{false, "M", tok.UNKNOWN},
			{false, "F", tok.UNKNOWN},
		}, nil),
		// T -> F
		NewProduction("T", []ProductionElement{
			{false, "F", tok.UNKNOWN},
		}, nil),
		// M -> *
		NewProduction("M", []ProductionElement{
			{true, "", tok.STAR},
		}, nil),
		// F -> ( E )
		NewProduction("F", []ProductionElement{
			{true, "", tok.LPAR},
			{false, "E", tok.UNKNOWN},
			{true, "", tok.RPAR},
		}, nil),
		// F -> num
		NewProduction("F", []ProductionElement{
			{true, "", tok.IDENTIFIER},
		}, nil),
	}

	a := &Automaton{}

	a.init(prods)
	a.buildState([]int{0})

	for i, state := range a.states {
		fmt.Printf("State %d:\n", i)

		for _, itemId := range state.itemIds {
			fmt.Println(a.items[itemId])
		}

		fmt.Printf("neigh: %v\n", state.neigh)

		fmt.Println("")
	}

	a.buildTable()
	tbl, gotoTbl := a.buildTable()
	for stateId, line := range tbl {
		fmt.Printf("State %d: ", stateId)
		fmt.Printf("%+v\n", line)
		fmt.Println("goto:", gotoTbl[stateId])
	}
}
