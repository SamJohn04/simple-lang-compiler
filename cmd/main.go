package main

import (
	"fmt"
	"os"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
	"github.com/SamJohn04/simple-lang-compiler/internal/frontend"
)

func main() {
	// Create an unbuffered channel for lexical tokens
	lex := make(chan common.Token)

	filename := os.Args[1]
	file, _ := os.Open(filename)
	defer file.Close()

	go frontend.Lexer(file, lex)
	programRoot, err := frontend.Parser(lex)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	programRoot, err = frontend.TypeChecker(programRoot)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	programRoot.Display("")
}
