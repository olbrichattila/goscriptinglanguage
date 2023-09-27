package main

import "fmt"

type Interpreter struct {
}

func newInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) evaluate(astNode Stmter, env *Environments) (RuntimeVal, *CustomError) {
	kind := astNode.Kind()
	switch kind {
	case NodeTypeNumericLiteral:
		return makeNumber(astNode.(*NumericLiteral).value), nil
	case NodeTypeStringLIteral:
		return makeString(astNode.(*StringLiteral).value), nil
	case NodeTypeBinaryExpession:
		return i.evalBinaryExpression(astNode.(*BinaryExpession), env)
	case NodeTypeProgram:
		return i.evalProgram(astNode.(*Program), env)
	case NodeTypeIdentifier:
		return i.evalIdentifier(astNode.(*Identifier), env)
	case NodeTypeVariableDeclaration:
		return i.evalVarDeclaration(astNode.(*VariableDeclaration), env)
	case NodeTypeAssigmentExpression:
		return i.evalAssignment(astNode.(*AssignmentExpr), env)
	case NodeTypeObjectLiteral:
		return i.evalObjectExpr(astNode.(*ObjectLiteral), env)
	case NodeTypeCallExpression:
		return i.evalCallExpr(astNode.(*CallExpression), env)
	case NodeTypeFunctionDeclaration:
		return i.evalFunctionDeclaration(astNode.(*FunctionDeclaration), env)
	case NodeTypeConditionExpression:
		return i.evalConditionDeclaration(astNode.(*ConditionDeclaration), env)
	case NodeTypeIfExpression:
		return i.evalIfExpr(astNode.(*IfExpression), env)
	case NodeTypeForExpression:
		return i.evalForExpr(astNode.(*ForExpression), env)
	case NodeTypeBreakExpression:
		return i.evalBreakExpr(astNode.(*BreakExpression), env)
	case NodeTypeContinueExpression:
		return i.evalContinueExpr(astNode.(*ContinueExpression), env)
	default:
		return nil, newCustomError(fmt.Sprintf("This AST node has not yet been setup for interpretation %s", kind)).addTrace(astNode.Pos())
	}
}
