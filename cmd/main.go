package main

import (
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
	frontend.Parser(lex)
}
