package a

type unexported struct{}

func UnexportedFuncArg(a unexported) {} // no-report; skip-func-args
