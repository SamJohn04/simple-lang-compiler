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
	inputFileName := os.Args[1]
	file, err := os.Open(inputFileName)
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
	loweredProgram, identifiers, err := frontend.SemanticAnalyzer(programRoot)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	loweredProgram, err = frontend.TypeChecker(loweredProgram, identifiers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	intermediateCodes, identifiers, err := backend.IntermediateCodeGenerator(
		loweredProgram,
		identifiers,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	finalCode, err := backend.CodeGenerator(intermediateCodes, identifiers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// output file name is the input file with the sj removed and c added
	// just in case the file name has no extension,
	//	a "." is (potentially) removed and added again
	outputFileName := fmt.Sprintf("%v.c", strings.TrimSuffix(inputFileName, ".sj"))
	if len(os.Args) >= 3 {
		outputFileName = os.Args[2]
	}
	os.WriteFile(outputFileName, []byte(finalCode), 0o664)
}
