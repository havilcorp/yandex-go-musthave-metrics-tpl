// Package fmtprint анализатор для проверки fmt.Print
package fmtprint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "fmtPrint",
	Doc:  "check fmt.Print in code",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "fmt" {
				if fun.Sel.Name == "Print" {
					pass.Reportf(callExpr.Pos(), "calling fmt.Print, use logrus")
				}
				if fun.Sel.Name == "Println" {
					pass.Reportf(callExpr.Pos(), "calling fmt.Println, use logrus")
				}
			}

			return false
		})
	}
	return nil, nil
}
