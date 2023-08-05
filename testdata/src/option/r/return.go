package r

type unexported struct{}

func UnexportedReturnParam() *unexported { return nil } // no-report; skip-func-args
