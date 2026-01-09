package parser

import (
	"testing"

	"github.com/cowellmi/stint/ast"
	"github.com/cowellmi/stint/lexer"
)

func TestStringInterpolation(t *testing.T) {
	input := "foo%bar:int:len(2)%"

	l := lexer.New(input)
	p := New(l)

	tmpl := p.ParseTemplate()
	if tmpl == nil {
		t.Fatalf("ParseTemplate() returned nil")
	}

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors: %v", len(p.Errors()), p.Errors())
	}

	if len(tmpl.Nodes) != 2 {
		t.Fatalf("tmpl.Nodes does not contain 2 nodes. got=%d",
			len(tmpl.Nodes))
	}

	// Raw node
	rnode, ok := tmpl.Nodes[0].(*ast.RawNode)
	if !ok {
		t.Fatalf("node[0] is not *ast.RawNode. got=%T", tmpl.Nodes[0])
	}
	if rnode.Value != "foo" {
		t.Errorf("rnode.Value wrong. expected=%q, got=%q", "foo", rnode.Value)
	}

	// Interpolation node
	inode, ok := tmpl.Nodes[1].(*ast.InterpolationNode)
	if !ok {
		t.Fatalf("node[1] is not *ast.InterpolationNode. got=%T", tmpl.Nodes[1])
	}

	if inode.Name.Value != "bar" {
		t.Errorf("inode.Name wrong. expected=%q, got=%q", "abc", inode.Name.Value)
	}

	// Constraints
	tests := []struct {
		expectedName string
		expectedArgs []string
	}{
		{ast.ConstraintInteger, []string{}},
		{ast.ConstraintLength, []string{"2"}},
	}

	if len(inode.Constraints) != len(tests) {
		t.Fatalf("constraints length wrong. expected=%d, got=%d",
			len(tests), len(inode.Constraints))
	}

	for i, tt := range tests {
		cons := inode.Constraints[i]
		if cons.Name != tt.expectedName {
			t.Errorf("tests[%d] - name wrong. expected=%q, got=%q",
				i, tt.expectedName, cons.Name)
		}

		if len(cons.Args) != len(tt.expectedArgs) {
			t.Fatalf("tests[%d] - args length wrong. expected=%d, got=%d",
				i, len(tt.expectedArgs), len(cons.Args))
		}

		for j, arg := range tt.expectedArgs {
			if cons.Args[j] != arg {
				t.Errorf("tests[%d] - arg[%d] wrong. expected=%q, got=%q",
					i, j, arg, cons.Args[j])
			}
		}
	}
}
