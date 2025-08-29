package compoundassign

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "compoundassign",
	Doc:  "suggests using assignment operators like += instead of x = x + y",
	Run:  run,
}

var binToAssign = map[token.Token]token.Token{
	token.ADD:     token.ADD_ASSIGN,
	token.SUB:     token.SUB_ASSIGN,
	token.MUL:     token.MUL_ASSIGN,
	token.QUO:     token.QUO_ASSIGN,
	token.REM:     token.REM_ASSIGN,
	token.AND:     token.AND_ASSIGN,
	token.OR:      token.OR_ASSIGN,
	token.XOR:     token.XOR_ASSIGN,
	token.SHL:     token.SHL_ASSIGN,
	token.SHR:     token.SHR_ASSIGN,
	token.AND_NOT: token.AND_NOT_ASSIGN,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			// single-variable assignments
			if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				return true
			}
			// The right-hand side must be a binary expression
			bin, ok := assign.Rhs[0].(*ast.BinaryExpr)
			if !ok {
				return true
			}

			if !matchesIdentOrSelector(assign.Lhs[0], bin.X) {
				return true
			}

			// Only consider simple assignment (=)
			if assign.Tok != token.ASSIGN {
				return true
			}
			// Only handle supported operators
			switch bin.Op {
			case token.ADD,
				token.SUB,
				token.MUL,
				token.QUO,
				token.REM,
				token.AND,
				token.OR,
				token.XOR,
				token.SHL,
				token.SHR,
				token.AND_NOT:
				oldExpr := render(pass.Fset, assign)
				assign.Tok = binToAssign[bin.Op]
				assign.Rhs[0] = bin.Y
				newExpr := render(pass.Fset, assign)
				pass.Report(analysis.Diagnostic{
					Pos: assign.Pos(),
					Message: fmt.Sprintf(
						"use '%s' instead of '%s'",
						newExpr,
						oldExpr,
					),
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: fmt.Sprintf("should replace `%s` with `%s`", oldExpr, newExpr),
							TextEdits: []analysis.TextEdit{
								{
									Pos:     assign.Pos(),
									End:     assign.End(),
									NewText: []byte(newExpr),
								},
							},
						},
					},
				})
			}

			return true
		})
	}

	return nil, nil
}

func matchesIdentOrSelector(a, b ast.Expr) bool {
	switch a := a.(type) {
	case *ast.Ident:
		if b, ok := b.(*ast.Ident); ok {
			return a.Name == b.Name
		}

		return false

	case *ast.SelectorExpr:
		if b, ok := b.(*ast.SelectorExpr); ok {
			return a.Sel.Name == b.Sel.Name && matchesIdentOrSelector(a.X, b.X)
		}

		return false
	}

	return false
}

// render returns the pretty-print of the given node.
func render(fset *token.FileSet, x any) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}

	return buf.String()
}
