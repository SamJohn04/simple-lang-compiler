package main

import (
	"fmt"
	"os"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
	"github.com/SamJohn04/simple-lang-compiler/internal/frontend"
)

func main() {
	filename := os.Args[1]
	file, _ := os.Open(filename)

	lex := make(chan common.Token)
	go frontend.Lexer(file, lex)

	for o := range lex {
		fmt.Println(o.Token, "Token", common.NameMapWithTokenKind[o.TokenKind])
	}
}
