package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func main() {

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {

		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText))

	}

}
