package main

import "fmt"

type Environments struct {
	parent    *Environments
	variables map[string]RuntimeVal
	constants map[string]interface{}
}

func newEnvironments(parent *Environments) *Environments {
	return &Environments{
		parent:    parent,
		variables: make(map[string]RuntimeVal),
		constants: make(map[string]interface{}),
	}
}

func (e *Environments) declareVar(varName string, value RuntimeVal, constant bool) (RuntimeVal, error) {
	_, exist := e.variables[varName]
	if exist {
		return nil, fmt.Errorf("Variable %s already exists", varName)
	}
	e.variables[varName] = value
	if constant {
		e.constants[varName] = nil
	}

	return value, nil
}

func (e *Environments) assignVar(varName string, value RuntimeVal) (RuntimeVal, error) {
	env, err := e.resolve(varName)
	if err != nil {
		return nil, err
	}

	_, exist := e.constants[varName]
	if exist {
		return nil, fmt.Errorf("Constant variable %s cannot be updated", varName)
	}

	env.variables[varName] = value

	return value, nil
}

func (e *Environments) resolve(varName string) (*Environments, error) {
	_, exist := e.variables[varName]
	if exist {
		return e, nil
	}
	if e.parent == nil {
		return nil, fmt.Errorf("Variable %s could not be resolved", varName)
	}

	return e.parent.resolve(varName)
}

func (e *Environments) lookupVar(varName string) (RuntimeVal, error) {
	env, err := e.resolve(varName)
	if err != nil {
		return nil, err
	}

	value := env.variables[varName]

	return value, nil
}
