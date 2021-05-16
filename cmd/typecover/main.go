package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"typecover"
)

func main() { singlechecker.Main(typecover.Analyzer) }
