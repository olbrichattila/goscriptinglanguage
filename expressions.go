package main

import "fmt"

func (i *Interpreter) evalBinaryExpression(binop *BinaryExpession, env *Environments) (RuntimeVal, error) {
	lhs, err := i.evaluate(binop.left, env)
	if err != nil {
		return nil, err
	}
	rhs, err := i.evaluate(binop.right, env)
	if err != nil {
		return nil, err
	}

	lhsVal, okLhs := lhs.(*NumberVal)
	rhsVal, okRhs := rhs.(*NumberVal)
	if okLhs && okRhs {
		return i.evalNumericBinaryExpr(*lhsVal, *rhsVal, binop.operator)
	}

	return makeNull(), nil
}

func (i *Interpreter) evalIdentifier(ident *Identifier, env *Environments) (RuntimeVal, error) {
	val, err := env.lookupVar(ident.symbol)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (i *Interpreter) evalNumericBinaryExpr(lhs, rhs NumberVal, operator string) (*NumberVal, error) {
	var result float64
	switch operator {
	case "+":
		result = lhs.Value + rhs.Value

	case "-":
		result = lhs.Value - rhs.Value
	case "*":
		result = lhs.Value * rhs.Value
	case "/":
		if rhs.Value == 0 {
			return nil, fmt.Errorf("Division by 0")
		}
		result = lhs.Value / rhs.Value
	case "%":
		result = float64(int(lhs.Value) % int(rhs.Value))
	default:
		return nil, fmt.Errorf("Operator %s not implemented", operator)
	}

	return makeNumber(result), nil
}

func (i *Interpreter) evalAssignment(node *AssignmentExpr, env *Environments) (RuntimeVal, error) {
	if node.assigne.Kind() != NodeTypeIdentifier {
		return nil, fmt.Errorf("Invalid LHS iside assignment expression")
	}

	varname := node.assigne.(*Identifier).symbol
	evaulated, err := i.evaluate(node.value, env)
	if err != nil {
		return nil, err
	}
	return env.assignVar(varname, evaulated)
}
