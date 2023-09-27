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
	TokenTypeAnd
	TokenTypeOr
	TokenTypeBinaryOperator
	TokenTypeLet
	TokenTypeConst
	TokenTypeSemicolon
	TokenTypeFn
	TokenTypeIf
	TokenTypeElse
	TokenTypeElseIf
	TokenTypeFor
	TokenTypeSwitch
	TokenTypeCase
	TokenTypeDefault
	TokenTypeBreak
	TokenTypeContinue
	TokenTypeEOF
)

type Token struct {
	Value string
	Type  TokenType
	Pos   int
}

type Tokenizer struct {
	keywords map[string]TokenType
}

func newTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) tokenize(sourceCode string) ([]Token, *CustomError) {
	t.keywords = map[string]TokenType{
		"let":      TokenTypeLet,
		"const":    TokenTypeConst,
		"fn":       TokenTypeFn,
		"if":       TokenTypeIf,
		"else":     TokenTypeElse,
		"elseif":   TokenTypeElseIf,
		"for":      TokenTypeFor,
		"switch":   TokenTypeSwitch,
		"case":     TokenTypeCase,
		"default":  TokenTypeDefault,
		"break":    TokenTypeBreak,
		"continue": TokenTypeContinue,
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
			tokens = append(tokens, Token{Type: TokenTypeOpenParen, Pos: i})
			i++
		case ")":
			tokens = append(tokens, Token{Type: TokenTypeCloseParen, Pos: i})
			i++
		case "{":
			tokens = append(tokens, Token{Type: TokenTypeOpenBrace, Pos: i})
			i++
		case "}":
			tokens = append(tokens, Token{Type: TokenTypeCloseBrace, Pos: i})
			i++
		case "[":
			tokens = append(tokens, Token{Type: TokenTypeOpenBracket, Pos: i})
			i++
		case "]":
			tokens = append(tokens, Token{Type: TokenTypeCloseBracket, Pos: i})
			i++
		case "+", "-", "/", "*", "%":
			if i < srcLen-1 && src[i+1] == "/" {
				for {
					if i == len(src) || src[i] == "\n" {
						break
					}
					i++
				}
			} else {
				tokens = append(tokens, Token{Type: TokenTypeBinaryOperator, Value: src[i], Pos: i})
				i++
			}
		case "=":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeDoubeEqual, Value: "=", Pos: i})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeEquals, Pos: i})
			}
			i++
		case ";":
			tokens = append(tokens, Token{Type: TokenTypeSemicolon, Pos: i})
			i++
		case ":":
			tokens = append(tokens, Token{Type: TokenTypeColon, Pos: i})
			i++
		case ",":
			tokens = append(tokens, Token{Type: TokenTypeComma, Pos: i})
			i++
		case ".":
			tokens = append(tokens, Token{Type: TokenTypeDot, Pos: i})
			i++
		case "<":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeSmallerEqual, Value: "<=", Pos: i})
				i++
			} else if i < srcLen-1 && src[i+1] == ">" {
				tokens = append(tokens, Token{Type: TokenTypeNotEqual, Value: "!=", Pos: i})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeSmaller, Value: "<", Pos: i})
			}
			i++
		case "!":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeNotEqual, Value: "!=", Pos: i})
				i++
			} else {
				// This is not yet pharsed
				tokens = append(tokens, Token{Type: TokenTypeNot, Value: "!", Pos: i})
			}
			i++
		case ">":
			if i < srcLen-1 && src[i+1] == "=" {
				tokens = append(tokens, Token{Type: TokenTypeGreaterEqual, Value: ">=", Pos: i})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenTypeGreater, Value: ">", Pos: i})
			}
			i++
		case "&":
			if i < srcLen-1 && src[i+1] == "&" {
				tokens = append(tokens, Token{Type: TokenTypeAnd, Value: "&", Pos: i})
				i++
				i++
			} else {
				return nil, newCustomError("Condition requries doube &").addTrace(i)
			}
		case "|":
			if i < srcLen-1 && src[i+1] == "|" {
				tokens = append(tokens, Token{Type: TokenTypeOr, Value: "|", Pos: i})
				i++
				i++
			} else {
				return nil, newCustomError("Condition requries doube |").addTrace(i)
			}
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

	tokens = append(tokens, Token{Type: TokenTypeEOF, Value: "EndOfFile", Pos: i})

	return tokens, nil
}

func (t *Tokenizer) tokenizeComplex(src []string, i int) (*Token, int, *CustomError) {
	if t.isInt(src[i]) {
		num := ""
		for {
			if i == len(src) || !t.isInt(src[i]) {
				break
			}
			num += src[i]
			i++
		}

		return &Token{Type: TokenTypeNumber, Value: num, Pos: i}, i, nil
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
			return &Token{Type: keywordTokenType, Value: alpha, Pos: i}, i, nil
		}

		return &Token{Type: TokenTypeIdentifier, Value: alpha, Pos: i}, i, nil
	}

	if src[i] == "\"" {
		str := ""
		if i == len(src) {
			err := newCustomError("After opening quote there should be at least one closing quote")
			err.addTrace(i)
			return nil, i, err
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

		return &Token{Type: TokenTypeString, Value: str, Pos: i}, i, nil
	}

	if t.isSkippable(src[i]) {
		i++
		return nil, i, nil
	}

	return nil, i, newCustomError(fmt.Sprintf("Uncrecoginized charecter found in source %s", src[i])).addTrace(i)
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
