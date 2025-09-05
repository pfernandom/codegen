package main

import (
	"flag"
	"fmt"

	"github.com/pfernandom/codegen/questions"
)

func main() {
	questionFile := flag.String("questions", "questions.json", "a questions file")
	flag.Parse()
	if *questionFile != "" {
		fmt.Printf("Asking questions from %s\n", *questionFile)
		answers := questions.AskQuestionsFromFile(*questionFile)
		fmt.Println(answers)
		return
	}
	answers := questions.AskQuestions(
		questions.StringPrompt("name", "What is your name?"),
		questions.StringPrompt("age", "What is your age?"),
	)

	fmt.Println(answers)
}
