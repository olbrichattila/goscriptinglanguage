package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	index  int
}

func newParser() *Parser {
	return &Parser{}
}

func (p *Parser) produceAST(sourceCode string) (*Program, error) {
	t := newTokenizer()
	tokens, err := t.tokenize(sourceCode)
	if err != nil {
		return nil, err
	}

	p.tokens = tokens

	pr := &Program{Stmt: &Stmt{kind: NodeTypeProgram}}

	for {
		if p.eof() {
			break
		}

		node, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		pr.body = append(pr.body, node)

	}

	return pr, nil
}

func (p *Parser) eof() bool {
	return p.tokens[p.index].Type == TokenTypeEOF
}

func (p *Parser) at() Token {
	return p.tokens[p.index]
}

func (p *Parser) next() Token {
	token := p.tokens[p.index]
	p.index++
	return token
}

func (p *Parser) expect(t TokenType, errMsg string) error {
	prev := p.next()
	if prev.Type != t {
		return fmt.Errorf(errMsg)
	}

	return nil
}

func (p *Parser) parseStmt() (Stmter, error) {
	switch p.at().Type {
	case TokenTypeLet, TokenTypeConst:
		return p.parseVarDeclaration()
	default:
		return p.parseExpr()
	}
}

func (p *Parser) parseVarDeclaration() (Stmter, error) {
	tokenType := p.next().Type
	token := p.at()
	isConstant := tokenType == TokenTypeConst
	err := p.expect(TokenTypeIdentifier, "Expected identifier name following let or const keywords")
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeSemicolon {
		p.next()
		if isConstant {
			return nil, fmt.Errorf("Must assign value to constant experssion, no value provided")
		}

		return &VariableDeclaration{
			Stmt:       &Stmt{kind: NodeTypeVariableDeclaration},
			identifier: token.Value,
			constant:   false,
		}, nil

	}

	err = p.expect(TokenTypeEquals, "Expected equals token following identifier in var declaration")
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	declaration, err := &VariableDeclaration{
		Stmt:       &Stmt{kind: NodeTypeVariableDeclaration},
		value:      expr,
		identifier: token.Value,
		constant:   isConstant,
	}, nil

	err = p.expect(TokenTypeSemicolon, "Variable declaration must end with semilolon")
	if err != nil {
		return nil, err
	}

	return declaration, nil
}

func (p *Parser) parseExpr() (Stmter, error) {
	return p.parseAssignmentExpr()
}

func (p *Parser) parseAssignmentExpr() (Stmter, error) {
	left, err := p.parseAdditiveExpr() // @todo swith this out with objects
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeEquals {
		p.next()
		value, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		return &AssignmentExpr{Stmt: &Stmt{kind: NodeTypeAssigmentExpr}, value: value, assigne: left}, nil

	}

	return left, nil
}

func (p *Parser) parseAdditiveExpr() (Stmter, error) {
	left, err := p.parseMultiplicativeExpr()
	if err != nil {
		return nil, err
	}

	for {
		v := p.at().Value
		if v == "+" || v == "-" {
			operator := p.next().Value
			right, err := p.parseMultiplicativeExpr()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpession{
				Stmt:     &Stmt{kind: NodeTypeBinaryExpession},
				left:     left,
				right:    right,
				operator: operator,
			}
			continue
		}
		break
	}

	return left, nil
}

func (p *Parser) parseMultiplicativeExpr() (Stmter, error) {
	left, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		v := p.at().Value
		if v == "/" || v == "*" || v == "%" {
			operator := p.next().Value
			right, err := p.parsePrimaryExpr()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpession{
				Stmt:     &Stmt{kind: NodeTypeBinaryExpession},
				left:     left,
				right:    right,
				operator: operator,
			}
			continue
		}
		break
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpr() (Stmter, error) {
	tk := p.at().Type

	switch tk {
	case TokenTypeIdentifier:
		return &Identifier{Stmt: &Stmt{kind: NodeTypeIdentifier}, symbol: p.next().Value}, nil
	case TokenTypeNumber:
		value, err := strconv.ParseFloat(p.next().Value, 64)
		if err != nil {
			return nil, err
		}
		return &NumericLiteral{Stmt: &Stmt{kind: NodeTypeNumericLiteral}, value: value}, nil
	case TokenTypeOpenParen:
		p.next()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		err = p.expect(TokenTypeColseParen, "Unexpected token found inside parenthesised expression, expected closing parenthesis")
		if err != nil {
			return nil, err
		}

		return value, nil
	default:
		return nil, fmt.Errorf("Invalid token type %d", p.at().Type)
	}
}
