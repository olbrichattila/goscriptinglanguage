package main

import (
	"fmt"
	"strconv"
	"time"
)

type FunctionCall func([]RuntimeVal, *Environments) RuntimeVal

func ntPrint(args []RuntimeVal, env *Environments) RuntimeVal {
	return ntPrinter(args, env, false)
}

func ntPrintLn(args []RuntimeVal, env *Environments) RuntimeVal {
	return ntPrinter(args, env, true)
}

func ntPrinter(args []RuntimeVal, env *Environments, ln bool) RuntimeVal {
	succ := false
	for _, arg := range args {
		if v, ok := arg.(*NumberVal); ok {
			fmt.Print(v.Value)
			succ = true
		}

		if v, ok := arg.(*StringVal); ok {
			fmt.Print(v.Value)
			succ = true
		}

		if v, ok := arg.(*NullVal); ok {
			fmt.Print(v.Value)
			succ = true
		}

		if v, ok := arg.(*BoolVal); ok {
			fmt.Print(v.Value)
			succ = true
		}

		if v, ok := arg.(*ObjectVal); ok {
			fmt.Print(v.properties)
			succ = true
		}

		if succ == false {
			fmt.Print("Cannot print this data type")
		}

		if ln == true {
			fmt.Println()
		}

	}

	return makeNull()
}

func ntTime(args []RuntimeVal, env *Environments) RuntimeVal {
	currentTime := time.Now().Unix()

	return makeNumber(float64(currentTime))
}

func ntNumToString(args []RuntimeVal, env *Environments) RuntimeVal {
	if n, ok := args[0].(*NumberVal); ok {
		s := strconv.FormatFloat(n.Value, 'f', -1, 64)
		return makeString(s)
	}

	return makeNull()
}
