package a

type unexported struct{}

type Exported struct{}

func UnexportedType(a unexported) {} // want `unexported type`

func ExportedType(a Exported) {}

func UnexportedPointer(a *unexported) {} // want `unexported type`

func ExportedPointer(*Exported) {}

func UnexportedStruct(struct{ U unexported }) {} // want `unexported type`

func ExportedStruct(struct{ e Exported }) {}

func UnexportedFunctionType(func() unexported) {} // want `unexported type`

func ExportedFunctionType(func() Exported) {}

func UnexportedFunctionArgType(func(u unexported)) {} // want `unexported type`

func ExportedFunctionArgType(func(exported Exported)) {}

func UnexportedInterface(interface{ E() unexported }) {} // want `unexported type`

func ExportedInterface(interface{ e() Exported }) {}

func UnexportedMapKey(map[unexported]Exported) {} // want `unexported type`

func UnexportedMapValue(map[Exported]unexported) {} // want `unexported type`

func ExportedMap(map[Exported]Exported) {}

func unexportedFunctionsAreSkipped(unexported) {}

func (Exported) unexportedMethodsAreSkipped(unexported) {}

func (Exported) UnexportedType(unexported) {} // want `unexported type`

func (Exported) unexportedType(unexported) {} // no report for private methods

func (Exported) ExportedType(Exported) {}

func (unexported) UnexportedType(unexported) {} // no report for methods of unexported types

func (unexported) ExportedType(Exported) {}
