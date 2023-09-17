package main

import (
	"fmt"
	"time"
)

type FunctionCall func([]RuntimeVal, *Environments) RuntimeVal

func ntPrint(args []RuntimeVal, env *Environments) RuntimeVal {
	for _, arg := range args {
		fmt.Print(arg)
	}

	return makeNull()
}

func ntTime(args []RuntimeVal, env *Environments) RuntimeVal {
	currentTime := time.Now().Unix()

	return makeNumber(float64(currentTime))
}
