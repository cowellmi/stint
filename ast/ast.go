package ast

import "github.com/cowellmi/stint/token"

type Node interface {
	TokenLiteral() string
}

type TemplateNode struct {
	Nodes []Node
}

func (t *TemplateNode) TokenLiteral() string {
	if len(t.Nodes) > 0 {
		return t.Nodes[0].TokenLiteral()
	} else {
		return ""
	}
}

type RawNode struct {
	Token token.Token // the token.RAW token
	Value string
}

func (n *RawNode) TokenLiteral() string {
	return n.Token.Literal
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

const (
	ConstraintInteger = "int"
	ConstraintLength  = "len"
)

type Constraint struct {
	Name string
	Args []string
}

type InterpolationNode struct {
	Token       token.Token // the token.TAG token
	Name        *Identifier
	Constraints []Constraint
}

func (n *InterpolationNode) TokenLiteral() string {
	return n.Token.Literal
}
