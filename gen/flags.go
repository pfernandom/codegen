package gen

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CmdFlags struct {
	DataFile      string
	TemplateDir   string
	OutDir        string
	QuestionsFile string
}

func ParseFlags() CmdFlags {

	data_file := flag.String("data", "data.json", "a data file")
	template_dir := flag.String("templates", "templates/", "a template directory")
	out_dir := flag.String("out", "out", "an output directory")
	questions_file := flag.String("questions", "questions.json", "a questions file")
	flag.Parse()

	base_dir := filepath.Dir(parseArg())

	fmt.Printf("base_dir: %s\n", base_dir)

	*data_file = filepath.Join(base_dir, *data_file)
	*template_dir = filepath.Join(base_dir, *template_dir)
	*questions_file = filepath.Join(base_dir, *questions_file)
	// *out_dir = filepath.Join(base_dir, *out_dir)

	if _, err := os.Stat(*data_file); os.IsNotExist(err) {
		log.Fatalf("error reading template data: data file not found: %s", *data_file)
	}
	if _, err := os.Stat(*template_dir); os.IsNotExist(err) {
		log.Fatalf("error reading template data: template directory not found: %s", *template_dir)
	}

	return CmdFlags{
		DataFile:      *data_file,
		TemplateDir:   *template_dir,
		OutDir:        *out_dir,
		QuestionsFile: *questions_file,
	}
}

func parseArg() string {
	if len(flag.Args()) < 1 {
		return "."
	}
	first_arg := flag.Args()[0]
	if !strings.HasSuffix(first_arg, "/") {
		return first_arg + "/"
	}
	return first_arg
}
