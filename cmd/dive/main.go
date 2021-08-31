package main

import (
	"github.com/gostaticanalysis/dive"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(dive.Analyzer) }
