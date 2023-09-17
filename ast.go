package main

type NodeType string

const (
	// STATEMENTS
	NodeTypeProgram             = "Program"
	NodeTypeVariableDeclaration = "VariableDeclaration"

	// EXPRESSIONS
	NodeTypeBinaryExpession     = "BinaryExpession"
	NodeTypeAssigmentExpression = "AssignmentExpr"
	NodeTypeMemberExpression    = "MemberExpression"
	NodeTypeCallExpression      = "CallExpression"

	// Literals
	NodeTypeProperty       = "Property"
	NodeTypeObjectLiteral  = "ObjectLiteral"
	NodeTypeNumericLiteral = "NumericLiteral"
	NodeTypeIdentifier     = "Identifier"
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

type AssignmentExpr struct {
	*Stmt
	assigne Stmter
	value   Stmter
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

type Property struct {
	*Stmt
	key   string
	value Stmter
}

type ObjectLiteral struct {
	*Stmt
	properties []*Property
}

type CallExpression struct {
	*Stmt
	args   []*Stmter
	caller Stmter
}

type MemberExpression struct {
	*Stmt
	object   Stmter
	propert  Stmter
	computed bool
}
