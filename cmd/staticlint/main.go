package main

import (
	ast "go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

const (
	mainFunc = "main"
	exitFunc = "Exit"
	osPkg    = "os"
)

// exitInMainAnalyzer проверяет использование os.Exit.
var exitInMainAnalyzer = &analysis.Analyzer{
	Name: "exitinmain",
	Doc:  "проверяет использование os.Exit в функции main пакета main",
	Run:  runExitInMainAnalyzer,
}

func main() {
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	checks := map[string]bool{
		"S1000": true,
		"S1001": true,
	}

	for _, v := range simple.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, shift.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, exitInMainAnalyzer)

	multichecker.Main(
		mychecks...,
	)
}

// runExitInMainAnalyzer функция проверки использования os.Exit.
func runExitInMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != mainFunc || pass.Pkg.Name() != mainFunc {
				continue
			}

			ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
				callExpr, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selExpr.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == osPkg && selExpr.Sel.Name == exitFunc {
					pass.Reportf(callExpr.Pos(), "использование os.Exit в функции main пакета main запрещено")
				}

				return true
			})
		}
	}
	return nil, nil
}
