package stint

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cowellmi/stint/ast"
	"github.com/cowellmi/stint/lexer"
	"github.com/cowellmi/stint/parser"
)

type Template struct {
	tree *ast.TemplateNode
}

func NewTemplate(input string) (*Template, error) {
	l := lexer.New(input)
	p := parser.New(l)
	tree := p.ParseTemplate()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parser errors: %s", strings.Join(p.Errors(), "; "))
	}

	return &Template{tree: tree}, nil
}

func (t *Template) Execute(w io.Writer, env map[string]string) error {
	for _, node := range t.tree.Nodes {
		out, err := eval(node, env)
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)
	}
	return nil
}

func (t *Template) Vars() []*ast.InterpolationNode {
	var vars []*ast.InterpolationNode
	for _, n := range t.tree.Nodes {
		if v, ok := n.(*ast.InterpolationNode); ok {
			vars = append(vars, v)
		}
	}
	return vars
}

func eval(node ast.Node, env map[string]string) (string, error) {
	switch n := node.(type) {
	case *ast.RawNode:
		return n.Value, nil
	case *ast.InterpolationNode:
		return evalInterpolation(n, env)
	default:
		return "", fmt.Errorf("unknown node type: %T", node)
	}
}

func evalInterpolation(n *ast.InterpolationNode, env map[string]string) (string, error) {
	val, ok := env[n.Name.Value]
	if !ok {
		return "", fmt.Errorf("variable %s not found", n.Name.Value)
	}

	for _, c := range n.Constraints {
		if err := checkConstraint(c, val); err != nil {
			return "", err
		}
	}

	return val, nil
}

func checkConstraint(c ast.Constraint, val string) error {
	switch c.Name {
	case ast.ConstraintInteger:
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("int constraint: value %s must be an integer", val)
		}
	case ast.ConstraintLength:
		if err := checkLength(c.Args, val); err != nil {
			return fmt.Errorf("len constraint: %w", err)
		}
	default:
		args := strings.Join(c.Args, ",")
		return fmt.Errorf("unknown constraint: %s(%s)", c.Name, args)
	}
	return nil
}

func checkLength(args []string, val string) error {
	l := len(val)
	switch len(args) {
	case 0:
		return errors.New("missing argument(s)")
	case 1:
		target, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("argument %s must be an integer", args[0])
		}
		if l != target {
			return fmt.Errorf("expected length %d, got %d", target, l)
		}
	case 2:
		min, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("argument %s must be an integer", args[0])
		}
		max, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("argument %s must be an integer", args[1])
		}
		if l < min || l > max {
			return fmt.Errorf("length %d outside range %d-%d", l, min, max)
		}
	default:
		return fmt.Errorf("invalid arguments:")
	}
	return nil
}
