package main

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenTypeNumber TokenType = iota
	TokenTypeString
	TokenTypeIdentifier
	TokenTypeEquals
	TokenTypeComma
	TokenTypeDot
	TokenTypeColon
	TokenTypeOpenParen
	TokenTypeCloseParen
	TokenTypeOpenBrace
	TokenTypeCloseBrace
	TokenTypeOpenBracket
	TokenTypeCloseBracket
	TokenTypeSmaller
	TokenTypeGreater
	TokenTypeSmallerEqual
	TokenTypeGreaterEqual
	TokenTypeDoubeEqual
	TokenTypeNotEqual
	TokenTypeNot
	TokenTypeBinaryOperator
	TokenTypeLet
	TokenTypeConst
	TokenTypeSemicolon
	TokenTypeFn
	TokenIf
	TokenFor
	TokenTypeEOF
)

type Token struct {
	Value string
	Type  TokenType
}

type Tokenizer struct {
	keywords map[string]TokenType
}

func newTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) tokenize(sourceCode string) ([]Token, error) {
	t.keywords = map[string]TokenType{
		"let":   TokenTypeLet,
		"const": TokenTypeConst,
		"fn":    TokenTypeFn,
		"if":    TokenIf,
		"for":   TokenFor,
	}

	var tokens []Token
	src := strings.Split(sourceCode, "")

	i := 0
	srcLen := len(src)
	for {
		if i == srcLen {
			break
		}

		switch src[i] {
		case "(":
			tokens = append(tokens, Token{Type: TokenTypeOpenParen})
			i++
		case ")":
			tokens = append(tokens, Token{Type: TokenTypeCloseParen})
			i++
		case "{":
			tokens = append(tokens, Token{Type: TokenTypeOpenBrace})
			i++
		case "}":
			tokens = append(tokens, Token{Type: TokenTypeCloseBrace})
			i++
		case "[":
			tokens = append(tokens, Token{Type: TokenTypeOpenBracket})
			i++
		case "]":
			tokens = append(tokens, Token{Type: TokenTypeCloseBracket})
			i++
		case "+", "-", "/", "*", "%":
			tokens = append(tokens, Token{Type: TokenTypeBinaryOperator, Value: src[i]})
			i++
		case "=":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeDoubeEqual, Value: "="})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeEquals})
			}
			i++
		case ";":
			tokens = append(tokens, Token{Type: TokenTypeSemicolon})
			i++
		case ":":
			tokens = append(tokens, Token{Type: TokenTypeColon})
			i++
		case ",":
			tokens = append(tokens, Token{Type: TokenTypeComma})
			i++
		case ".":
			tokens = append(tokens, Token{Type: TokenTypeDot})
			i++
		case "<":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeSmallerEqual, Value: "<="})
				i++
			} else if i < srcLen-1 && src[i+1] == ">" {
				tokens = append(tokens, Token{Type: TokenTypeNotEqual, Value: "!="})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeSmaller, Value: "<"})
			}
			i++
		case "!":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeNotEqual, Value: "!="})
				i++
			} else {
				// This is not yet pharsed
				tokens = append(tokens, Token{Type: TokenTypeNot, Value: "!"})
			}
			i++
		case ">":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeGreaterEqual, Value: ">="})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeGreater, Value: ">"})
			}
			i++
		default:
			tk, index, err := t.tokenizeComplex(src, i)
			if err != nil {
				return nil, err
			}

			if tk != nil {
				tokens = append(tokens, *tk)
			}

			i = index
		}
	}

	tokens = append(tokens, Token{Type: TokenTypeEOF, Value: "EndOfFile"})

	return tokens, nil
}

func (t *Tokenizer) tokenizeComplex(src []string, i int) (*Token, int, error) {
	if t.isInt(src[i]) {
		num := ""
		for {
			if i == len(src) || !t.isInt(src[i]) {
				break
			}
			num += src[i]
			i++
		}

		return &Token{Type: TokenTypeNumber, Value: num}, i, nil
	}

	if t.isAlpha(src[i]) {
		alpha := ""
		for {
			if i == len(src) || !t.isAlpha(src[i]) {
				break
			}
			alpha += src[i]
			i++
		}

		if keywordTokenType, exist := t.keywords[alpha]; exist {
			return &Token{Type: keywordTokenType, Value: alpha}, i, nil
		}

		return &Token{Type: TokenTypeIdentifier, Value: alpha}, i, nil
	}

	if src[i] == "\"" {
		str := ""
		if i == len(src) {
			return nil, i, fmt.Errorf("After opening quote there should be at least one closin quote")
		}
		i++

		for {
			if i < len(src)-1 && src[i] == "\"" && src[i+1] == "\"" {
				str += "\""
				i++
				i++
			}

			if i == len(src) || src[i] == "\"" {
				break
			}
			str += src[i]
			i++
		}

		if i < len(src) {
			i++
		}

		return &Token{Type: TokenTypeString, Value: str}, i, nil
	}

	if t.isSkippable(src[i]) {
		i++
		return nil, i, nil
	}

	return nil, i, fmt.Errorf("Uncrecoginized charecter found in source %s", src[i])
}

func (t *Tokenizer) isSkippable(s string) bool {
	return s == " " || s == "\n" || s == "\t" || s == "\r"
}

func (t *Tokenizer) isInt(s string) bool {
	return s >= "0" && s <= "9"
}

func (t *Tokenizer) isAlpha(s string) bool {
	return strings.ToLower(s) != strings.ToUpper(s)
}
