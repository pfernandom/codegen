package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/questions"
)

type dataGen struct {
	Data        map[string]string
	OutDir      string
	TemplateDir string
}

func NewDataGen(flags CmdFlags) *dataGen {
	data := readJsonFile(flags.DataFile)
	for key, value := range questions.AskQuestionsFromFile(flags.QuestionsFile) {
		data[string(key)] = value
	}
	return &dataGen{
		Data:        data,
		OutDir:      flags.OutDir,
		TemplateDir: flags.TemplateDir,
	}
}

func readJsonFile(file string) map[string]string {
	var data map[string]string
	b, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return data
		}
		Must(err)
	}
	Must(json.Unmarshal(b, &data)).OrFail()
	return data
}

func (dg *dataGen) Generate() error {
	return filepath.Walk(dg.TemplateDir, dg.Walk)
}

func (dg *dataGen) Walk(s string, d fs.FileInfo, err error) error {
	MustOrMessage(err, "walk error: %w", err)
	st := MustReturn(template.New(s).Parse(s)).OrFail()

	var buf bytes.Buffer
	Must(st.Execute(&buf, dg.Data)).OrFail()

	out_file := s
	if buf.String() != s {
		out_file = buf.String()
	}

	if d.IsDir() {
		return nil
	}
	out_file = strings.TrimPrefix(out_file, "templates/")
	out_file = filepath.Join(dg.OutDir, strings.TrimSuffix(out_file, ".tmpl"))
	Must(os.MkdirAll(filepath.Dir(out_file), 0755)).OrFail()

	f := MustReturn(os.Create(out_file)).OrFailWith("create file: %w", out_file)
	defer f.Close()
	if strings.HasSuffix(s, ".tmpl") {
		tmpl := MustReturn(template.ParseFiles(s)).OrFailWith("parse template: %w", s)
		Must(tmpl.Execute(f, dg.Data)).OrFail()

		err = tmpl.Execute(f, dg.Data)
		Must(err).OrFailWith("execute template: %w", err)
	} else if !d.IsDir() {
		inf := MustReturn(os.Open(s)).OrFailWith("open out file: %w", s)
		defer inf.Close()
		MustReturn(io.Copy(f, inf)).OrFailWith("copy file: %w", s)
	}
	return nil
}
