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

	out := make(chan common.Token)
	go frontend.Lexer(file, out)

	for o := range out {
		fmt.Println(o.Token, "Token", common.NameMapWithTokenKind[o.TokenKind])
	}
}
