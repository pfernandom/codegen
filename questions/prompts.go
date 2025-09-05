package questions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PromptKey string

type Prompt interface {
	Prompt() (promptKey PromptKey, value string)
}

type stringPrompt struct {
	Key   PromptKey
	Label string
}

func StringPrompt(key string, label string) Prompt {
	return &stringPrompt{Key: PromptKey(key), Label: label}
}

func (p *stringPrompt) Prompt() (PromptKey, string) {
	return p.Key, sprompt(p.Label)
}

// StringPrompt asks for a string value using the label
func sprompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
