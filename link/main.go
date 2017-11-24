package main

import (
	"fmt"
	"os"
)

import (
	"github.com/ajm188/gophercise/link/link"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println(link.FindLinks(file))
}
