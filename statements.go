package main

func (i *Interpreter) formatError(e *CustomError, p int) *CustomError {
	return e.addTrace(p)
}

func (i *Interpreter) evalProgram(program *Program, env *Environments) (RuntimeVal, *CustomError) {
	var lastEvaulatedValue RuntimeVal
	lastEvaulatedValue = makeNull()

	for _, statements := range program.body {
		evaluated, err := i.evaluate(statements, env)
		if err != nil {
			return nil, i.formatError(err, program.Pos())
		}

		lastEvaulatedValue = evaluated
	}

	return lastEvaulatedValue, nil
}

func (i *Interpreter) evalVarDeclaration(declaration *VariableDeclaration, env *Environments) (RuntimeVal, *CustomError) {
	if declaration.value == nil {
		return env.declareVar(declaration.identifier, makeNull(), declaration.constant)
	} else {
		value, err := i.evaluate(declaration.value, env)
		if err != nil {
			return nil, i.formatError(err, declaration.Pos())
		}
		return env.declareVar(declaration.identifier, value, declaration.constant)
	}
}

func (i *Interpreter) evalConditionDeclaration(cnd *ConditionDeclaration, env *Environments) (RuntimeVal, *CustomError) {
	lhs, err := i.evaluate(cnd.left, env)
	if err != nil {
		return nil, i.formatError(err, cnd.Pos())
	}
	rhs, err := i.evaluate(cnd.right, env)
	if err != nil {
		return nil, i.formatError(err, cnd.Pos())
	}

	lhsVal, okLhs := lhs.(*NumberVal)
	rhsVal, okRhs := rhs.(*NumberVal)
	if okLhs && okRhs {
		return i.evalNumericConditionExpr(*lhsVal, *rhsVal, cnd.operator)
	}

	lsVal, okLs := lhs.(*StringVal)
	rsVal, okRs := rhs.(*StringVal)
	if okLs && okRs {
		return i.evalStringConditionExpr(*lsVal, *rsVal, cnd.operator)
	}

	lsBVal, okBLs := lhs.(*BoolVal)
	rsBVal, okBRs := rhs.(*BoolVal)
	if okBLs && okBRs {
		return i.evalBoolConditionExpr(*lsBVal, *rsBVal, cnd.operator)
	}

	return makeNull(), nil
}

func (i *Interpreter) evalFunctionDeclaration(declaration *FunctionDeclaration, env *Environments) (RuntimeVal, *CustomError) {
	fn := &FnValue{
		Type:           ValueFunction,
		name:           declaration.name,
		declarationEnv: env,
		paramaters:     declaration.parameters,
		body:           declaration.body,
	}

	return env.declareVar(declaration.name, fn, true)
}
