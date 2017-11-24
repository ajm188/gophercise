package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
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
			answer:   line[comma+1:],
		}

		rawText = rawText[newline+1:]
		problems[i] = problem
	}

	return problems
}

func quiz(problems []Problem, timeLimit int) int {
	reader := bufio.NewReader(os.Stdin)
	problemChan := make(chan Problem, len(problems))
	scoreChan := make(chan int)

	go func(problems chan Problem, scores chan int) {
		for {
			problem := <-problems
			question := problem.question
			fmt.Printf("%s: ", question)
			answer, _ := reader.ReadString('\n')
			answer = answer[:len(answer)-1]
			score := 0
			if problem.answer == answer {
				score = 1
			}
			scores <- score
		}
	}(problemChan, scoreChan)
	defer close(scoreChan)
	defer close(problemChan)
	defer os.Stdin.Close()

	fmt.Println("Press enter to start")
	_, err := reader.ReadString('\n')
	if err != nil {
		return 0
	}

	stopTime := time.Now().Add(time.Duration(timeLimit) * time.Second)
	for _, problem := range problems {
		problemChan <- problem
	}
	score := 0
	for {
		if !stopTime.After(time.Now()) {
			return score
		}
		select {
		case nextScore := <-scoreChan:
			score += nextScore
		default:
			continue
		}
	}
	return score
}

func main() {
	problemFilePath := flag.String("problem-file", "problems.csv", "Path to CSV problems file")
	timeLimit := flag.Int("time-limit", 30, "Time limit in seconds.")
	flag.Parse()

	problemBytes, err := ioutil.ReadFile(*problemFilePath)
	if err != nil {
		return
	}
	problems := parseProblems(string(problemBytes))

	score := quiz(problems, *timeLimit)
	fmt.Printf("\nFinal score: %d / %d\n", score, len(problems))
}
