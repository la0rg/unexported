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
		// SkipInterfaces
		// RegexpFileSkip
		// SkipTypes
		// SkipFunctions
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

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil), (*ast.TypeSpec)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			a.analyzeFuncDecl(pass, n)
		case *ast.TypeSpec:
			a.analyzeTypeSpec(pass, n)
		}
	})
	return nil, nil
}

func (a *analyzer) analyzeFuncDecl(pass *analysis.Pass, f *ast.FuncDecl) {
	if !f.Name.IsExported() {
		return
	}

	// skip methods of unexported types
	if _, unexported := isUnexported(a.receiverType(pass, f)); unexported {
		return
	}

	description := fmt.Sprintf("function %s", f.Name)
	if f.Recv != nil {
		description = fmt.Sprintf("method %s", f.Name)
	}

	a.analyzeFieldList(pass, description, f.Type.Results)
	a.analyzeFieldList(pass, description, f.Type.Params)
}

func (a *analyzer) analyzeTypeSpec(pass *analysis.Pass, t *ast.TypeSpec) {
	if !t.Name.IsExported() {
		return
	}

	declType := pass.TypesInfo.TypeOf(t.Type)
	if typeName, unexported := isUnexported(declType); unexported {
		pass.Reportf(t.Pos(), "unexported type %s is used in the exported type declaration %s", typeName, t.Name)
	}
}

func (a *analyzer) analyzeFieldList(pass *analysis.Pass, description string, fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		fieldType := pass.TypesInfo.TypeOf(field.Type)
		if typeName, unexported := isUnexported(fieldType); unexported {
			pass.Reportf(field.Pos(), "unexported type %s is used in the exported %s", typeName, description)
		}
	}
}

func (a *analyzer) receiverType(pass *analysis.Pass, f *ast.FuncDecl) types.Type {
	if f.Recv == nil || len(f.Recv.List) == 0 {
		return nil
	}

	return pass.TypesInfo.TypeOf(f.Recv.List[0].Type)
}

func isUnexportedTypes(ts ...types.Type) (string, bool) {
	for _, t := range ts {
		if name, unexported := isUnexported(t); unexported {
			return name, true
		}
	}
	return "", false
}

func isUnexported(t types.Type) (string, bool) {
	switch T := t.(type) {
	case *types.Named:
		// skip builtins
		if T.Obj().Pkg() == nil {
			return "", false
		}

		return T.Obj().Name(), !T.Obj().Exported()

	case *types.Struct:
		var fields = make([]types.Type, 0, T.NumFields())
		for i := 0; i < T.NumFields(); i++ {
			if T.Field(i).Exported() {
				fields = append(fields, T.Field(i).Type())
			}
		}
		return isUnexportedTypes(fields...)

	case *types.Tuple:
		var vars = make([]types.Type, 0, T.Len())
		for i := 0; i < T.Len(); i++ {
			vars = append(vars, T.At(i).Type())
		}
		return isUnexportedTypes(vars...)

	case *types.Signature:
		return isUnexportedTypes(T.Params(), T.Results())

	case *types.Interface:
		var methods = make([]types.Type, 0, T.NumMethods())
		for i := 0; i < T.NumMethods(); i++ {
			if T.Method(i).Exported() {
				methods = append(methods, T.Method(i).Type())
			}
		}
		return isUnexportedTypes(methods...)

	case *types.Map:
		return isUnexportedTypes(T.Key(), T.Elem())

	case interface{ Elem() types.Type }:
		return isUnexported(T.Elem())
	}

	// otherwise assume it's exported to avoid false positives
	return "", false
}
