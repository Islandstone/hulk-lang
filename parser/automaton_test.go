package parse

import (
	tok "../tokenizer"
	// "fmt"
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

	a.Init(prods)

	assert.Len(t, a.items, 26)
}

func TestGenerateEpsilonItems(t *testing.T) {
	prods := []Production{
		// Functions
		NewProduction("S'", []ProductionElement{
			NewNonTerminal("Func_list"),
		}, nil),
		NewProduction("Func_list", []ProductionElement{
			NewNonTerminal("Func_nelist"),
		}, nil),
		NewProduction("Func_list", []ProductionElement{
			NewTerminal(tok.EPSILON),
		}, nil),
		NewProduction("Func_nelist", []ProductionElement{
			NewNonTerminal("Func_nelist"),
			NewNonTerminal("Func"),
		}, nil),
		NewProduction("Func_nelist", []ProductionElement{
			NewNonTerminal("Func"),
		}, nil),
		NewProduction("Func", []ProductionElement{
			NewTerminal(tok.FUNCTION),
			NewTerminal(tok.IDENTIFIER),
			NewTerminal(tok.LPAR),
			NewTerminal(tok.RPAR),
			NewTerminal(tok.LBRACE),
			NewTerminal(tok.RBRACE),
		}, nil),
	}

	a := &Automaton{}
	a.Init(prods)

	/*
		a.BuildState([]int{0})
		for i, state := range a.states {
			fmt.Printf("State %d:\n", i)

			for _, itemId := range state.itemIds {
				fmt.Println(a.items[itemId])
			}

			fmt.Println()
		}

		assert.Len(t, a.states, 11)
	*/
}

func TestMakeStartState(t *testing.T) {
	t.Skip()

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

	a.Init(prods)
	a.BuildState([]int{0})

	/*
		for i, state := range a.states {
			fmt.Printf("State %d:\n", i)

			for _, itemId := range state.itemIds {
				fmt.Println(a.items[itemId])
			}

			fmt.Printf("neigh: %v\n", state.neigh)

			fmt.Println("")
		}

		tbl, gotoTbl := a.BuildTable()
		for stateId, line := range tbl {
			fmt.Printf("State %d: ", stateId)
			fmt.Printf("%+v\n", line)
			fmt.Println("goto:", gotoTbl[stateId])
		}
	*/
}

// TODO: Items with <EPSILON> are not handled correctly
func TestListProductions(t *testing.T) {
	// t.Skip()
	prods := []Production{
		// Functions
		NewProduction("S'", []ProductionElement{
			NewNonTerminal("Func_nelist"),
		}, nil),
		NewProduction("Func_list", []ProductionElement{
			NewNonTerminal("Func_nelist"),
		}, nil),
		NewProduction("Func_list", []ProductionElement{
			NewTerminal(tok.EPSILON),
		}, nil),
		NewProduction("Func_nelist", []ProductionElement{
			NewNonTerminal("Func_nelist"),
			NewNonTerminal("Func"),
		}, nil),
		NewProduction("Func_nelist", []ProductionElement{
			NewNonTerminal("Func"),
		}, nil),
		NewProduction("Func", []ProductionElement{
			NewTerminal(tok.FUNCTION),
			NewTerminal(tok.IDENTIFIER),
			NewTerminal(tok.LPAR),
			NewTerminal(tok.RPAR),
			NewTerminal(tok.LBRACE),
			NewNonTerminal("Stmt_list"),
			NewTerminal(tok.RBRACE),
		}, nil),

		NewProduction("Stmt_list", []ProductionElement{
			NewNonTerminal("Stmt_nelist"),
		}, nil),
		NewProduction("Stmt_list", []ProductionElement{
			NewTerminal(tok.EPSILON),
		}, nil),
		NewProduction("Stmt_nelist", []ProductionElement{
			NewNonTerminal("Stmt_nelist"),
			NewNonTerminal("Stmt"),
		}, nil),
		NewProduction("Stmt_nelist", []ProductionElement{
			NewNonTerminal("Stmt"),
		}, nil),

		NewProduction("Stmt", []ProductionElement{
			NewTerminal(tok.PLUS),
		}, nil),
	}

	a := &Automaton{}

	a.Init(prods)
	a.BuildState([]int{0})

	/*
		for i, state := range a.states {
			fmt.Printf("State %d:\n", i)

			for _, itemId := range state.itemIds {
				fmt.Println(a.items[itemId])
			}

			// fmt.Printf("neigh: %v\n", state.neigh)
			fmt.Printf("neigh: ")
			for _, neigh := range state.neigh {
				fmt.Printf("%d ", neigh.id)
			}

			fmt.Println()
			fmt.Println()
		}

		tbl, _ := a.BuildTable()
		for i, _ := range a.states {
			state := tbl[i]
			fmt.Printf("State %d: ", i)
			fmt.Printf("%+v\n\n", state)
			// fmt.Println("goto:", gotoTbl[stateId])
		}
	*/
}
