package unexported

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, NewAnalyzer(), "a", "r", "t")
}

func TestAnalyzerOptions(t *testing.T) {
	cases := []struct {
		pkg  string
		flag string
	}{
		{pkg: "option/t", flag: "skip-types"},
		{pkg: "option/a", flag: "skip-func-args"},
		{pkg: "option/r", flag: "skip-func-returns"},
		{pkg: "option/i", flag: "skip-interfaces"},
	}

	for _, c := range cases {
		t.Run(c.pkg, func(t *testing.T) {
			testdata := analysistest.TestData()
			analyzer := NewAnalyzer()
			if err := analyzer.Flags.Set(c.flag, "true"); err != nil {
				t.Fatal(err)
			}
			analysistest.Run(t, testdata, analyzer, c.pkg)
		})
	}
}
