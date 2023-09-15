package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	env := declareDefaultEnv()

	mode := 1

	switch mode {
	case 1:
		propmt(env)
	case 2:
		testing(env)
	case 3:
		testTokenizer()
	default:
		propmt(env)
	}
}

func testTokenizer() {
	t := newTokenizer()
	tokens, err := t.tokenize("let x = 45 + (50 * 2) / (foo - vakk)")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tokens)
}

func testing(env *Environments) {
	s := "const x = 10;x = 5"
	p := newParser()
	parsed, err := p.produceAST(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	i := newInterpreter()
	e, err := i.evaluate(parsed, env)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println(e)
}

func propmt(env *Environments) {
	for {
		s := readFromConsole()
		if s == "\n" {
			break
		}
		// s := "let x = 10;"

		p := newParser()
		parsed, err := p.produceAST(s)
		if err != nil {
			fmt.Println(err)
			continue
		}

		i := newInterpreter()
		e, err := i.evaluate(parsed, env)
		if err != nil {
			fmt.Println(err)

			continue
		}
		fmt.Println(e)
		// fmt.Println(parsed)
	}
}

func declareDefaultEnv() *Environments {
	env := newEnvironments(nil)
	env.declareVar("null", MK_NULL(), true)
	env.declareVar("true", MK_BOOL(true), true)
	env.declareVar("false", MK_BOOL(false), true)

	return env
}

func readFromConsole() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Aty programming language")
	fmt.Println("-------------------------")

	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')

	return text
}
