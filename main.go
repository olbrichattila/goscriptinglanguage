package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	env, err := newEnvironments(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	mode := 2

	switch mode {
	case 1:
		propmt(env)
	case 2:
		executeScript(env)
	case 3:
		testing(env)
	case 4:
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

func executeScript(env *Environments) {
	s, err := readFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}

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

func readFromConsole() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Aty programming language")
	fmt.Println("-------------------------")

	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')

	return text
}

func readFromFile() (string, error) {
	if len(os.Args) < 2 {
		fmt.Println()
		return "", fmt.Errorf("Please provide the file name to run.")
	}

	return readFile(os.Args[1])
}
