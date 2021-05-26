package typecover_test

import (
	"testing"

	"github.com/kisunji/typecover"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestStructs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), typecover.Analyzer, "structs")
}

func TestInterfaces(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), typecover.Analyzer, "interfaces")
}
