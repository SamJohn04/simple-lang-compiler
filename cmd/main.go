package main

import (
	"fmt"
	"os"

	"github.com/SamJohn04/simple-lang-compiler/internal/backend"
	"github.com/SamJohn04/simple-lang-compiler/internal/common"
	"github.com/SamJohn04/simple-lang-compiler/internal/frontend"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Create an unbuffered channel for lexical tokens
	lex := make(chan common.Token)

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

	intermediateCodes, err := backend.IntermediateCodeGenerator(programRoot)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(intermediateCodes)
}
