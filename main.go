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
	mode := 3

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

	s, _ := readFile("./examples/expr1.gl")

	p := newParser()
	parsed, err := p.produceAST(s)
	if err != nil {
		displayError(err, &s)
		return
	}

	i := newInterpreter()
	_, cErr := i.evaluate(parsed, env)
	if cErr != nil {
		displayError(cErr, &s)
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
			displayError(err, &s)
			continue
		}

		i := newInterpreter()
		_, cErr := i.evaluate(parsed, env)
		if cErr != nil {
			displayError(cErr, &s)
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
		displayError(pErr, &s)
		return
	}

	i := newInterpreter()
	_, cErr := i.evaluate(parsed, env)
	if cErr != nil {
		displayError(cErr, &s)
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

func displayError(err *CustomError, src *string) {
	// @todo refactor this, maybe separat struct
	red := "\033[31m"
	green := "\033[32m"
	reset := "\033[0m"

	fmt.Println()
	fmt.Println(red + err.message + reset)
	l := len(*src)
	str := *src

	for _, tr := range err.trace {
		line := 1
		pos := 1
		for i, c := range str {
			s := string(c)
			if s == "\n" {
				line++
				pos = 1
			}

			if i == tr {
				startPos := i - 3
				if startPos < 0 {
					startPos = 0
				}

				endPos := i + 3
				if endPos > l {
					endPos = l
				}

				fmt.Printf(
					green+"Error at line (%d), position (%d) near at: `%s`\n"+reset,
					line,
					pos,
					str[startPos:endPos],
				)
			}
			pos++
		}

		if tr == l {
			startPos := l - 6
			if startPos < 0 {
				startPos = 0
			}

			fmt.Printf(
				green+"Error at the end of the file near at: `%s`\n"+reset,
				str[startPos:l],
			)
		}
	}

	fmt.Println()
}
