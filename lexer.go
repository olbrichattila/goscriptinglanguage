package main

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenTypeNumber TokenType = iota
	TokenTypeIdentifier
	TokenTypeEquals
	TokenTypeOpenParen
	TokenTypeColseParen
	TokenTypeBinaryOperator
	TokenTypeLet
	TokenTypeConst
	TokenTypeSemicolon
	TokenTypeEOF
)

type Token struct {
	Value string
	Type  TokenType
}

type Tokenizer struct {
}

func newTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) tokenize(sourceCode string) ([]Token, error) {
	keywords := map[string]TokenType{
		"let":   TokenTypeLet,
		"const": TokenTypeConst,
	}

	var tokens []Token
	src := strings.Split(sourceCode, "")

	i := 0
	for {
		if i == len(src) {
			break
		}

		switch src[i] {
		case "(":
			tokens = append(tokens, Token{Type: TokenTypeOpenParen})
			i++
		case ")":
			tokens = append(tokens, Token{Type: TokenTypeColseParen})
			i++
		case "+", "-", "/", "*", "%":
			tokens = append(tokens, Token{Type: TokenTypeBinaryOperator, Value: src[i]})
			i++
		case "=":
			tokens = append(tokens, Token{Type: TokenTypeEquals})
			i++
		case ";":
			tokens = append(tokens, Token{Type: TokenTypeSemicolon})
			i++
		default:
			if t.isInt(src[i]) {
				num := ""
				for {
					if i == len(src) || !t.isInt(src[i]) {
						break
					}
					num += src[i]
					i++
				}

				tokens = append(tokens, Token{Type: TokenTypeNumber, Value: num})

			} else if t.isAlpha(src[i]) {
				alpha := ""
				for {
					if i == len(src) || !t.isAlpha(src[i]) {
						break
					}
					alpha += src[i]
					i++
				}

				if keywordTokenType, exist := keywords[alpha]; exist {
					tokens = append(tokens, Token{Type: keywordTokenType, Value: alpha})
				} else {
					tokens = append(tokens, Token{Type: TokenTypeIdentifier, Value: alpha})
				}

			} else if t.isSkippable(src[i]) {
				i++
			} else {
				return nil, fmt.Errorf("Uncrecoginized charecter found in source %s", src[i])
			}

		}
	}

	tokens = append(tokens, Token{Type: TokenTypeEOF, Value: "EndOfFile"})

	return tokens, nil
}

func (t *Tokenizer) isSkippable(s string) bool {
	return s == " " || s == "\n" || s == "\t"
}

func (t *Tokenizer) isInt(s string) bool {
	return s >= "0" && s <= "9"
}

func (t *Tokenizer) isAlpha(s string) bool {
	return strings.ToLower(s) != strings.ToUpper(s)
}
