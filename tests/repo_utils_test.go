package tests

import (
	"path/filepath"
	"testing"

	"github.com/pfernandom/codegen/gen"
	"github.com/stretchr/testify/assert"
)

func TestIsGitRepoTemplate(t *testing.T) {
	assert.True(t, gen.IsGitRepoTemplate("git@github.com:pfernandom/codegen-example.git"))
	assert.True(t, gen.IsGitRepoTemplate("https://github.com/pfernandom/codegen-example.git"))
	assert.True(t, gen.IsGitRepoTemplate("git@github.com:pfernandom/codegen-example"))
	assert.True(t, gen.IsGitRepoTemplate("https://github.com/pfernandom/codegen-example"))
	assert.True(t, gen.IsGitRepoTemplate("git@github.com:pfernandom/codegen-example.git/"))
	assert.True(t, gen.IsGitRepoTemplate("https://github.com/pfernandom/codegen-example.git/"))
	assert.False(t, gen.IsGitRepoTemplate("example"))
	assert.False(t, gen.IsGitRepoTemplate("gitexample"))
}

func TestFormatGitRepoUrl(t *testing.T) {
	for _, repo := range []string{
		"https://github.com/pfernandom/codegen-example",
		"git@github.com:pfernandom/codegen-example.git",
		"git@github.com:pfernandom/codegen-example",
		"https://github.com/pfernandom/codegen-example.git",
		"https://github.com/pfernandom/codegen-example",
		"https://github.com/pfernandom/codegen-example/tree/otherbranch#",
	} {
		t.Run(repo, func(t *testing.T) {
			gitRepoUrl, _, err := gen.FormatGitRepoUrl(repo)
			assert.NoError(t, err)
			assert.NotEmpty(t, gitRepoUrl)
			assert.Equal(t, gitRepoUrl, "git@github.com:pfernandom/codegen-example.git")
		})
	}
}

func TestCloneGitRepo(t *testing.T) {
	repo := "https://github.com/pfernandom/codegen-example"
	tmpDir, closer, err := gen.CloneGitRepo(repo)
	defer closer.Close()
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpDir)
	assert.DirExists(t, tmpDir)
	assert.FileExists(t, filepath.Join(tmpDir, "data.json"))
	assert.DirExists(t, filepath.Join(tmpDir, "templates"))
	assert.FileExists(t, filepath.Join(tmpDir, "questions.json"))
}
