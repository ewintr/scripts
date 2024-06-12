package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputFile = flag.String("i", "", "input file")
	model     = flag.String("m", "llama3", "model file")
)

const (
	promptStart = `Could you translate the following text into English? Try te preserve the original tone of voice as much as possible. Only answer with the translation itself, no additional comments needed.`
)

func main() {
	flag.Parse()
	if *inputFile == "" {
		fmt.Println("No input file specified")
		os.Exit(1)
	}
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		fmt.Println("OLLAMA_HOST environment variable not set")
		os.Exit(1)
	}
	ollama := NewOllama(ollamaHost)

	doc, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	translated := make([]string, 0)
	for _, chunk := range strings.Split(string(doc), "\n\n") {
		prompt := fmt.Sprintf("%s\n---\n%s", promptStart, chunk)
		res, err := ollama.Generate(*model, prompt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf(".")
		translated = append(translated, res)
	}

	fmt.Printf("\n\n%s\n", strings.Join(translated, "\n"))
}
