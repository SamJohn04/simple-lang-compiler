package main

import (
	"fmt"
	"os"
	"os/exec"
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

	cCode, err := backend.CodeGenerator(intermediateCodes, identifiers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// output file name is the input file with the sl removed and 'out' added.
	// Just in case the file name has no extension, a "." is (potentially) removed and added again
	// Expects gcc in your system
	outputFileName := fmt.Sprintf("%v.out", strings.TrimSuffix(inputFileName, ".sl"))
	if len(os.Args) >= 3 {
		outputFileName = os.Args[2]
	}
	err = toObjectFile(outputFileName, cCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func toObjectFile(outputFileName, cCode string) error {
	tmpFile, err := os.CreateTemp("", "prog-*.c")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(cCode); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write c code to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	cmd := exec.Command("gcc", tmpFile.Name(), "-o", outputFileName)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gcc failed: %w", err)
	}
	return nil
}
