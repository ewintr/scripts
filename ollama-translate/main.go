package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	inputFile = flag.String("i", "", "input file")
	model     = flag.String("m", "llama3", "model file")
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
	res, err := ollama.Generate(*model, "Could you translate the following text into English? 'Mijn fietsband is lek. Wat moet ik nu doen'")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(res)
}
