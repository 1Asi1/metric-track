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

func runExitInMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "main" && pass.Pkg.Name() == "main" {
					ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
						if callExpr, ok := node.(*ast.CallExpr); ok {
							if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								if ident, ok := selExpr.X.(*ast.Ident); ok {
									if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
										pass.Reportf(callExpr.Pos(), "использование os.Exit в функции main пакета main запрещено")
									}
								}
							}
						}
						return true
					})
				}
			}
		}
	}
	return nil, nil
}
