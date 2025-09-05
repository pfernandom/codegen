package main

import (
	"flag"
	"fmt"

	"github.com/pfernandom/codegen/questions"
)

func main() {
	questionFile := flag.String("questions", "questions.json", "a questions file")
	flag.Parse()
	questionsHandler := questions.NewQuestionsHandler()
	if *questionFile != "" {
		fmt.Printf("Asking questions from %s\n", *questionFile)
		answers := questionsHandler.AskQuestionsFromFile(*questionFile)
		fmt.Println(answers)
		return
	}
	answers := questionsHandler.AskQuestions(
		questions.StringPrompt("name", "What is your name?"),
		questions.StringPrompt("age", "What is your age?"),
	)

	fmt.Println(answers)
}
