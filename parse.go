package main

import (
	"./ast"
	"./parser"
	tok "./tokenizer"
	"bytes"
	"fmt"
	"io"
)

// http://www.cs.cornell.edu/Courses/cs412/2006sp/lectures/lec09.pdf
// https://www.cs.colostate.edu/~cs453/yr2014/Slides/07-shift-reduce.ppt.pdf

var table map[int]map[tok.Terminal]interface{}
var gototable map[int]map[string]int
var initial int = 0

func init() {

	prods := []parse.Production{
		// S' -> E
		parse.NewProduction("S'", []parse.ProductionElement{
			{false, "E", tok.UNKNOWN},
		}, nil),
		// E -> E A T
		parse.NewProduction("E", []parse.ProductionElement{
			{false, "E", tok.UNKNOWN},
			{false, "A", tok.UNKNOWN},
			{false, "T", tok.UNKNOWN},
		}, func(stack []parse.Elem) interface{} {
			fmt.Printf("%#v ", stack[len(stack)-3].Tree)
			fmt.Printf("%#v ", stack[len(stack)-2].Tree)
			fmt.Printf("%#v ", stack[len(stack)-1].Tree)
			fmt.Println()

			t := ast.BinOp{
				stack[len(stack)-3].Tree.(ast.Expr),
				stack[len(stack)-2].Tree.(ast.Operator),
				stack[len(stack)-1].Tree.(ast.Expr),
			}

			return t
		}),
		// E -> T
		parse.NewProduction("E", []parse.ProductionElement{
			{false, "T", tok.UNKNOWN},
		}, func(stack []parse.Elem) interface{} {
			return stack[len(stack)-1].Tree
		}),
		// A -> +
		parse.NewProduction("A", []parse.ProductionElement{
			{true, "", tok.PLUS},
		}, func(stack []parse.Elem) interface{} {
			return ast.PLUS
		}),
		// A -> -
		parse.NewProduction("A", []parse.ProductionElement{
			{true, "", tok.MINUS},
		}, func(stack []parse.Elem) interface{} {
			return ast.MINUS
		}),
		// T -> T M F
		parse.NewProduction("T", []parse.ProductionElement{
			{false, "T", tok.UNKNOWN},
			{false, "M", tok.UNKNOWN},
			{false, "F", tok.UNKNOWN},
		}, func(stack []parse.Elem) interface{} {
			fmt.Printf("%#v ", stack[len(stack)-3].Tree)
			fmt.Printf("%#v ", stack[len(stack)-2].Tree)
			fmt.Printf("%#v ", stack[len(stack)-1].Tree)
			fmt.Println()

			t := ast.BinOp{
				stack[len(stack)-3].Tree.(ast.Expr),
				stack[len(stack)-2].Tree.(ast.Operator),
				stack[len(stack)-1].Tree.(ast.Expr),
			}

			return t
		}),
		// T -> F
		parse.NewProduction("T", []parse.ProductionElement{
			{false, "F", tok.UNKNOWN},
		}, func(stack []parse.Elem) interface{} {
			return stack[len(stack)-1].Tree
		}),
		// M -> *
		parse.NewProduction("M", []parse.ProductionElement{
			{true, "", tok.STAR},
		}, func(stack []parse.Elem) interface{} {
			return ast.MULTIPLY
		}),
		// F -> ( E )
		parse.NewProduction("F", []parse.ProductionElement{
			{true, "", tok.LPAR},
			{false, "E", tok.UNKNOWN},
			{true, "", tok.RPAR},
		}, func(stack []parse.Elem) interface{} {
			return stack[len(stack)-2].Tree
		}),
		// F -> num
		parse.NewProduction("F", []parse.ProductionElement{
			{true, "", tok.IDENTIFIER},
		}, func(stack []parse.Elem) interface{} {
			varname := stack[len(stack)-1].Token.Text
			return ast.Variable{varname}
		}),
	}

	a := &parse.Automaton{}

	a.Init(prods)
	a.BuildState([]int{0})
	table, gototable = a.BuildTable()
}

func Parse(input io.Reader) {
	stack := make([]parse.Elem, 1)

	stack[0] = parse.Elem{initial, tok.Token{tok.UNKNOWN, "top"}, nil}

	tknzr := tok.NewTokenizer(input)

	token := tknzr.GetNextToken()

	for {
		state := stack[len(stack)-1]

		action := table[state.State][token.Type]

		fmt.Println(stack)

		if s, ok := action.(parse.Shift); ok {
			fmt.Println("Shifting", token)
			newState := parse.Elem{s.Goto, token, nil}
			stack = append(stack, newState)

			token = tknzr.GetNextToken()

			continue
		} else if r, ok := action.(parse.Reduce); ok {
			fmt.Println("Reducing with rule", r.Rule)

			var astNode ast.ASTnode = nil
			if r.Create != nil {
				astNode = r.Create(stack)
			}

			stack = stack[:len(stack)-r.Count]
			curState := stack[len(stack)-1].State
			fmt.Println("New state after reduction: ", gototable[curState][r.Nonterminal])
			newState := parse.Elem{gototable[curState][r.Nonterminal], tok.Token{tok.UNKNOWN, r.Nonterminal}, astNode}
			stack = append(stack, newState)
		} else if _, ok := action.(parse.Accept); ok {
			fmt.Println("Accepted")
			return
		} else {
			fmt.Printf("Error: Found no action for token %s in state %d\n", token, state.State)
			return
		}
	}
}

func main() {
	r := bytes.NewBufferString("a +b * c - shazbot")
	Parse(r)
}
