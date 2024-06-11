package main

import (
	"flag"
	"fmt"
)

var (
	inputFile = flag.String("i", "", "input file")
)

func main() {
	flag.Parse()

	fmt.Println("Hello World")
}
