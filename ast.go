package main

type NodeType string

const (
	// STATEMENTS
	NodeTypeProgram             = "Program"
	NodeTypeVariableDeclaration = "VariableDeclaration"

	// EXPRESSIONS
	NodeTypeNumericLiteral  = "NumericLiteral"
	NodeTypeIdentifier      = "Identifier"
	NodeTypeBinaryExpession = "BinaryExpession"
)

type Stmter interface {
	Kind() NodeType
}

type Stmt struct {
	kind NodeType
}

func (s *Stmt) Kind() NodeType {
	return s.kind
}

type Program struct {
	*Stmt
	body []Stmter
}

type VariableDeclaration struct {
	*Stmt
	constant   bool
	identifier string
	value      Stmter
}

type Expr struct {
	Stmt
}

type BinaryExpession struct {
	*Stmt
	left     Stmter
	right    Stmter
	operator string
}

type Identifier struct {
	*Stmt
	symbol string
}

type NumericLiteral struct {
	*Stmt
	value float64
}
