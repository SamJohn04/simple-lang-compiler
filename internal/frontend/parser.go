package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var clrTable [][]int = [][]int{}

type clrParseItem struct {
	state int
	token common.Token
}

// Parsing is done using LL(1) method.
func Parser(input <-chan common.Token) {
	// Parser needs to shift, reduce, and goto
	for inToken := range input {
		fmt.Println(inToken.Token, common.NameMapWithTokenKind[inToken.TokenKind])
	}
}
