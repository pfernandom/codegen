package main

import (
	"fmt"
	"os"

	"github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/gen"
)

var Must = errors.MustNow

func main() {

	flags := gen.ParseFlags()

	fmt.Printf("Creating project from %s and %s to %s\n", flags.DataFile, flags.TemplateDir, flags.OutDir)

	dg := gen.NewDataGen(flags)

	Must(os.RemoveAll(flags.OutDir))
	Must(os.MkdirAll(flags.OutDir, 0755))
	Must(dg.Generate())
	fmt.Println("done")
}
