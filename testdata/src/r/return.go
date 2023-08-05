package r

type unexported struct{}

type Exported struct{}

func NoReturn() {}

func UnexportedType() unexported { // want `unexported type`
	return unexported{}
}

func ExportedType() Exported {
	return Exported{}
}

func UnexportedPointer() *unexported { // want `unexported type`
	return nil
}

func ExportedPointer() *Exported {
	return nil
}

func UnexportedStruct() struct{ U unexported } { // want `unexported type`
	return struct{ U unexported }{}
}

func ExportedStruct() struct{ e Exported } {
	return struct{ e Exported }{}
}

func UnexportedTuple() (unexported, error) { // want `unexported type`
	return unexported{}, nil
}

func ExportedTuple() (Exported, error) {
	return Exported{}, nil
}

func UnexportedFunctionType() func() unexported { // want `unexported type`
	return nil
}

func ExportedFunctionType() func() Exported {
	return nil
}

func UnexportedFunctionArgType() func(u unexported) { // want `unexported type`
	return nil
}

func ExportedFunctionArgType() func(exported Exported) {
	return nil
}

func UnexportedInterface() interface{ E() unexported } { // want `unexported type`
	return nil
}

func ExportedInterface() interface{ e() Exported } {
	return nil
}

func UnexportedMapKey() map[unexported]Exported { // want `unexported type`
	return nil
}

func UnexportedMapValue() map[Exported]unexported { // want `unexported type`
	return nil
}

func ExportedMap() map[Exported]Exported {
	return nil
}

func unexportedFunctionsAreSkipped() unexported {
	return unexported{}
}

func (Exported) unexportedMethodsAreSkipped() unexported {
	return unexported{}
}

func (Exported) UnexportedType() unexported { // want `unexported type`
	return unexported{}
}

func (Exported) unexportedType() unexported { // no report for private methods
	return unexported{}
}

func (Exported) ExportedType() Exported {
	return Exported{}
}

func (unexported) UnexportedType() unexported { // no report for methods of unexported types
	return unexported{}
}

func (unexported) ExportedType() Exported { // no report for methods of unexported types
	return Exported{}
}
