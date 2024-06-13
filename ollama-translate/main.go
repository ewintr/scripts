package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputFile  = flag.String("i", "", "input file (markdown, or plain text)")
	model      = flag.String("m", "llama3", "llm model")
	outputFile = flag.String("o", "", "output file")
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

	chunks := make([]string, 0)
	for _, par := range strings.Split(string(doc), "\n\n") {
		last := len(chunks) - 1
		switch {
		case last == -1:
			chunks = append(chunks, par)
		case last >= 0 && len(par)+len(chunks[last]) > 500:
			chunks = append(chunks, par)
		default:
			chunks[last] = fmt.Sprintf("%s\n\n%s", chunks[last], par)
		}
	}

	fmt.Printf("translating %d chunks\n", len(chunks))
	translated := make([]string, 0)
	for _, chunk := range chunks {
		prompt := fmt.Sprintf("%s\n---\n%s", promptStart, chunk)
		res, err := ollama.Generate(*model, prompt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf(".")
		translated = append(translated, res)
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, []byte(strings.Join(translated, "\n")), 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nfile written to %s\n", *outputFile)
		os.Exit(0)
	}

	fmt.Printf("\n\n%s\n", strings.Join(translated, "\n"))
}
