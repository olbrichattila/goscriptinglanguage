package main

import "fmt"

func (i *Interpreter) evalBinaryExpression(binop *BinaryExpession, env *Environments) (RuntimeVal, *CustomError) {
	lhs, err := i.evaluate(binop.left, env)
	if err != nil {
		return nil, i.formatError(err, binop.Pos())
	}
	rhs, err := i.evaluate(binop.right, env)
	if err != nil {
		return nil, i.formatError(err, binop.Pos())
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

func (i *Interpreter) evalIdentifier(ident *Identifier, env *Environments) (RuntimeVal, *CustomError) {
	val, err := env.lookupVar(ident.symbol)
	if err != nil {
		return nil, i.formatError(err, ident.Pos())
	}

	return val, nil
}

func (i *Interpreter) evalNumericBinaryExpr(lhs, rhs NumberVal, operator string) (*NumberVal, *CustomError) {
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
			return nil, newCustomError("Division by 0")
		}
		result = lhs.Value / rhs.Value
	case "%":
		result = float64(int(lhs.Value) % int(rhs.Value))
	default:
		return nil, newCustomError(fmt.Sprintf("Operator %s not implemented", operator))
	}

	return makeNumber(result), nil
}

func (i *Interpreter) evalStringBinaryExpr(lhs, rhs StringVal, operator string) (*StringVal, *CustomError) {
	var result string
	switch operator {
	case "+":
		result = lhs.Value + rhs.Value
	default:
		return nil, newCustomError(fmt.Sprintf("Operator %s not implemented", operator))
	}

	return makeString(result), nil
}

func (i *Interpreter) evalAssignment(node *AssignmentExpr, env *Environments) (RuntimeVal, *CustomError) {
	if node.assigne.Kind() != NodeTypeIdentifier {
		return nil, newCustomError("Invalid LHS iside assignment expression")
	}

	varname := node.assigne.(*Identifier).symbol
	evaulated, err := i.evaluate(node.value, env)
	if err != nil {
		return nil, i.formatError(err, node.Pos())
	}

	result, err := env.assignVar(varname, evaulated)
	if err != nil {
		err.addTrace(node.Pos())
		err.addTrace(node.assigne.Pos())
	}

	return result, err
}

func (i *Interpreter) evalObjectExpr(node *ObjectLiteral, env *Environments) (RuntimeVal, *CustomError) {
	object := &ObjectVal{Type: ValueObject, properties: make(map[string]RuntimeVal)}

	for _, property := range node.properties {
		if property.value == nil {
			runtimeVal, err := env.lookupVar(property.key)
			if err != nil {
				return nil, i.formatError(err, node.Pos())
			}
			object.properties[property.key] = runtimeVal
		} else {

			runtimeVal, err := i.evaluate(property.value, env)
			if err != nil {
				return nil, i.formatError(err, node.Pos())
			}
			object.properties[property.key] = runtimeVal
		}
	}

	return object, nil
}

func (i *Interpreter) evalCallExpr(expr *CallExpression, env *Environments) (RuntimeVal, *CustomError) {
	var args []RuntimeVal

	for _, arg := range expr.args {
		ev, err := i.evaluate(*arg, env)
		if err != nil {
			return nil, i.formatError(err, expr.Pos())
		}
		args = append(args, ev)
	}

	f, err := i.evaluate(expr.caller, env)
	if err != nil {
		return nil, i.formatError(err, expr.Pos())
	}

	if fn, ok := f.(*NativeFnValue); ok {
		return fn.call(args, env), nil
	}

	if fnc, ok := f.(*FnValue); ok {
		scope, err := newEnvironments(fnc.declarationEnv)
		if err != nil {
			return nil, i.formatError(err, expr.Pos())
		}

		for ind, varName := range fnc.paramaters {
			// @TODO check the bouds here, verify the airity of the function
			_, err := scope.declareVar(varName, args[ind], false)
			if err != nil {
				return nil, i.formatError(err, expr.Pos())
			}
		}

		var result RuntimeVal = makeNull()
		for _, statement := range fnc.body {
			result, err = i.evaluate(statement, scope)
			if err != nil {
				return nil, i.formatError(err, expr.Pos())
			}
		}

		return result, nil
	}

	return nil, newCustomError("cannot call value which is not a function").addTrace(expr.Pos())
}

func (i *Interpreter) evalNumericConditionExpr(lhs, rhs NumberVal, operator string) (*BoolVal, *CustomError) {
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
		return nil, newCustomError(fmt.Sprintf("Conditional Operator %s not implemented", operator))
	}

	return makeBool(result), nil
}

func (i *Interpreter) evalStringConditionExpr(lhs, rhs StringVal, operator string) (*BoolVal, *CustomError) {
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
		return nil, newCustomError(fmt.Sprintf("String Conditional Operator %s not implemented", operator))
	}

	return makeBool(result), nil
}

func (i *Interpreter) evalBoolConditionExpr(lhs, rhs BoolVal, operator string) (*BoolVal, *CustomError) {
	var result bool
	switch operator {
	case "&":
		result = lhs.Value && rhs.Value
	case "|":
		result = lhs.Value || rhs.Value
	default:
		return nil, newCustomError(fmt.Sprintf("Logical conditional operator %s not implemented", operator))
	}

	return makeBool(result), nil
}

func (i *Interpreter) evalIfExpr(ifE *IfExpression, env *Environments) (RuntimeVal, *CustomError) {
	var cond RuntimeVal
	var err *CustomError
	if ifE.condition == nil {
		cond = makeBool(true)
	} else {
		cond, err = i.evaluate(ifE.condition, env)
		if err != nil {
			return nil, i.formatError(err, ifE.Pos())
		}
	}

	var result RuntimeVal = makeNull()
	if cond.(*BoolVal).Value == true {
		for _, statement := range ifE.body {
			result, err = i.evaluate(statement, env)
			if err != nil {
				return nil, i.formatError(err, ifE.Pos())
			}
		}
	} else if ifE.elseExpression != nil {
		return i.evalIfExpr(ifE.elseExpression.(*IfExpression), env)
	}

	return result, nil
}

func (i *Interpreter) evalForExpr(forE *ForExpression, env *Environments) (RuntimeVal, *CustomError) {
	var err *CustomError
	var result RuntimeVal = makeNull()
	braked := false
	continued := false

	if forE.declaration != nil {
		_, err := i.evaluate(forE.declaration, env)
		if err != nil {
			return nil, i.formatError(err, forE.Pos())
		}
	}

	for {
		if forE.condition != nil {
			cond, err := i.evaluate(forE.condition, env)
			if err != nil {
				return nil, i.formatError(err, forE.Pos())
			}

			if cond.(*BoolVal).Value == false {
				break
			}
		}

		for _, statement := range forE.body {
			result, err = i.evaluate(statement, env)
			if _, ok := result.(*BreakVal); ok {
				braked = true
				break
			}

			if _, ok := result.(*ContinueVal); ok {
				continued = true
				break
			}
			if err != nil {
				return nil, i.formatError(err, forE.Pos())
			}
		}

		if braked == true {
			break
		}

		if continued == true {
			continue
		}

		if forE.afterCondition != nil {
			cond, err := i.evaluate(forE.afterCondition, env)
			if err != nil {
				return nil, i.formatError(err, forE.Pos())
			}

			if cond.(*BoolVal).Value == false {
				break
			}
		}

		if forE.incrementalExpression != nil {
			_, err = i.evaluate(forE.incrementalExpression, env)
			if err != nil {
				return nil, i.formatError(err, forE.Pos())
			}
		}
	}

	return result, nil
}

func (i *Interpreter) evalBreakExpr(forE *BreakExpression, env *Environments) (RuntimeVal, *CustomError) {
	return makeBreak(), nil
}

func (i *Interpreter) evalContinueExpr(forE *ContinueExpression, env *Environments) (RuntimeVal, *CustomError) {
	return makeContinue(), nil
}

func (i *Interpreter) evalSwitchExpr(sw *SwitchExpression, env *Environments) (RuntimeVal, *CustomError) {
	// @Todo refactor this, too complex
	cv, err := i.evaluate(sw.value, env)
	if err != nil {
		err.addTrace(sw.Pos())
		return nil, err
	}

	for _, swcase := range sw.body {

		if _, ok := swcase.compare.(*BoolVal); ok {
			isBreak, err := i.evalBody(swcase.body, env)
			if err != nil {
				return nil, err
			}
			if isBreak {
				break
			}
			continue
		}
		if a, ok := cv.(*NumberVal); ok {
			if b, ok := swcase.compare.(*NumberVal); ok {
				if a.Value == b.Value {
					isBreak, err := i.evalBody(swcase.body, env)
					if err != nil {
						return nil, err
					}
					if isBreak {
						break
					}
				}
				continue
			}
		}

		if a, ok := cv.(*StringVal); ok {
			if b, ok := swcase.compare.(*StringVal); ok {
				if a.Value == b.Value {
					isBreak, err := i.evalBody(swcase.body, env)
					if err != nil {
						return nil, err
					}
					if isBreak {
						break
					}
				}
				continue
			}
		}

		runError := newCustomError(fmt.Sprintf("Type error in switch - case %T, %T", cv, swcase.compare)).addTrace(sw.Pos())
		runError.addTrace(swcase.pos)

		return nil, runError
	}
	return makeNull(), nil
}

func (i *Interpreter) evalBody(items []Stmter, env *Environments) (bool, *CustomError) {
	var lastRValue RuntimeVal
	var err *CustomError
	for _, item := range items {
		lastRValue, err = i.evaluate(item, env)
		if err != nil {
			return false, err
		}
		if _, ok := lastRValue.(*BreakVal); ok {
			return true, nil
		}
	}

	return false, nil

}
