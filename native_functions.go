package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

func ntStringToNum(args []RuntimeVal, env *Environments) RuntimeVal {
	if s, ok := args[0].(*StringVal); ok {
		n, err := strconv.ParseFloat(s.Value, 64)
		if err == nil {
			return makeNumber(n)
		}

	}

	return makeNull()
}

func ntInput(args []RuntimeVal, env *Environments) RuntimeVal {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return makeString(strings.TrimSuffix(text, "\n"))
}

func ntRound(args []RuntimeVal, env *Environments) RuntimeVal {
	if n, ok := args[0].(*NumberVal); ok {
		d := 1.0
		if len(args) > 1 {
			if n2, ok := args[1].(*NumberVal); ok {
				d = n2.Value * 10
			}
		}
		n := math.Round(n.Value*d) / d
		return makeNumber(n)
	}

	return makeNull()
}

func ntRand(args []RuntimeVal, env *Environments) RuntimeVal {
	rng := 100
	if n, ok := args[0].(*NumberVal); ok {
		rng = int(n.Value)
	}

	return makeNumber(float64(rand.Intn(rng)))
}
