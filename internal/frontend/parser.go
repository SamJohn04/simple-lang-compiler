package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

func Parser(input <-chan common.Token) {
	for inToken := range input {
		fmt.Println(inToken.Token, common.NameMapWithTokenKind[inToken.TokenKind])
	}
}
