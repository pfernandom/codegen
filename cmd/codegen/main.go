package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/gen"
	"github.com/pfernandom/codegen/questions"
)

var Must = errors.MustNow

func main() {

	flags, closer := gen.ParseFlags()
	defer closer.Close()
	questionsHandler := questions.NewQuestionsHandler()
	if flags.IsGitRepo {
		fmt.Printf("Creating project from Git repo %s to %s\n", flags.GitRepo, flags.OutDir)
	} else {
		fmt.Printf("Creating project from %s and %s to %s\n", flags.DataFile, flags.TemplateDir, flags.OutDir)

	}
	dg := gen.NewDataGen(flags, questionsHandler)

	Must(os.RemoveAll(flags.OutDir))
	Must(os.MkdirAll(flags.OutDir, 0755))
	Must(dg.Generate())

	abs_out_dir, err := filepath.Abs(flags.OutDir)
	Must(err)
	fmt.Println("Generated project in", abs_out_dir)
	fmt.Println("done")
}
