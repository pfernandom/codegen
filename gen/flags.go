package gen

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	giturl "github.com/kubescape/go-git-url"
)

type CmdFlags struct {
	DataFile      string
	TemplateDir   string
	OutDir        string
	QuestionsFile string
	IsGitRepo     bool
	GitRepo       string
}

func ParseFlags() (CmdFlags, io.Closer) {
	data_file := flag.String("data", "data.json", "a data file")
	template_dir := flag.String("templates", "templates/", "a template directory")
	out_dir := flag.String("out", "out", "an output directory")
	questions_file := flag.String("questions", "questions.json", "a questions file")

	var is_git_repo bool = false
	var git_repo string = ""
	flag.Parse()

	base_dir := "."
	closer := io.Closer(nil)
	if isArgGitRepo() {
		is_git_repo = true
		git_repo = flag.Arg(0)
		var err error
		base_dir, closer, err = CloneGitRepo(flag.Arg(0))
		if err != nil {
			log.Fatalf("error cloning git repo: %s", err)
		}
		if _, err := os.Stat(base_dir); os.IsNotExist(err) {
			log.Fatalf("error cloning git repo: %s", err)
		}

	} else {
		base_dir = filepath.Dir(parseArg())
	}

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
		IsGitRepo:     is_git_repo,
		GitRepo:       git_repo,
	}, closer
}

func IsGitRepoTemplate(first_arg string) bool {
	if strings.HasPrefix(first_arg, "http") || strings.HasPrefix(first_arg, "git") {
		_, err := giturl.NewGitURL(first_arg)
		return err == nil
	}
	return false
}

func isArgGitRepo() bool {
	if len(flag.Args()) < 1 {
		return false
	}
	return IsGitRepoTemplate(flag.Args()[0])
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
