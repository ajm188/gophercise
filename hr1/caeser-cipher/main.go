package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func caeserCipher(k int32) func(rune) rune {
	return func(c rune) rune {
		switch {
		case c >= 'A' && c <= 'Z':
			return 'A' + (c-'A'+k)%26
		case c >= 'a' && c <= 'z':
			return 'a' + (c-'a'+k)%26
		default:
			return c
		}
	}
}

func encrypt(plaintext string, k int32) string {
	return strings.Map(caeserCipher(k), plaintext)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	defer os.Stdin.Close()
	rotationAmount, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	plaintext, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	k, err := strconv.Atoi(rotationAmount[:len(rotationAmount)-1])
	if err != nil {
		panic(err)
	}
	fmt.Println(encrypt(plaintext[:len(plaintext)-1], int32(k)))
}
