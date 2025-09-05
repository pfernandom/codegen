package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/pfernandom/codegen/gen"
	"github.com/pfernandom/codegen/questions"
	"github.com/stretchr/testify/assert"
)

func TestCodegen(t *testing.T) {
	withExampleDir(t, "example", func(t *outputTester, flags gen.CmdFlags) {
		dg := gen.NewDataGen(flags, &testQuestionsHandler{
			answers: map[questions.PromptKey]string{
				"name": "John Doe",
				"age":  "30",
			},
		})

		assert.NoError(t, dg.Validate())
		assert.NoError(t, dg.Generate())
		out_dir_contents, err := os.ReadDir(flags.OutDir)
		assert.NoError(t, err)
		fmt.Println(out_dir_contents)

		t.assertOutFileExistsAndContains(
			"index.html",
			"<title>This is a title</title>",
		)
		t.assertOutFileExists(
			"js/my_app.js",
		)
		t.assertOutFileExistsAndContains(
			"js/index.js",
			"console.log(\"Hello John Doe, World!\");",
		)
	})
}
