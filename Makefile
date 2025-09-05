codegen = bin/codegen
questions = bin/questions

.PHONY: install

$(codegen):
	go build -o $(codegen) cmd/codegen/main.go

$(questions):
	go build -o $(questions) cmd/questions/main.go

clean:
	rm -f $(codegen) $(questions)

install:
	go install ./cmd/...

all: $(codegen) $(questions)