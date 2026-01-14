package parser

import (
	"fmt"

	"github.com/cowellmi/stint/ast"
	"github.com/cowellmi/stint/lexer"
	"github.com/cowellmi/stint/token"
)

type Parser struct {
	l       *lexer.Lexer
	curTok  token.Token
	peekTok token.Token
	errors  []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read curTok and peekTok
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseTemplate() *ast.TemplateNode {
	t := &ast.TemplateNode{}
	for p.curTok.Type != token.EOF {
		n := p.parseNode()
		if n != nil {
			t.Nodes = append(t.Nodes, n)
		}
		p.nextToken()
	}
	return t
}

func (p *Parser) parseNode() ast.Node {
	switch p.curTok.Type {
	case token.RAW:
		return &ast.RawNode{
			Token: p.curTok,
			Value: p.curTok.Literal,
		}
	case token.TAG:
		return p.parseInterpolation()
	default:
		return nil
	}
}

func (p *Parser) parseInterpolation() *ast.InterpolationNode {
	node := &ast.InterpolationNode{Token: p.curTok}

	// Identifier
	if !p.expectPeek(token.IDENT) {
		p.recoverUntil(token.TAG)
		return nil
	}
	node.Name = &ast.Identifier{
		Token: p.curTok,
		Value: p.curTok.Literal,
	}

	// Constraints
	for p.peekTok.Type == token.COLON {
		p.nextToken() // move to ':'
		p.nextToken() // move to the thing after ':'

		switch p.curTok.Literal {
		case ast.ConstraintInteger:
			c := ast.Constraint{Name: ast.ConstraintInteger}
			node.Constraints = append(node.Constraints, c)
		case ast.ConstraintLength:
			c := p.parseConstraintArgs(ast.ConstraintLength)
			node.Constraints = append(node.Constraints, c)
		default:
			msg := fmt.Sprintf("unknown constraint: %s", p.curTok.Literal)
			p.errors = append(p.errors, msg)
		}
	}

	// Consume closing tag
	if !p.expectPeek(token.TAG) {
		return nil
	}

	return node
}

func (p *Parser) parseConstraintArgs(name string) ast.Constraint {
	c := ast.Constraint{Name: name}

	if !p.expectPeek(token.LPAREN) {
		return c
	}

	for p.peekTok.Type != token.RPAREN && p.peekTok.Type != token.EOF {
		p.nextToken()
		switch p.curTok.Type {
		case token.INT:
			c.Args = append(c.Args, p.curTok.Literal)
		case token.COMMA:
			continue
		default:
			msg := fmt.Sprintf("illegal token %s (%q) in constraint arguments",
				p.curTok.Type, p.curTok.Literal)
			p.errors = append(p.errors, msg)
			p.recoverUntil(token.RPAREN)
			return c
		}
	}

	p.expectPeek(token.RPAREN)
	return c
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTok.Type != t {
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekTok.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) recoverUntil(t token.TokenType) {
	for p.peekTok.Type != t && p.peekTok.Type != token.EOF {
		p.nextToken()
	}
}
