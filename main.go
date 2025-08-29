package main

import (
	"github.com/ajinbaby01/compoundassign/compoundassign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(compoundassign.Analyzer)
}
