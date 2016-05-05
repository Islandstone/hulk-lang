package parse

import (
	"../ast"
	"../tokenizer"
)

type ProductionElement struct {
	Terminal bool

	// name for non-term
	Name string

	// Token type for terminals
	Token tokenizer.Terminal
}

type Production struct {
	left  ProductionElement
	right []ProductionElement
	// astType interface{}
	// astType reflect.Type
	Create func([]Elem) interface{}
}

type Elem struct {
	State int
	Token tokenizer.Token
	Tree  ast.ASTnode
}

const (
	EPSILON = "ε"
)

func NewProduction(name string, rhs []ProductionElement, create func([]Elem) interface{}) Production {
	return Production{
		ProductionElement{false, name, tokenizer.UNKNOWN},
		rhs,
		create,
	}
}

func NewTerminal(token tokenizer.Terminal) (p ProductionElement) {
	p.Terminal = true
	p.Name = ""
	p.Token = token

	return
}

func NewNonTerminal(name string) (p ProductionElement) {
	p.Terminal = true
	p.Name = name
	p.Token = tokenizer.UNKNOWN

	return
}
