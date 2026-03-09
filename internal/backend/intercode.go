package backend

import "github.com/SamJohn04/simple-lang-compiler/internal/common"

// returns 3-Address Code
func IntermediateCodeGenerator(
	input common.ProgramAST,
	identifiers []common.IdentifierInformation,
) ([]string, []common.IdentifierInformation, error) {
	numberOfGotos := 0
	codes, identifiers, err := input.ThreeAddressCode(identifiers, &numberOfGotos)
	if err != nil {
		return []string{}, identifiers, &common.InternalError{
			PointOfFailure: "Intermediate Code Generator",
			Message:        err.Error(),
		}
	}
	return codes, identifiers, nil
}
