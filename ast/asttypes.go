package ast

type ASTnode interface {
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
