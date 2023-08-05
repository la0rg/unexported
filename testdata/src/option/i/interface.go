package i

type unexported interface{ U() }

func UnexporterInterfaceFuncArgument(u unexported) {} // no report; skip-interfaces

func UnexporterInterfaceFuncReturnArgument() unexported { return nil } // no report; skip-interfaces

type UnexportedInterfaceTypeField struct{ U unexported } // no report; skip-interfaces
