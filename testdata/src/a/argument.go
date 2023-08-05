package a

type unexported struct{}

type Exported struct{}

func UnexportedType(a unexported) {} // want `unexported`

func ExportedType(a Exported) {}

func UnexportedPointer(a *unexported) {} // want `unexported`

func ExportedPointer(*Exported) {}

func UnexportedStruct(struct{ U unexported }) {} // want `unexported`

func ExportedStruct(struct{ e Exported }) {}

func UnexportedFunctionType(func() unexported) {} // want `unexported`

func ExportedFunctionType(func() Exported) {}

func UnexportedFunctionArgType(func(u unexported)) {} // want `unexported`

func ExportedFunctionArgType(func(exported Exported)) {}

func UnexportedInterface(interface{ E() unexported }) {} // want `unexported`

func ExportedInterface(interface{ e() Exported }) {}

func UnexportedMapKey(map[unexported]Exported) {} // want `unexported`

func UnexportedMapValue(map[Exported]unexported) {} // want `unexported`

func ExportedMap(map[Exported]Exported) {}

func unexportedFunctionsAreSkipped(unexported) {}

func (Exported) unexportedMethodsAreSkipped(unexported) {}

func (Exported) UnexportedType(unexported) {} // want `unexported`

func (Exported) unexportedType(unexported) {} // no report for private methods

func (Exported) ExportedType(Exported) {}

func (unexported) UnexportedType(unexported) {} // no report for methods of unexported types

func (unexported) ExportedType(Exported) {}

type unexportedInterfaceType interface {
	Method()
}

func UnexportedNamedInterface(unexportedInterfaceType) {} // want "unexported"

type ExportedInterfaceType interface {
	Exported() unexported
}

func ExportedNamedInterface(ExportedInterfaceType) {}
