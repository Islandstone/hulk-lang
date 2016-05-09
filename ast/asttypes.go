package ast

type ASTnode interface {
}

type Program struct {
	Functions []Function
}

type Expr interface {
}

type Const struct {
}

type Variable struct {
	Name string
}

type BinOp struct {
	Left  Expr
	Op    Operator
	Right Expr
}

type Stmt interface {
}

type RegisterAssignStmt struct {
	Register string

	Right Expr
}

type Function struct {
	Name   string
	Params []Variable
	Stmts  []Stmt
}

type Operator int

const (
	INVALID Operator = iota
	PLUS
	MINUS
	MULTIPLY
	DIV
	EXP
	DOT
)

type ExprStmt struct {
	Expr
}
