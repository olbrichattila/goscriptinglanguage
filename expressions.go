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

	lsVal, okLs := lhs.(*StringVal)
	rsVal, okRs := rhs.(*StringVal)
	if okLs && okRs {
		return i.evalStringBinaryExpr(*lsVal, *rsVal, binop.operator)
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

func (i *Interpreter) evalStringBinaryExpr(lhs, rhs StringVal, operator string) (*StringVal, error) {
	var result string
	switch operator {
	case "+":
		result = lhs.Value + rhs.Value
	default:
		return nil, fmt.Errorf("Operator %s not implemented", operator)
	}

	return makeString(result), nil
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

func (i *Interpreter) evalObjectExpr(node *ObjectLiteral, env *Environments) (RuntimeVal, error) {
	object := &ObjectVal{Type: ValueObject, properties: make(map[string]RuntimeVal)}

	for _, property := range node.properties {
		if property.value == nil {
			runtimeVal, err := env.lookupVar(property.key)
			if err != nil {
				return nil, err
			}
			object.properties[property.key] = runtimeVal
		} else {

			runtimeVal, err := i.evaluate(property.value, env)
			if err != nil {
				return nil, err
			}
			object.properties[property.key] = runtimeVal
		}
	}

	return object, nil
}

func (i *Interpreter) evalCallExpr(expr *CallExpression, env *Environments) (RuntimeVal, error) {
	var args []RuntimeVal

	for _, arg := range expr.args {
		ev, err := i.evaluate(*arg, env)
		if err != nil {
			return nil, err
		}
		args = append(args, ev)
	}

	f, err := i.evaluate(expr.caller, env)
	if err != nil {
		return nil, err
	}

	if fn, ok := f.(*NativeFnValue); ok {
		return fn.call(args, env), nil
	}

	if fnc, ok := f.(*FnValue); ok {
		scope, err := newEnvironments(fnc.declarationEnv)
		if err != nil {
			return nil, err
		}

		for i, varName := range fnc.paramaters {
			// @TODO check the bouds here, verify the airity of the function
			_, err := scope.declareVar(varName, args[i], false)
			if err != nil {
				return nil, err
			}
		}

		var result RuntimeVal = makeNull()
		for _, statement := range fnc.body {
			result, err = i.evaluate(statement, scope)
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("cannot call value which is not a function")
}

func (i *Interpreter) evalNumericConditionExpr(lhs, rhs NumberVal, operator string) (*BoolVal, error) {
	var result bool
	switch operator {
	case "=":
		result = lhs.Value == rhs.Value
	case ">":
		result = lhs.Value > rhs.Value
	case ">=":
		result = lhs.Value >= rhs.Value
	case "<":
		result = lhs.Value < rhs.Value
	case "<=":
		result = lhs.Value <= rhs.Value
	case "!=":
		result = lhs.Value != rhs.Value
	default:
		return nil, fmt.Errorf("Conditional Operator %s not implemented", operator)
	}

	return makeBool(result), nil
}

func (i *Interpreter) evalStringConditionExpr(lhs, rhs StringVal, operator string) (*BoolVal, error) {
	var result bool
	switch operator {
	case "=":
		result = lhs.Value == rhs.Value
	case ">":
		result = lhs.Value > rhs.Value
	case ">=":
		result = lhs.Value >= rhs.Value
	case "<":
		result = lhs.Value < rhs.Value
	case "<=":
		result = lhs.Value <= rhs.Value
	case "!=":
		result = lhs.Value != rhs.Value
	default:
		return nil, fmt.Errorf("String Conditional Operator %s not implemented", operator)
	}

	return makeBool(result), nil
}

func (i *Interpreter) evalIfExpr(ifE *IfExpression, env *Environments) (RuntimeVal, error) {
	cond, err := i.evaluate(ifE.condition, env)
	if err != nil {
		return nil, err
	}

	var result RuntimeVal = makeNull()
	if cond.(*BoolVal).Value == true {
		for _, statement := range ifE.body {
			result, err = i.evaluate(statement, env)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (i *Interpreter) evalForExpr(forE *ForExpression, env *Environments) (RuntimeVal, error) {

	var result RuntimeVal = makeNull()

	_, err := i.evaluate(forE.declaration, env)
	if err != nil {
		return nil, err
	}

	for {
		cond, err := i.evaluate(forE.condition, env)
		if err != nil {
			return nil, err
		}

		if cond.(*BoolVal).Value == false {
			break
		}

		_, err = i.evaluate(forE.incrementalExpression, env)

		for _, statement := range forE.body {
			result, err = i.evaluate(statement, env)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
