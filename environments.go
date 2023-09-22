package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Environments struct {
	parent    *Environments
	variables map[string]RuntimeVal
	constants map[string]interface{}
}

func newEnvironments(parent *Environments) (*Environments, *CustomError) {
	e := &Environments{
		parent:    parent,
		variables: make(map[string]RuntimeVal),
		constants: make(map[string]interface{}),
	}

	if parent == nil {
		err := e.declareDefaultEnv()
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Environments) declareVar(varName string, value RuntimeVal, constant bool) (RuntimeVal, *CustomError) {
	_, exist := e.variables[varName]
	if exist {
		return nil, newCustomError(
			fmt.Sprintf("Variable %s already exists", varName),
		)
	}
	e.variables[varName] = value
	if constant {
		e.constants[varName] = nil
	}

	return value, nil
}

func (e *Environments) assignVar(varName string, value RuntimeVal) (RuntimeVal, *CustomError) {
	env, err := e.resolve(varName)
	if err != nil {
		return nil, err
	}

	_, exist := e.constants[varName]
	if exist {
		return nil, newCustomError(fmt.Sprintf("Constant variable %s cannot be updated", varName))
	}

	env.variables[varName] = value

	return value, nil
}

func (e *Environments) resolve(varName string) (*Environments, *CustomError) {
	_, exist := e.variables[varName]
	if exist {
		return e, nil
	}
	if e.parent == nil {
		return nil, newCustomError(fmt.Sprintf("Variable %s could not be resolved", varName))
	}

	return e.parent.resolve(varName)
}

func (e *Environments) lookupVar(varName string) (RuntimeVal, *CustomError) {
	env, err := e.resolve(varName)
	if err != nil {
		return nil, err
	}

	value := env.variables[varName]

	return value, nil
}

func (e *Environments) declareDefaultEnv() *CustomError {
	rand.Seed(time.Now().UnixNano())
	_, err := e.declareVar("null", makeNull(), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("true", makeBool(true), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("false", makeBool(false), true)
	if err != nil {
		return err
	}

	// Define native function
	_, err = e.declareVar("print", makeNativeFn(ntPrint), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("println", makeNativeFn(ntPrintLn), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("time", makeNativeFn(ntTime), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("numToStr", makeNativeFn(ntNumToString), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("strToNum", makeNativeFn(ntStringToNum), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("input", makeNativeFn(ntInput), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("round", makeNativeFn(ntRound), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("rand", makeNativeFn(ntRand), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("len", makeNativeFn(ntLen), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("substr", makeNativeFn(ntSubstr), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("sleep", makeNativeFn(ntSleep), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("fileRead", makeNativeFn(ntFileRead), true)
	if err != nil {
		return err
	}

	_, err = e.declareVar("fileWrite", makeNativeFn(ntFileWrite), true)
	if err != nil {
		return err
	}

	return nil
}
