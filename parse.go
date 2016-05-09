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

	prods := parse.ReadGrammar("grammar.yaml", ReduceMap)

	a := &parse.Automaton{}

	a.Init(prods)
	a.BuildState([]int{0})
	table, gototable = a.BuildTable()
}

func Parse(input io.Reader) ast.ASTnode {
	stack := make([]parse.Elem, 1)

	stack[0] = parse.Elem{initial, tok.Token{tok.UNKNOWN, "top"}, nil}

	tknzr := tok.NewTokenizer(input)

	token := tknzr.GetNextToken()

	for {
		state := stack[len(stack)-1]

		action := table[state.State][token.Type]

		if s, ok := action.(parse.Shift); ok {
			// fmt.Println("Shifting", token)
			newState := parse.Elem{s.Goto, token, nil}
			stack = append(stack, newState)

			token = tknzr.GetNextToken()

			continue
		} else if r, ok := action.(parse.Reduce); ok {
			fmt.Println("Reducing with rule", r.Rule)

			var astNode ast.ASTnode = nil
			if r.Create != nil {
				astNode = r.Create(stack)
			} else {
				fmt.Println("Note: Missing create function for", r.Rule)
			}

			stack = stack[:len(stack)-r.Count]
			curState := stack[len(stack)-1].State
			// fmt.Println("New state after reduction: ", gototable[curState][r.Nonterminal])
			newState := parse.Elem{gototable[curState][r.Nonterminal], tok.Token{tok.UNKNOWN, r.Nonterminal}, astNode}
			stack = append(stack, newState)
		} else if _, ok := action.(parse.Accept); ok {
			fmt.Println("Accepted")
			return stack[1].Tree
		} else {
			fmt.Println(stack)
			fmt.Printf("Error: Found no action for token %s in state %d\n", token, state.State)
			return nil
		}
	}
}

func main() {
	// r := bytes.NewBufferString("a + b ; c - shazbot;")

	/*
		for i, _ := range table {
			state := table[i]
			fmt.Printf("State %d: ", i)
			fmt.Printf("%+v\n\n", state)
			// fmt.Println("goto:", gotoTbl[stateId])
		}
	*/

	// r := bytes.NewBufferString("func foobar(){} func shazbot(){}")

	t := Parse(bytes.NewBufferString("func foobar(param1, param2){ var a; CR3 <- asdf; } "))
	fmt.Printf("%#v\n", t)
}
