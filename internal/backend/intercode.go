package backend

import "github.com/SamJohn04/simple-lang-compiler/internal/common"

// INTEND returns 3-Address Code
func IntermediateCodeGenerator(input common.SyntaxTreeNode) ([]string, error) {
	return []string{}, &common.UnderConstructionError{
		PointOfFailure: "Intermediate Code Generator",
		Message:        "",
	}
}
