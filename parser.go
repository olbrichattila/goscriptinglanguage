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

func (p *Parser) expect(t TokenType, errMsg string) (*Token, error) {
	prev := p.next()
	if prev.Type != t {
		return nil, fmt.Errorf(errMsg)
	}

	return &prev, nil
}

func (p *Parser) parseStmt() (Stmter, error) {
	switch p.at().Type {
	case TokenTypeLet, TokenTypeConst:
		return p.parseVarDeclaration()
	case TokenTypeFn:
		return p.parseFunctionDeclaration()
	case TokenTypeIf:
		return p.parseIfExpression()
	case TokenFor:
		return p.parseForExpression()
	default:
		return p.parseExpr()
	}
}

func (p *Parser) parseVarDeclaration() (Stmter, error) {
	tokenType := p.next().Type
	token := p.at()
	isConstant := tokenType == TokenTypeConst
	_, err := p.expect(TokenTypeIdentifier, "Expected identifier name following let or const keywords")
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

	_, err = p.expect(TokenTypeEquals, "Expected equals token following identifier in var declaration")
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

	_, err = p.expect(TokenTypeSemicolon, "Variable declaration must end with semilolon")
	if err != nil {
		return nil, err
	}

	return declaration, nil
}

func (p *Parser) parseFunctionDeclaration() (Stmter, error) {
	p.next()
	token, err := p.expect(TokenTypeIdentifier, "Expected function name following fn keyword")
	if err != nil {
		return nil, err
	}

	name := token.Value
	var params []string

	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	for _, arg := range args {
		a := *arg
		if a.Kind() != NodeTypeIdentifier {
			return nil, fmt.Errorf("Inside function declatation expected paramegters to be string")
		}

		params = append(params, a.(*Identifier).symbol)
	}

	_, err = p.expect(TokenTypeOpenBrace, "Expected function body declaration ")
	if err != nil {
		return nil, err
	}

	var body []Stmter

	for {
		if p.at().Type == TokenTypeEOF || p.at().Type == TokenTypeCloseBrace {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		body = append(body, s)
	}

	_, err = p.expect(TokenTypeCloseBrace, "Closing brace expected inside function declaration")
	if err != nil {
		return nil, err
	}

	return &FunctionDeclaration{
		Stmt:       &Stmt{kind: NodeTypeFunctionDeclaration},
		parameters: params,
		name:       name,
		body:       body,
	}, nil
}

func (p *Parser) parseIfExpression() (Stmter, error) {
	tType := p.next().Type
	var cond Stmter
	var err error

	if tType != TokenTypeElse {
		_, err := p.expect(TokenTypeOpenParen, "Open parenthesis expected after if statement")
		if err != nil {
			return nil, err
		}

		cond, err = p.parseConditionalExpr()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(TokenTypeCloseParen, "Close parenthesis expected after if statement conditions")
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(TokenTypeOpenBrace, "Expected open brace after if condition")
	if err != nil {
		return nil, err
	}

	var body []Stmter

	for {
		if p.at().Type == TokenTypeEOF || p.at().Type == TokenTypeCloseBrace {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		body = append(body, s)
	}

	_, err = p.expect(TokenTypeCloseBrace, "Closing brace expected inside function declaration")
	if err != nil {
		return nil, err
	}

	var elseExpression Stmter

	if p.at().Type == TokenTypeElse || p.at().Type == TokenTypeElseIf {
		elseExpression, err = p.parseIfExpression()
		if err != nil {
			return nil, err
		}
	}

	return &IfExpression{
		Stmt:           &Stmt{kind: NodeTypeIfExpression},
		condition:      cond,
		body:           body,
		elseExpression: elseExpression,
	}, nil
}

func (p *Parser) parseForExpression() (Stmter, error) {
	// @TODO refactor to decrease complexity, factor out syntax variations phraser
	var err error
	var afterCondition, declaration, condition, incrementalExpression Stmter
	p.next()

	if p.at().Type != TokenTypeOpenBrace {
		_, err = p.expect(TokenTypeOpenParen, "Open parenthesis expected after for statement")
		if err != nil {
			return nil, err
		}

		parCount := p.countFor(TokenTypeSemicolon)

		if parCount == 0 {
			condition, err = p.parseConditionalExpr()
			if err != nil {
				return nil, err
			}
		}

		if parCount > 0 {
			declaration, err = p.parseVarDeclaration()
			if err != nil {
				return nil, err
			}

			condition, err = p.parseConditionalExpr()
			if err != nil {
				return nil, err
			}

			_, err = p.expect(TokenTypeSemicolon, "Semicolon expected after for variable condition")
			if err != nil {
				return nil, err
			}

			incrementalExpression, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		}

		_, err := p.expect(TokenTypeCloseParen, "Close parenthesis expected after for statement conditions")
		if err != nil {
			return nil, err
		}

		_, err = p.expect(TokenTypeOpenBrace, "Expected open brace after if condition")
		if err != nil {
			return nil, err
		}
	} else {
		p.next()
	}

	var body []Stmter

	for {
		if p.at().Type == TokenTypeEOF || p.at().Type == TokenTypeCloseBrace {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		body = append(body, s)
	}

	_, err = p.expect(TokenTypeCloseBrace, "Closing brace expected inside function declaration")
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeOpenParen {
		p.next()
		afterCondition, err = p.parseConditionalExpr()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(TokenTypeCloseParen, "Close parenthesis expected after if statement closing conditions")
		if err != nil {
			return nil, err
		}
	}

	return &ForExpression{
		Stmt:                  &Stmt{kind: NodeTypeForExpression},
		declaration:           declaration,
		condition:             condition,
		afterCondition:        afterCondition,
		incrementalExpression: incrementalExpression,
		body:                  body,
	}, nil
}

func (p *Parser) parseExpr() (Stmter, error) {
	return p.parseConditionalExpr()
}

func (p *Parser) parseConditionalExpr() (Stmter, error) {
	left, err := p.parseAssignmentExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeSmaller || p.at().Type == TokenTypeSmallerEqual || p.at().Type == TokenTypeGreater || p.at().Type == TokenTypeGreaterEqual || p.at().Type == TokenTypeDoubeEqual || p.at().Type == TokenTypeNotEqual {
		operator := p.next().Value
		right, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		return &ConditionDeclaration{
			Stmt:     &Stmt{kind: NodeTypeConditionExpression},
			left:     left,
			right:    right,
			operator: operator,
		}, nil
	}

	return left, nil
}

func (p *Parser) parseAssignmentExpr() (Stmter, error) {
	left, err := p.parseObjectExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeEquals {
		p.next()
		value, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		return &AssignmentExpr{Stmt: &Stmt{kind: NodeTypeAssigmentExpression}, value: value, assigne: left}, nil

	}

	return left, nil
}

func (p *Parser) parseObjectExpr() (Stmter, error) {

	if p.at().Type != TokenTypeOpenBrace {
		return p.parseAdditiveExpr()
	}

	p.next()

	var properties []*Property

	for {
		if p.eof() || p.at().Type == TokenTypeCloseBrace {
			break
		}

		t, err := p.expect(TokenTypeIdentifier, "Object literal key expected")
		if err != nil {
			return nil, err
		}
		key := t.Value

		if p.at().Type == TokenTypeComma {
			p.next()
			properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty}, key: key})
			continue
		}

		if p.at().Type == TokenTypeCloseBrace {
			properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty}, key: key})
			continue
		}

		_, err = p.expect(TokenTypeColon, "Missing colon following in object expression")
		if err != nil {
			return nil, err
		}

		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty}, key: key, value: value})

		if p.at().Type != TokenTypeCloseBrace {
			_, err := p.expect(TokenTypeComma, "Expected comma or closing bracket following property")
			if err != nil {
				return nil, err
			}
		}

	}

	_, err := p.expect(TokenTypeCloseBrace, "Objects literal missing closing brace.")
	if err != nil {
		return nil, err
	}

	return &ObjectLiteral{Stmt: &Stmt{kind: NodeTypeObjectLiteral}, properties: properties}, nil
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
	left, err := p.parseCallMemberExpr()
	if err != nil {
		return nil, err
	}

	for {
		v := p.at().Value
		if v == "/" || v == "*" || v == "%" {
			operator := p.next().Value
			right, err := p.parseCallMemberExpr()
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

func (p *Parser) parseCallMemberExpr() (Stmter, error) {
	member, err := p.parseMemberExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeOpenParen {
		return p.parseCallExpr(member)
	}

	return member, nil
}

func (p *Parser) parseCallExpr(caller Stmter) (Stmter, error) {

	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	callExpr := &CallExpression{
		Stmt:   &Stmt{kind: NodeTypeCallExpression},
		caller: caller,
		args:   args,
	}

	if p.at().Type == TokenTypeOpenBrace {

		e, err := p.parseCallExpr(callExpr)
		if err != nil {
			return nil, err

		}
		callExpr = e.(*CallExpression)
	}

	return callExpr, nil
}

func (p *Parser) parseMemberExpr() (Stmter, error) {
	object, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		if p.at().Type != TokenTypeDot && p.at().Type != TokenTypeOpenBracket {
			break
		}

		operator := p.next()

		var property Stmter
		var computed bool

		if operator.Type == TokenTypeDot {
			computed = false
			property, err = p.parsePrimaryExpr()
			if err != nil {
				return nil, err
			}

			if property.Kind() != NodeTypeIdentifier {
				return nil, fmt.Errorf("Cannot use operatior without right hand side being an identifier")
			}
		} else {
			computed = true

			property, err = p.parseExpr()
			if err != nil {
				return nil, err
			}

			_, err := p.expect(TokenTypeCloseBracket, "Missing cosing bracket in computed value")
			if err != nil {
				return nil, err
			}
		}

		object = &MemberExpression{
			Stmt:     &Stmt{kind: NodeTypeMemberExpression},
			object:   object,
			propert:  property,
			computed: computed,
		}
	}

	return object, nil
}

func (p *Parser) parseArgs() ([]*Stmter, error) {
	var err error
	_, err = p.expect(TokenTypeOpenParen, "Expected open parantesis")
	if err != nil {
		return nil, err
	}
	var args []*Stmter

	if p.at().Type != TokenTypeCloseParen {
		args, err = p.parseArgumentsLists()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(TokenTypeCloseParen, "Missing closing parenthesis inside arguement list")
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (p *Parser) parseArgumentsLists() ([]*Stmter, error) {
	var args []*Stmter
	arg, err := p.parseAssignmentExpr()
	if err != nil {
		return nil, err
	}

	args = append(args, &arg)

	for {
		if p.at().Type != TokenTypeComma {
			break
		}

		p.next()
		fArg, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		args = append(args, &fArg)
	}

	return args, nil
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
	case TokenTypeString:
		return &StringLiteral{Stmt: &Stmt{kind: NodeTypeStringLIteral}, value: p.next().Value}, nil
	case TokenTypeOpenParen:
		p.next()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(TokenTypeCloseParen, "Unexpected token found inside parenthesised expression, expected closing parenthesis")
		if err != nil {
			return nil, err
		}

		return value, nil
	default:
		return nil, fmt.Errorf("Invalid token type %d %T", p.at().Type, p.at().Value)
	}
}

func (p *Parser) countFor(divider TokenType) int {
	i := p.index
	c := 0

	for {
		if p.eof() || p.tokens[i].Type == TokenTypeCloseParen {
			break
		}

		if p.tokens[i].Type == divider {
			c++
		}

		i++
	}

	return c
}
