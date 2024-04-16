// Package osexit анализатор для проверки os.Exit в пакете main в функции main
package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osExit",
	Doc:  "check calling os.Exit",
	Run:  run,
}

// run функция анализатора
func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	// исключаем вызов os.Exit в caches
	for _, pkg := range pass.Pkg.Imports() {
		if pkg.Name() == "testing" {
			return nil, nil
		}
	}
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			return nil, nil
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if fd, ok := n.(*ast.FuncDecl); ok && fd.Name.Name != "main" {
				return false
			}

			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" {
				pass.Reportf(callExpr.Pos(), "calling os.Exit")
			}

			return true
		})
	}
	return nil, nil
}
