package backend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// returns 3-Address Code
func IntermediateCodeGenerator(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) ([]string, error) {
	numberOfIdentifiers := 0
	numberOfGotos := 0

	return generateForProgram(
		input,
		identifierTable,
		&numberOfIdentifiers,
		&numberOfGotos,
	)
}

func generateForProgram(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
	numberOfIdentifiers,
	numberOfGotos *int,
) ([]string, error) {
	intermediateCodes := []string{}
	for _, child := range input.ChildNodes {
		codesFromChild, err := generateNextInstructionSet(
			child,
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return intermediateCodes, err
		}
		intermediateCodes = append(intermediateCodes, codesFromChild...)
	}
	return intermediateCodes, nil
}

func generateNextInstructionSet(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
	numberOfIdentifiers,
	numberOfGotos *int,
) ([]string, error) {
	if input.InnerToken.Token == "noop" {
		// a noop operation
		// as of right now, no flags like this should pass here
		return []string{}, intermediateCodeGeneratorInternalError(
			"found a noop statement when none should exist",
		)
	}
	switch input.InnerToken.TokenKind {
	case common.TokenDeclare:
		// derived from mut v
		return []string{}, nil

	case common.TokenAssignment:
		// v = E
		codesFromChild, outputVariable, err := generateForExpression(
			input.ChildNodes[1],
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return []string{}, err
		}
		codesFromChild = append(
			codesFromChild,
			fmt.Sprintf("%v = %v", input.ChildNodes[0].InnerToken.Token, outputVariable),
		)
		return codesFromChild, nil

	case common.TokenIf:
		// if R { I } ...
		return generateForIfStatement(
			input,
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)

	case common.TokenWhile:
		// while R { I }
		return generateForWhileStatement(
			input,
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)

	case common.TokenOutput:
		// output str C
		codes := []string{}
		param := []string{input.ChildNodes[0].InnerToken.Token}
		for _, child := range input.ChildNodes[1:] {
			codesFromChild, outputVariable, err := generateForExpression(
				child,
				identifierTable,
				numberOfIdentifiers,
				numberOfGotos,
			)
			if err != nil {
				return []string{}, err
			}
			codes = append(codes, codesFromChild...)
			param = append(param, outputVariable)
		}
		for _, p := range param {
			codes = append(codes, fmt.Sprintf("param %v", p))
		}
		codes = append(codes, "call printf")
		return codes, nil

	default:
		return []string{}, intermediateCodeGeneratorInternalError(
			fmt.Sprintf(
				"Unexpected token received: %v",
				common.NameMapWithTokenKind[input.InnerToken.TokenKind],
			),
		)
	}
}

func generateForIfStatement(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
	numberOfIdentifiers,
	numberOfGotos *int,
) ([]string, error) {
	// if the condition is true goto L1
	// goto L2
	// L1: --- if block ---
	// goto LX
	// L2: if the condition is true goto L3 \\ in the case of else if
	// goto L4
	// ...
	// LX-1: --- else block --- \\ if it exists
	// LX: --- continue ---
	codes := []string{}
	finalGotoLink := getNextGoto(numberOfGotos)
	for _, child := range input.ChildNodes {
		if child.InnerToken.TokenKind == common.TokenElse {
			childCodes, err := generateForProgram(
				child.ChildNodes[0],
				identifierTable,
				numberOfIdentifiers,
				numberOfGotos,
			)
			if err != nil {
				return []string{}, err
			}
			codes = append(codes, childCodes...)
			break
		}
		codesFromRelation, identifier, err := generateForExpression(
			child.ChildNodes[0],
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return []string{}, err
		}
		childCodes, err := generateForProgram(
			child.ChildNodes[1],
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return []string{}, err
		}
		holdGoto := getNextGoto(numberOfGotos)
		nextGoto := getNextGoto(numberOfGotos)
		codes = append(codes, codesFromRelation...)
		codes = append(codes, fmt.Sprintf("if %v goto %v", identifier, holdGoto))
		codes = append(codes, fmt.Sprintf("goto %v", nextGoto))
		codes = append(codes, fmt.Sprintf("%v:", holdGoto))
		codes = append(codes, childCodes...)
		codes = append(codes, fmt.Sprintf("goto %v", finalGotoLink))
		codes = append(codes, fmt.Sprintf("%v:", nextGoto))
	}
	codes = append(codes, fmt.Sprintf("%v:", finalGotoLink))
	return codes, nil
}

func generateForWhileStatement(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
	numberOfIdentifiers,
	numberOfGotos *int,
) ([]string, error) {
	whileGoto := getNextGoto(numberOfGotos)
	holdGoto := getNextGoto(numberOfGotos)
	nextGoto := getNextGoto(numberOfGotos)
	codes := []string{fmt.Sprintf("%v:", whileGoto)}
	codesFromRelation, identifier, err := generateForExpression(
		input.ChildNodes[0],
		identifierTable,
		numberOfIdentifiers,
		numberOfGotos,
	)
	if err != nil {
		return []string{}, err
	}
	childCodes, err := generateForProgram(
		input.ChildNodes[1],
		identifierTable,
		numberOfIdentifiers,
		numberOfGotos,
	)
	if err != nil {
		return childCodes, err
	}
	codes = append(codes, codesFromRelation...)
	codes = append(codes, fmt.Sprintf("if %v goto %v", identifier, holdGoto))
	codes = append(codes, fmt.Sprintf("goto %v", nextGoto))
	codes = append(codes, fmt.Sprintf("%v:", holdGoto))
	codes = append(codes, childCodes...)
	codes = append(codes, fmt.Sprintf("goto %v", whileGoto))
	codes = append(codes, fmt.Sprintf("%v:", nextGoto))

	return codes, nil
}

func generateForExpression(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
	numberOfIdentifiers,
	numberOfGotos *int,
) ([]string, string, error) {
	switch input.InnerToken.TokenKind {
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenIdent:
		return []string{}, input.InnerToken.Token, nil
	case common.TokenInput:
		identifier := getNextIdentifier(numberOfIdentifiers)
		identifierTable[identifier] = common.IdentifierInformation{
			DataType: input.Datatype,
			Mutable:  false,
		}
		codes := []string{
			fmt.Sprintf(
				"%v = input",
				identifier,
			),
		}
		return codes, identifier, nil

	case common.TokenNot:
		fallthrough
	case common.TokenExpressionSub:
		if len(input.ChildNodes) == 1 {
			// t = op t
			childCodes, outputVariable, err := generateForExpression(
				input.ChildNodes[0],
				identifierTable,
				numberOfIdentifiers,
				numberOfGotos,
			)
			if err != nil {
				return []string{}, "", err
			}

			identifier := getNextIdentifier(numberOfIdentifiers)
			identifierTable[identifier] = common.IdentifierInformation{
				DataType: input.Datatype,
				Mutable:  false,
			}
			childCodes = append(
				childCodes,
				fmt.Sprintf(
					"%v = %v %v",
					identifier,
					common.Operators[common.TokenNot],
					outputVariable,
				),
			)
			return childCodes, identifier, nil
		}
		fallthrough
	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionMul:
		fallthrough
	case common.TokenExpressionDiv:
		fallthrough
	case common.TokenExpressionModulo:
		fallthrough
	case common.TokenRelationalEquals:
		fallthrough
	case common.TokenRelationalGreaterThan:
		fallthrough
	case common.TokenRelationalGreaterThanOrEquals:
		fallthrough
	case common.TokenRelationalLesserThan:
		fallthrough
	case common.TokenRelationalLesserThanOrEquals:
		fallthrough
	case common.TokenRelationalNotEquals:
		fallthrough
	case common.TokenAnd:
		fallthrough
	case common.TokenOr:
		// t = t op t
		firstChildCodes, firstOutputVariable, err := generateForExpression(
			input.ChildNodes[0],
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return []string{}, "", err
		}
		secondChildCodes, secondOutputVariable, err := generateForExpression(
			input.ChildNodes[1],
			identifierTable,
			numberOfIdentifiers,
			numberOfGotos,
		)
		if err != nil {
			return []string{}, "", err
		}

		identifier := getNextIdentifier(numberOfIdentifiers)
		identifierTable[identifier] = common.IdentifierInformation{
			DataType: input.Datatype,
			Mutable:  false,
		}
		codes := []string{}
		codes = append(codes, firstChildCodes...)
		codes = append(codes, secondChildCodes...)
		codes = append(
			codes,
			fmt.Sprintf(
				"%v = %v %v %v",
				identifier,
				firstOutputVariable,
				common.Operators[input.InnerToken.TokenKind],
				secondOutputVariable,
			),
		)
		return codes, identifier, nil

	default:
		return []string{}, "", intermediateCodeGeneratorInternalError(
			fmt.Sprintf(
				"unknown token at expression %v",
				common.NameMapWithTokenKind[input.InnerToken.TokenKind],
			),
		)
	}
}

func getNextGoto(numberOfGotos *int) string {
	(*numberOfGotos)++
	return fmt.Sprintf("L%d", *numberOfGotos)
}

func getNextIdentifier(numberOfIdentifiers *int) string {
	(*numberOfIdentifiers)++
	return fmt.Sprintf("t%d", *numberOfIdentifiers)
}

func intermediateCodeGeneratorInternalError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Intermediate Code Generator",
		Message:        message,
	}
}
