package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func isCapital(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func countWords(text string) int {
	words := strings.FieldsFunc(text, isCapital)
	return len(words)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	defer os.Stdin.Close()
	for {
		text, err := reader.ReadString('\n')
		switch {
		case err == io.EOF:
			return
		case err != nil:
			panic(err)
		}
		fmt.Println(countWords(text))
	}
}
