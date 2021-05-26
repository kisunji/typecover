package typecover_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"typecover"
)

func TestStructs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), typecover.Analyzer, "structs")
}

func TestInterfaces(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), typecover.Analyzer, "interfaces")
}
