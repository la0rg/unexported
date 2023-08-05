package t

type unexported struct{}

type UnexportedType struct{ U unexported } // no report; skip-types
