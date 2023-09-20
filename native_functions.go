package main

import (
	"fmt"
	"strconv"
	"time"
)

type FunctionCall func([]RuntimeVal, *Environments) RuntimeVal

func ntPrint(args []RuntimeVal, env *Environments) RuntimeVal {
	for _, arg := range args {
		if v, ok := arg.(*NumberVal); ok {
			fmt.Println(v.Value)
			continue
		}

		if v, ok := arg.(*StringVal); ok {
			fmt.Println(v.Value)
			continue
		}

		if v, ok := arg.(*NullVal); ok {
			fmt.Println(v.Value)
			continue
		}

		if v, ok := arg.(*BoolVal); ok {
			fmt.Println(v.Value)
			continue
		}

		if v, ok := arg.(*ObjectVal); ok {
			fmt.Println(v.properties)
			continue
		}

		fmt.Print("Cannot print this data type")
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
