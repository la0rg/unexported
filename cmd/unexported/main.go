package unexported

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/la0rg/unexported"
)

func main() {
	singlechecker.Main(unexported.Analyzer)
}
