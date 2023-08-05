package unexported

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a new analyzer that checks that exported functions and types use only exported types in their signatures.
func NewAnalyzer() *analysis.Analyzer {
	opts := new(options)
	flags := flag.FlagSet{}
	flags.BoolVar(&opts.SkipInterfaces, "skip-interfaces", false, "when set, interfaces are excluded from analysis")
	flags.BoolVar(&opts.SkipTypes, "skip-types", false, "when set, types are excluded from analysis")
	flags.BoolVar(&opts.SkipFuncArgs, "skip-func-args", false, "when set, function arguments are excluded from analysis")
	flags.BoolVar(&opts.SkipFuncReturns, "skip-func-returns", false, "when set, function return parameters are excluded from analysis")

	return &analysis.Analyzer{
		Name:     "unexported",
		Doc:      "check that exported functions do not accept/return unexported types",
		URL:      "https://github.com/la0rg/unexported",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run(opts),
		Flags:    flags,
	}
}

type options struct {
	SkipInterfaces  bool
	SkipTypes       bool
	SkipFuncArgs    bool
	SkipFuncReturns bool
}

func run(opts *options) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		analyzer := &analyzer{pass: pass, opts: opts}

		nodeFilter := []ast.Node{(*ast.FuncDecl)(nil), (*ast.TypeSpec)(nil)}
		inspector.Preorder(nodeFilter, func(n ast.Node) {
			switch n := n.(type) {
			case *ast.FuncDecl:
				analyzer.funcDecl(n)
			case *ast.TypeSpec:
				analyzer.typeSpec(n)
			}
		})
		return nil, nil
	}
}

type analyzer struct {
	pass *analysis.Pass
	opts *options
}

func (a analyzer) funcDecl(f *ast.FuncDecl) {
	if !f.Name.IsExported() {
		return
	}

	// skip methods of unexported types
	if _, unexported := a.isUnexported(a.receiverType(f)); unexported {
		return
	}

	description := fmt.Sprintf("function %s", f.Name)
	if f.Recv != nil {
		description = fmt.Sprintf("method %s", f.Name)
	}

	if !a.opts.SkipFuncArgs {
		a.fieldList(description, f.Type.Params)
	}

	if !a.opts.SkipFuncReturns {
		a.fieldList(description, f.Type.Results)
	}
}

func (a analyzer) typeSpec(t *ast.TypeSpec) {
	if !t.Name.IsExported() || a.opts.SkipTypes {
		return
	}

	declType := a.pass.TypesInfo.TypeOf(t.Type)
	if typeName, unexported := a.isUnexported(declType); unexported {
		a.pass.Reportf(t.Pos(), "unexported type %s is used in the exported type declaration %s", typeName, t.Name)
	}
}

func (a *analyzer) fieldList(description string, fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		fieldType := a.pass.TypesInfo.TypeOf(field.Type)
		if typeName, unexported := a.isUnexported(fieldType); unexported {
			a.pass.Reportf(field.Pos(), "unexported type %s is used in the exported %s", typeName, description)
		}
	}
}

func (a analyzer) receiverType(f *ast.FuncDecl) types.Type {
	if f.Recv == nil || len(f.Recv.List) == 0 {
		return nil
	}

	return a.pass.TypesInfo.TypeOf(f.Recv.List[0].Type)
}

func (a analyzer) isUnexportedTypes(ts ...types.Type) (string, bool) {
	for _, t := range ts {
		if name, unexported := a.isUnexported(t); unexported {
			return name, true
		}
	}
	return "", false
}

func (a analyzer) isUnexported(t types.Type) (string, bool) {
	switch T := t.(type) {
	case *types.Named:
		// skip builtins
		if T.Obj().Pkg() == nil {
			return "", false
		}

		if _, isInterface := T.Underlying().(*types.Interface); isInterface && a.opts.SkipInterfaces {
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
		return a.isUnexportedTypes(fields...)

	case *types.Tuple:
		var vars = make([]types.Type, 0, T.Len())
		for i := 0; i < T.Len(); i++ {
			vars = append(vars, T.At(i).Type())
		}
		return a.isUnexportedTypes(vars...)

	case *types.Signature:
		return a.isUnexportedTypes(T.Params(), T.Results())

	case *types.Interface:
		var methods = make([]types.Type, 0, T.NumMethods())
		for i := 0; i < T.NumMethods(); i++ {
			if T.Method(i).Exported() {
				methods = append(methods, T.Method(i).Type())
			}
		}
		return a.isUnexportedTypes(methods...)

	case *types.Map:
		return a.isUnexportedTypes(T.Key(), T.Elem())

	case interface{ Elem() types.Type }:
		return a.isUnexported(T.Elem())
	}

	// otherwise assume it's exported to avoid false positives
	return "", false
}
