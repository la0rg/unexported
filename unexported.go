package unexported

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "unexported",
	Doc:      "check that exported functions do not accept/return unexported types",
	URL:      "https://github.com/la0rg/unexported",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		f := n.(*ast.FuncDecl)
		if !f.Name.IsExported() {
			return
		}

		// skip methods of unexported types
		if isUnexported(receiverType(pass, f)) {
			return
		}

		if f.Type.Results != nil {
			for _, ret := range f.Type.Results.List {
				retType := pass.TypesInfo.TypeOf(ret.Type)
				if isUnexported(retType) {
					pass.Reportf(ret.Pos(), "unexported type %s is used in the exported function %s", retType, f.Name)
				}
			}
		}
	})
	return nil, nil
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

func receiverType(pass *analysis.Pass, f *ast.FuncDecl) types.Type {
	if f.Recv == nil || len(f.Recv.List) == 0 {
		return nil
	}

	return pass.TypesInfo.TypeOf(f.Recv.List[0].Type)
}
