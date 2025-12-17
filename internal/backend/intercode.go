package backend

import (
	"errors"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// INTEND returns 3-Address Code
func IntermediateCodeGenerator(input common.SyntaxTreeNode, output chan<- common.GeneratorOutput) {
	for _, node := range input.ChildNodes {
		node.Display("")
	}

	output <- common.GeneratorOutput{
		Result: "",
		Err:    errors.New("incomplete"),
	}
}
