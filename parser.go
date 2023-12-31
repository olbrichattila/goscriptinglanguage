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

func (p *Parser) produceAST(sourceCode string) (*Program, *CustomError) {
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

func (p *Parser) expect(t TokenType, errMsg string) (*Token, *CustomError) {
	prev := p.next()
	if prev.Type != t {
		return nil, newCustomError(errMsg).addTrace(prev.Pos)
	}

	return &prev, nil
}

func (p *Parser) parseStmt() (Stmter, *CustomError) {
	switch p.at().Type {
	case TokenTypeLet, TokenTypeConst:
		return p.parseVarDeclaration()
	case TokenTypeFn:
		return p.parseFunctionDeclaration()
	case TokenTypeIf:
		return p.parseIfExpression()
	case TokenTypeFor:
		return p.parseForExpression()
	case TokenTypeBreak:
		return p.parseBreakExpression()
	case TokenTypeContinue:
		return p.parseContinueExpression()
	case TokenTypeSwitch:
		return p.parseSwitchExpression()
	default:
		return p.parseExpr()
	}
}

func (p *Parser) parseVarDeclaration() (Stmter, *CustomError) {
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
			return nil, newCustomError("Must assign value to constant experssion, no value provided")
		}

		return &VariableDeclaration{
			Stmt:       &Stmt{kind: NodeTypeVariableDeclaration, pos: p.at().Pos},
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
		Stmt:       &Stmt{kind: NodeTypeVariableDeclaration, pos: p.at().Pos},
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

func (p *Parser) parseFunctionDeclaration() (Stmter, *CustomError) {
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
			return nil, newCustomError("Inside function declatation expected paramegters to be string")
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
		Stmt:       &Stmt{kind: NodeTypeFunctionDeclaration, pos: p.at().Pos},
		parameters: params,
		name:       name,
		body:       body,
	}, nil
}

func (p *Parser) parseIfExpression() (Stmter, *CustomError) {
	tType := p.next().Type
	var cond Stmter
	var err *CustomError

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
		Stmt:           &Stmt{kind: NodeTypeIfExpression, pos: p.at().Pos},
		condition:      cond,
		body:           body,
		elseExpression: elseExpression,
	}, nil
}

func (p *Parser) parseForExpression() (Stmter, *CustomError) {
	// @TODO refactor to decrease complexity, factor out syntax variations phraser
	var err *CustomError
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
		Stmt:                  &Stmt{kind: NodeTypeForExpression, pos: p.at().Pos},
		declaration:           declaration,
		condition:             condition,
		afterCondition:        afterCondition,
		incrementalExpression: incrementalExpression,
		body:                  body,
	}, nil
}

func (p *Parser) parseSwitchExpression() (Stmter, *CustomError) {
	// @Todo refactor this, too complex
	p.next()
	_, err := p.expect(TokenTypeOpenParen, "Open parenthesis expected after switch")
	if err != nil {
		return nil, err
	}

	v, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(TokenTypeCloseParen, "Close parenthesis expected after switch statement value definition")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(TokenTypeOpenBrace, "Open brace expected after switch condition")
	if err != nil {
		return nil, err
	}

	var body []SwitchCaseExpression
	for {
		if p.at().Type == TokenTypeEOF || p.at().Type == TokenTypeCloseBrace {
			break
		}

		if p.at().Type == TokenTypeCase || p.at().Type == TokenTypeDefault {
			prevToken := p.next()
			valueToken := p.at()

			var comp RuntimeVal

			if prevToken.Type == TokenTypeDefault {
				comp = makeBool(true)
			} else if valueToken.Type == TokenTypeNumber {
				n, err := strconv.ParseFloat(valueToken.Value, 64)
				if err != nil {
					return nil, newCustomError(fmt.Sprintf("%s is not a number", valueToken.Value)).addTrace(valueToken.Pos)
				}
				comp = makeNumber(n)
				p.next()
			} else if valueToken.Type == TokenTypeString {
				comp = makeString(valueToken.Value)
				p.next()
			} else if valueToken.Type != TokenTypeDefault {
				return nil, newCustomError("Swich case condition can be string or numbar only").addTrace(valueToken.Pos)
			}

			_, err := p.expect(TokenTypeColon, "Colon expected after case value")
			if err != nil {
				return nil, err
			}

			var swBody []Stmter

			for {
				if p.at().Type == TokenTypeEOF || p.at().Type == TokenTypeCase || p.at().Type == TokenTypeDefault || p.at().Type == TokenTypeCloseBrace {
					break
				}

				s, err := p.parseStmt()
				if err != nil {
					return nil, err
				}

				swBody = append(swBody, s)

			}
			body = append(body, SwitchCaseExpression{compare: comp, body: swBody, pos: p.at().Pos})
		}
	}

	_, err = p.expect(TokenTypeCloseBrace, "CV")
	if err != nil {
		return nil, err
	}

	return &SwitchExpression{
		Stmt:  &Stmt{kind: NodeTypeSwitchExpression, pos: p.at().Pos},
		value: v,
		body:  body,
	}, nil
}

func (p *Parser) parseBreakExpression() (Stmter, *CustomError) {
	p.next()
	return &BreakExpression{
		Stmt: &Stmt{kind: NodeTypeBreakExpression, pos: p.at().Pos},
	}, nil
}

func (p *Parser) parseContinueExpression() (Stmter, *CustomError) {
	p.next()
	return &ContinueExpression{
		Stmt: &Stmt{kind: NodeTypeContinueExpression, pos: p.at().Pos},
	}, nil
}

func (p *Parser) parseExpr() (Stmter, *CustomError) {
	return p.parseConditionalExpr()
}

func (p *Parser) parseConditionalExpr() (Stmter, *CustomError) {
	left, err := p.parseAssignmentExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeSmaller ||
		p.at().Type == TokenTypeSmallerEqual ||
		p.at().Type == TokenTypeGreater ||
		p.at().Type == TokenTypeGreaterEqual ||
		p.at().Type == TokenTypeDoubeEqual ||
		p.at().Type == TokenTypeNotEqual ||
		p.at().Type == TokenTypeAnd ||
		p.at().Type == TokenTypeOr ||
		p.at().Type == TokenTypeNot {
		operator := p.next().Value
		right, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		return &ConditionDeclaration{
			Stmt:     &Stmt{kind: NodeTypeConditionExpression, pos: p.at().Pos},
			left:     left,
			right:    right,
			operator: operator,
		}, nil
	}

	return left, nil
}

func (p *Parser) parseAssignmentExpr() (Stmter, *CustomError) {
	left, err := p.parseObjectExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeEquals {
		p.next()
		value, err := p.parseObjectExpr()
		if err != nil {
			return nil, err
		}

		return &AssignmentExpr{Stmt: &Stmt{kind: NodeTypeAssigmentExpression, pos: p.at().Pos}, value: value, assigne: left}, nil

	}

	return left, nil
}

func (p *Parser) parseObjectExpr() (Stmter, *CustomError) {

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
			properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty, pos: p.at().Pos}, key: key})
			continue
		}

		if p.at().Type == TokenTypeCloseBrace {
			properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty, pos: p.at().Pos}, key: key})
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

		properties = append(properties, &Property{Stmt: &Stmt{kind: NodeTypeProperty, pos: p.at().Pos}, key: key, value: value})

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

	return &ObjectLiteral{Stmt: &Stmt{kind: NodeTypeObjectLiteral, pos: p.at().Pos}, properties: properties}, nil
}

func (p *Parser) parseAdditiveExpr() (Stmter, *CustomError) {
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
				Stmt:     &Stmt{kind: NodeTypeBinaryExpession, pos: p.at().Pos},
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

func (p *Parser) parseMultiplicativeExpr() (Stmter, *CustomError) {
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
				Stmt:     &Stmt{kind: NodeTypeBinaryExpession, pos: p.at().Pos},
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

func (p *Parser) parseCallMemberExpr() (Stmter, *CustomError) {
	member, err := p.parseMemberExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TokenTypeOpenParen {
		return p.parseCallExpr(member)
	}

	return member, nil
}

func (p *Parser) parseCallExpr(caller Stmter) (Stmter, *CustomError) {

	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	callExpr := &CallExpression{
		Stmt:   &Stmt{kind: NodeTypeCallExpression, pos: p.at().Pos},
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

func (p *Parser) parseMemberExpr() (Stmter, *CustomError) {
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
				return nil, newCustomError("Cannot use operatior without right hand side being an identifier")
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
			Stmt:     &Stmt{kind: NodeTypeMemberExpression, pos: p.at().Pos},
			object:   object,
			propert:  property,
			computed: computed,
		}
	}

	return object, nil
}

func (p *Parser) parseArgs() ([]*Stmter, *CustomError) {
	var err *CustomError
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

func (p *Parser) parseArgumentsLists() ([]*Stmter, *CustomError) {
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

func (p *Parser) parsePrimaryExpr() (Stmter, *CustomError) {
	tk := p.at().Type
	pos := p.at().Pos

	switch tk {
	case TokenTypeIdentifier:
		return &Identifier{Stmt: &Stmt{kind: NodeTypeIdentifier, pos: p.at().Pos}, symbol: p.next().Value}, nil
	case TokenTypeNumber:
		value, err := strconv.ParseFloat(p.next().Value, 64)
		if err != nil {
			return nil, newCustomError(err.Error())
		}
		return &NumericLiteral{Stmt: &Stmt{kind: NodeTypeNumericLiteral, pos: p.at().Pos}, value: value}, nil
	case TokenTypeString:
		return &StringLiteral{Stmt: &Stmt{kind: NodeTypeStringLIteral, pos: p.at().Pos}, value: p.next().Value}, nil
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
		return nil, newCustomError(fmt.Sprintf("Invalid token type %d %T", p.at().Type, p.at().Value)).addTrace(pos)
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
