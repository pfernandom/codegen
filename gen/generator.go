package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	. "github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/questions"
)

var LOG = slog.New(slog.NewTextHandler(os.Stderr, nil))

type dataGen struct {
	Data        map[string]string
	OutDir      string
	TemplateDir string
}

func NewDataGen(flags CmdFlags, questionsHandler questions.QuestionsHandler) *dataGen {
	data := readJsonFile(flags.DataFile)
	for key, value := range questionsHandler.AskQuestionsFromFile(flags.QuestionsFile) {
		data[string(key)] = value
	}
	return &dataGen{
		Data:        data,
		OutDir:      flags.OutDir,
		TemplateDir: flags.TemplateDir,
	}
}

func (dg *dataGen) Validate() error {
	if _, err := os.Stat(dg.TemplateDir); os.IsNotExist(err) {
		return fmt.Errorf("template directory not found: %s", dg.TemplateDir)
	}
	return nil
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

func (dg *dataGen) ParseToOut(s string, out_file string, d fs.FileInfo, err error) error {
	Must(os.MkdirAll(filepath.Dir(out_file), 0755)).OrFail()

	f := MustReturn(os.Create(out_file)).OrFailWith("create file: %w", out_file)
	defer f.Close()
	if strings.HasSuffix(s, ".tmpl") {
		tmpl := MustReturn(template.ParseFiles(s)).OrFailWith("parse template: %w", s)
		Must(tmpl.Execute(f, dg.Data)).OrFail()
	} else if !d.IsDir() {
		inf := MustReturn(os.Open(s)).OrFailWith("open out file: %w", s)
		defer inf.Close()
		MustReturn(io.Copy(f, inf)).OrFailWith("copy file: %w", s)
	}
	return nil
}
func (dg *dataGen) Walk(s string, d fs.FileInfo, err error) error {
	MustOrMessage(err, "walk error: %w", err)
	if strings.Contains(s, ".git") {
		return nil
	}
	if d.IsDir() {
		return nil
	}

	out_file := MustReturn(parseFilenameTemplate(s, dg.Data)).OrFail()
	out_file = MustReturn(filepath.Rel(dg.TemplateDir, out_file)).OrFail()
	out_file = filepath.Join(dg.OutDir, strings.TrimSuffix(out_file, ".tmpl"))

	return dg.ParseToOut(s, out_file, d, err)
}

func parseFilenameTemplate(s string, data map[string]string) (string, error) {
	st, err := template.New(fmt.Sprintf("filename-%s", s)).Parse(s)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := st.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
