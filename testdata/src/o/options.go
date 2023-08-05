package o

type unexported interface {
	U()
}

func UnexporterInterfaceFuncArgument(u unexported) {} // no report; skip-interfaces

func UnexporterInterfaceFuncReturnArgument(u unexported) {} // no report; skip-interfaces

type UnexportedInterfaceField struct { // no report; skip-interfaces
	U unexported
}
