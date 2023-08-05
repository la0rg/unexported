package t

type unexported struct{}

type Exported struct{}

type ExportedStructField struct {
	E Exported
}

type UnexportedStructField struct { // want `unexported`
	U unexported
}

type UnexportedFieldDoesNotTriggerReport struct {
	u unexported
}

type ExportedMapValue map[int]Exported

type UnexportedMapValue map[int]unexported // want `unexported`

type ExportedMapKey map[Exported]string

type UnexportedMapKey map[unexported]string // want `unexported`

type ExportedSlice []Exported

type UnexportedSlice []unexported // want `unexported`

type ExportedFunctionArg func(exported Exported)

type UnexportedFunctionArg func(unexported) // want `unexported`

type ExportedFunctionReturn func() Exported

type UnexportedFunctionReturn func() unexported // want `unexported`

type UnexportedInterface interface { // want `unexported`
	U() unexported
}

type ExportedInterface interface {
	E() Exported
}

type UnexportedMethodShouldNotTriggerReport interface {
	u() unexported
}

type unexportedTypeDoesNotTriggerReport struct {
	U unexported
}
