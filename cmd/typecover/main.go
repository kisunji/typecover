package main

import (
	"github.com/kisunji/typecover"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(typecover.Analyzer) }
