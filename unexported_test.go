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
	testdata := analysistest.TestData()

	analyzer := NewAnalyzer()
	analyzer.Flags.Set("skip-interfaces", "true")

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "o")
}
