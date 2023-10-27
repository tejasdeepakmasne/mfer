package main

import (
	"fmt"
	"os"

	lexer "github.com/tejasdeepakmasne/mfer/lexer"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	progArgs := os.Args

	if len(progArgs) == 3 {
		if progArgs[1] == "lex" {
			outputTokens := lexer.RunLexerFromFile(progArgs[2])
			lexer.PrintTokens(outputTokens)
		} else {
			fmt.Println("USAGE: mfer lex <filepath>")
		}
	} else if len(progArgs) == 2 {
		if progArgs[1] == "lex" {
			lexer.RunLexerPrompt()
		} else {
			fmt.Println("USAGE: mfer <option> where <option> can be lex")
		}
	} else {
		fmt.Println("USAGE: mfer <option>")
	}

}
