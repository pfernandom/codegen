package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pfernandom/codegen/gen"
	"github.com/pfernandom/codegen/questions"
	"github.com/stretchr/testify/assert"
)

func getTestFlags(t *testing.T, base_dir string) gen.CmdFlags {
	return gen.CmdFlags{
		DataFile:      filepath.Join(base_dir, "data.json"),
		TemplateDir:   filepath.Join(base_dir, "templates"),
		OutDir:        filepath.Join(t.TempDir(), "out"),
		QuestionsFile: filepath.Join(base_dir, "questions.json"),
	}
}

type testQuestionsHandler struct {
	answers map[questions.PromptKey]string
}

func (h *testQuestionsHandler) AskQuestions(prompt questions.Prompt, prompts ...questions.Prompt) map[questions.PromptKey]string {
	return h.answers
}

func (h *testQuestionsHandler) AskQuestionsFromFile(file string) map[questions.PromptKey]string {
	return h.answers
}

type outputTester struct {
	*testing.T
	out_dir string
}

func (ot *outputTester) assertFileExistsAndContains(file string, content string) {
	ot.T.Helper()
	assert.FileExists(ot.T, file)
	data, err := os.ReadFile(file)
	assert.NoError(ot.T, err)
	assert.Contains(ot.T, string(data), content)
}

func (ot *outputTester) assertOutFileExistsAndContains(file string, content string) {
	ot.T.Helper()
	file = filepath.Join(ot.out_dir, file)
	assert.FileExists(ot.T, file)
	data, err := os.ReadFile(file)
	assert.NoError(ot.T, err)
	assert.Contains(ot.T, string(data), content)
}

func withExampleDir(t *testing.T, base_dir string, fn func(t *outputTester, flags gen.CmdFlags)) {
	flags := getTestFlags(t, base_dir)
	out_tester := &outputTester{T: t, out_dir: flags.OutDir}

	fn(out_tester, flags)
}

func (ot *outputTester) assertOutFileExists(file string) {
	ot.T.Helper()
	file = filepath.Join(ot.out_dir, file)
	assert.FileExists(ot.T, file)
}
