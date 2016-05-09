package main

import "./parser"
import "./ast"

var ReduceMap map[string]func(stack []parse.Elem) interface{} = map[string]func(stack []parse.Elem) interface{}{
	"binop ~> expr": func(stack []parse.Elem) interface{} {
		return ast.BinOp{
			stack[len(stack)-3].Tree.(ast.Expr),
			stack[len(stack)-2].Tree.(ast.Operator),
			stack[len(stack)-1].Tree.(ast.Expr),
		}
	},

	"top": func(stack []parse.Elem) interface{} {
		return stack[len(stack)-1].Tree
	},

	"op:plus": func(stack []parse.Elem) interface{} {
		return ast.PLUS
	},

	"op:minus": func(stack []parse.Elem) interface{} {
		return ast.MINUS
	},
	"op:multiply": func(stack []parse.Elem) interface{} {
		return ast.MULTIPLY
	},
	"op:div": func(stack []parse.Elem) interface{} {
		return ast.DIV
	},

	"top-1": func(stack []parse.Elem) interface{} {
		return stack[len(stack)-2].Tree
	},

	"variable": func(stack []parse.Elem) interface{} {
		varname := stack[len(stack)-1].Token.Text
		return ast.Variable{varname}
	},

	"func": func(stack []parse.Elem) interface{} {
		var stmts []ast.Stmt = nil
		if stack[len(stack)-2].Tree != nil {
			stmts = stack[len(stack)-2].Tree.([]ast.Stmt)
		}
		return ast.Function{
			stack[len(stack)-7].Token.Text,            // Name
			stack[len(stack)-5].Tree.([]ast.Variable), // TODO: Params
			stmts,
		}
	},

	"func_list_create": func(stack []parse.Elem) interface{} {
		return make([]interface{}, 0)
	},

	"single_func": func(stack []parse.Elem) interface{} {
		return []ast.Function{stack[len(stack)-1].Tree.(ast.Function)}
	},

	"append_func": func(stack []parse.Elem) interface{} {
		return append(stack[len(stack)-2].Tree.([]ast.Function),
			stack[len(stack)-1].Tree.(ast.Function))
	},

	"append_stmt": func(stack []parse.Elem) interface{} {

		stmts := []ast.Stmt{stack[len(stack)-2].Tree.(ast.Stmt)}
		if stack[len(stack)-1].Tree != nil {
			stmts = append(stmts, stack[len(stack)-1].Tree.([]ast.Stmt)...)
		}

		return stmts
	},

	"empty_stmt": func(stack []parse.Elem) interface{} {
		return ast.ExprStmt{}
	},

	"expr_stmt": func(stack []parse.Elem) interface{} {
		return ast.ExprStmt{stack[len(stack)-2].Tree.(ast.Expr)}
	},

	"variable_decl": func(stack []parse.Elem) interface{} {
		return ast.Variable{"decl: " + stack[len(stack)-2].Token.Text}
	},

	"formal_param_prepend": func(stack []parse.Elem) interface{} {
		l := []ast.Variable{ast.Variable{stack[len(stack)-3].Token.Text}}
		return append(l, stack[len(stack)-1].Tree.([]ast.Variable)...)
	},

	"formal_param_single": func(stack []parse.Elem) interface{} {
		return []ast.Variable{ast.Variable{stack[len(stack)-1].Token.Text}}
	},

	"noop": func(stack []parse.Elem) interface{} { return nil },
}
