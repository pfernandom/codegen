package questions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	e "github.com/pfernandom/codegen/errors"
)

var Must = e.MustNow

type QuestionsHandler interface {
	AskQuestions(prompt Prompt, prompts ...Prompt) map[PromptKey]string
	AskQuestionsFromFile(file string) map[PromptKey]string
}

type questionsHandler struct {
}

func NewQuestionsHandler() QuestionsHandler {
	return &questionsHandler{}
}

func (h *questionsHandler) AskQuestions(prompt Prompt, prompts ...Prompt) map[PromptKey]string {
	answers := make(map[PromptKey]string)
	key, value := prompt.Prompt()
	answers[key] = value
	for _, prompt := range prompts {
		key, value := prompt.Prompt()
		answers[key] = value
	}
	return answers
}

func (h *questionsHandler) AskQuestionsFromFile(file string) map[PromptKey]string {
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

	if len(questions) == 0 {
		return answers
	}
	fmt.Println("Please answer the following questions:")
	for key, value := range questions {
		answers[key] = sprompt(value)
	}
	return answers
}
