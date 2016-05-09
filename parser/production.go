package parse

import (
	"../ast"
	"../tokenizer"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type ProductionElement struct {
	Terminal bool

	// name for non-term
	Name string

	// Token type for terminals
	Token tokenizer.Terminal
}

type Production struct {
	Left   ProductionElement
	Right  []ProductionElement
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
	p.Terminal = false
	p.Name = name
	p.Token = tokenizer.UNKNOWN

	return
}

type Config struct {
	Production []ProductionConfig
}

type ProductionConfig struct {
	Name   string
	Right  []string
	Create string
}

type RHSspec struct {
	Terminal bool
	Name     string
}

func ReadGrammar(file string, reduceMap map[string]func([]Elem) interface{}) []Production {
	var conf Config

	b, _ := ioutil.ReadFile(file)
	err := yaml.Unmarshal(b, &conf)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("\n\n%#v\n", conf)

	prods := make([]Production, 0)

	for _, spec := range conf.Production {
		pes := make([]ProductionElement, 0)
		for _, pe := range spec.Right {
			if pe[1] != ':' {
				fmt.Println("Error: Expected ':' on the right-hand side")
				return nil
			}
			if pe[0] == 'n' {
				pes = append(pes, NewNonTerminal(pe[2:]))
			} else if pe[0] == 't' {
				if term, ok := tokenizer.TerminalReflection[pe[2:]]; !ok {
					fmt.Println("Error: Unknown ternimal named ", pe[2:])
					return nil
				} else {
					pes = append(pes, NewTerminal(term))
				}
			}
		}
		if spec.Create != "" {
			if _, ok := reduceMap[spec.Create]; !ok {
				fmt.Println("Unknown function", spec.Create)
			}
		}
		prods = append(prods, NewProduction(spec.Name, pes, reduceMap[spec.Create]))
	}

	return prods
}
