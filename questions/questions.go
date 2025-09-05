package questions

import (
	"encoding/json"
	"errors"
	"os"

	e "github.com/pfernandom/codegen/errors"
)

var Must = e.MustNow

func AskQuestions(prompt Prompt, prompts ...Prompt) map[PromptKey]string {
	answers := make(map[PromptKey]string)
	key, value := prompt.Prompt()
	answers[key] = value
	for _, prompt := range prompts {
		key, value := prompt.Prompt()
		answers[key] = value
	}
	return answers
}

func AskQuestionsFromFile(file string) map[PromptKey]string {
	questions := make(map[PromptKey]string)
	b, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return questions
		}
		e.Must(err)
	}
	e.Must(json.Unmarshal(b, &questions))
	answers := make(map[PromptKey]string)

	for key, value := range questions {
		answers[key] = sprompt(value)
	}
	return answers
}
