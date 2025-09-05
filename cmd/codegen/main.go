package main

import (
	"fmt"
	"os"

	"github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/gen"
	"github.com/pfernandom/codegen/questions"
)

var Must = errors.MustNow

func main() {

	flags := gen.ParseFlags()
	questionsHandler := questions.NewQuestionsHandler()
	fmt.Printf("Creating project from %s and %s to %s\n", flags.DataFile, flags.TemplateDir, flags.OutDir)

	dg := gen.NewDataGen(flags, questionsHandler)

	Must(os.RemoveAll(flags.OutDir))
	Must(os.MkdirAll(flags.OutDir, 0755))
	Must(dg.Generate())
	fmt.Println("done")
}
