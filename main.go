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

	if len(os.Args) > 1 {
		if os.Args[1] == "prompt" {
			mode = 1
		}
	}

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
	tokens, err := t.tokenize("x >= 10")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tokens)
}

func testing(env *Environments) {

	s, _ := readFile("./examples/test.gl")

	p := newParser()
	parsed, err := p.produceAST(s)
	if err != nil {
		fmt.Println(err.message, err.trace)
		return
	}

	i := newInterpreter()
	_, cErr := i.evaluate(parsed, env)
	if cErr != nil {
		fmt.Println(cErr.message, cErr.trace)
		return
	}
}

func propmt(env *Environments) {
	for {
		s := readFromConsole()
		if s == "\n" {
			break
		}

		p := newParser()
		parsed, err := p.produceAST(s)
		if err != nil {
			fmt.Println(err.message, err.trace)
			continue
		}

		i := newInterpreter()
		_, cErr := i.evaluate(parsed, env)
		if cErr != nil {
			fmt.Println(cErr.message, cErr.trace)

			continue
		}
	}
}

func executeScript(env *Environments) {
	s, err := readFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	p := newParser()
	parsed, pErr := p.produceAST(s)
	if pErr != nil {
		fmt.Println(pErr.message, pErr.trace)
		return
	}

	i := newInterpreter()
	_, cErr := i.evaluate(parsed, env)
	if cErr != nil {
		fmt.Println(cErr.message, cErr.trace)
		return
	}
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

func prettyPrint(p *Program) {
	for _, val := range p.body {
		displayStruct(val)
	}
}

func displayStruct(s Stmter) {
	fmt.Printf("%s\n", s)
}
