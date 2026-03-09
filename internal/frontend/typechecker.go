package frontend

import "github.com/SamJohn04/simple-lang-compiler/internal/common"

func TypeChecker(
	input common.ProgramAST, identifiers []common.IdentifierInformation,
) (common.ProgramAST, error) {
	// golang always gives a 0 value if undefined
	// which we have maaped to the unknown type
	// as such, SyntaxTreeNode.Datatype is not necessary till here
	err := input.PerformAllChecks(identifiers)
	if err != nil {
		return common.ProgramAST{}, typeCheckerCompilationError(err.Error())
	}
	return input, nil
}

func typeCheckerCompilationError(message string) *common.CompilationError {
	return &common.CompilationError{
		PointOfFailure: "Type Checker",
		Message:        message,
	}
}
