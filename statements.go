package main

func (i *Interpreter) evalProgram(program *Program, env *Environments) (RuntimeVal, error) {
	var lastEvaulatedValue RuntimeVal
	lastEvaulatedValue = makeNull()

	for _, statements := range program.body {
		evaluated, err := i.evaluate(statements, env)
		if err != nil {
			return nil, err
		}

		lastEvaulatedValue = evaluated
	}

	return lastEvaulatedValue, nil
}

func (i *Interpreter) evalVarDeclaration(declaration *VariableDeclaration, env *Environments) (RuntimeVal, error) {
	if declaration.value == nil {
		return env.declareVar(declaration.identifier, makeNull(), declaration.constant)
	} else {
		value, err := i.evaluate(declaration.value, env)
		if err != nil {
			return nil, err
		}
		return env.declareVar(declaration.identifier, value, declaration.constant)
	}
}

func (i *Interpreter) evalFunctionDeclaration(declaration *FunctionDeclaration, env *Environments) (RuntimeVal, error) {
	fn := &FnValue{
		Type:           ValueFunction,
		name:           declaration.name,
		declarationEnv: env,
		paramaters:     declaration.parameters,
		body:           declaration.body,
	}

	return env.declareVar(declaration.name, fn, true)
}
