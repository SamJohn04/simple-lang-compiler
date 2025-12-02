package frontend

import (
	"errors"
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var clrTable [][]int = [][]int{}

type clrParseItem struct {
	state int
	token common.Token
}

// Parsing is done using CLR(1) method.
func Parser(input <-chan common.Token) {
	// Parser needs to shift, reduce, and goto
	for inToken := range input {
		fmt.Println(inToken.Token, common.NameMapWithTokenKind[inToken.TokenKind])
	}
}

func parserShift(stack []clrParseItem, input <-chan common.Token) error {
	if len(stack) == 0 {
		return errors.New("stack empty")
	}

	token, ok := <-input
	if !ok {
		return errors.New("shift attempt post channel closing")
	}

	_ = clrTable[token.TokenKind]

	return nil
}
