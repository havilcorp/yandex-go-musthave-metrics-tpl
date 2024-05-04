// Package commentcode анализатор для проверки заккоментированного кода
package commentcode

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "commentCode",
	Doc:  "check comment code",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, c := range file.Comments {
			if strings.HasPrefix(c.Text(), "func") {
				pass.Reportf(c.Pos(), "delete comment code")
			}
		}
	}
	return nil, nil
}
