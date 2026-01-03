package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/backend"
	"github.com/SamJohn04/simple-lang-compiler/internal/common"
	"github.com/SamJohn04/simple-lang-compiler/internal/frontend"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("at least 1 argument required")
		os.Exit(1)
	}
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

	intermediateCodes, err := backend.IntermediateCodeGenerator(programRoot)
	if err != nil {
		fmt.Println(strings.Join(intermediateCodes, "\n"))
		fmt.Println(err)
		os.Exit(1)
	}
	finalCode, err := backend.CodeGenerator(intermediateCodes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output_file_name := "output.c"
	if len(os.Args) >= 3 {
		output_file_name = os.Args[2]
	}
	os.WriteFile(output_file_name, []byte(finalCode), 0o664)
}
