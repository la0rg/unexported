package unexported

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a new analyzer that checks that exported functions and types use only exported types in their signatures.
func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analyzer{
		// TODO add skip settings
	}
	return &analysis.Analyzer{
		Name:     "unexported",
		Doc:      "check that exported functions do not accept/return unexported types",
		URL:      "https://github.com/la0rg/unexported",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      analyzer.run,
	}
}

type analyzer struct{}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			a.analyzeFuncDecl(pass, n)
		}
	})
	return nil, nil
}

func (a *analyzer) analyzeFuncDecl(pass *analysis.Pass, f *ast.FuncDecl) {
	if !f.Name.IsExported() {
		return
	}

	// skip methods of unexported types
	if isUnexported(a.receiverType(pass, f)) {
		return
	}

	description := fmt.Sprintf("function %s", f.Name)
	if f.Recv != nil {
		description = fmt.Sprintf("method %s", f.Name)
	}

	a.analyzeFieldList(pass, description, f.Type.Results)
	a.analyzeFieldList(pass, description, f.Type.Params)
	a.analyzeFieldList(pass, description, f.Type.TypeParams)
}

func (a *analyzer) analyzeFieldList(pass *analysis.Pass, description string, fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		fieldType := pass.TypesInfo.TypeOf(field.Type)
		if isUnexported(fieldType) {
			pass.Reportf(field.Pos(), "unexported %s is used in the exported %s", fieldType, description)
		}
	}
}

func (a *analyzer) receiverType(pass *analysis.Pass, f *ast.FuncDecl) types.Type {
	if f.Recv == nil || len(f.Recv.List) == 0 {
		return nil
	}

	return pass.TypesInfo.TypeOf(f.Recv.List[0].Type)
}

func isUnexported(t types.Type) bool {
	switch T := t.(type) {
	case *types.Named:
		// skip builtins
		if T.Obj().Pkg() == nil {
			return false
		}

		return !T.Obj().Exported()

	case *types.Struct:
		for i := 0; i < T.NumFields(); i++ {
			if isUnexported(T.Field(i).Type()) {
				return true
			}
		}

	case *types.Tuple:
		for i := 0; i < T.Len(); i++ {
			if isUnexported(T.At(i).Type()) {
				return true
			}
		}

	case *types.Signature:
		return isUnexported(T.Params()) || isUnexported(T.Results())

	case *types.Interface:
		for i := 0; i < T.NumMethods(); i++ {
			if isUnexported(T.Method(i).Type()) {
				return true
			}
		}

	case *types.Map:
		return isUnexported(T.Key()) || isUnexported(T.Elem())

	case interface{ Elem() types.Type }:
		return isUnexported(T.Elem())
	}

	// otherwise assume it's exported to avoid false positives
	return false
}
