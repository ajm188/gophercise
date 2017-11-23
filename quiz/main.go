package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
)

type Problem struct {
	question string
	answer string
}

func parseProblems(rawText string) []Problem {
	numLines := strings.Count(rawText, "\n")
	problems := make([]Problem, numLines)
	for i := 0; i < numLines; i++ {
		newline := strings.Index(rawText, "\n")
		line := rawText[:newline]

		comma := strings.LastIndex(line, ",")
		problem := Problem{
			question: line[:comma],
			answer: line[comma + 1:],
		}

		rawText = rawText[newline + 1:]
		problems[i] = problem
	}

	return problems
}

func quiz(problems []Problem) int {
	score := 0
	var answer string
	for _, problem := range problems {
		fmt.Printf("%s: ", problem.question)
		fmt.Scanf("%s\n", &answer)
		if answer == problem.answer {
			score++
		}
	}

	return score
}

func main() {
	problemFilePath := flag.String("problem-file", "problems.csv", "Path to CSV problems file")
	flag.Parse()

	problemBytes, err := ioutil.ReadFile(*problemFilePath)
	if err != nil {
		return
	}
	problems := parseProblems(string(problemBytes))

	score := quiz(problems)
	fmt.Println(score)
}
