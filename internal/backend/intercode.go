package backend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var (
	numberOfIdentifiers int
	identifierData      map[string]int
	stack               []int
)

// returns 3-Address Code
func IntermediateCodeGenerator(input common.SyntaxTreeNode) ([]string, error) {
	numberOfIdentifiers = 0
	identifierData = make(map[string]int)
	stack = []int{}

	intermediateCodes := []string{}
	for _, child := range input.ChildNodes {
		codesFromChild, err := generateNextInstructionSet(child)
		if err != nil {
			return intermediateCodes, err
		}
		intermediateCodes = append(intermediateCodes, codesFromChild...)
	}
	return intermediateCodes, nil
}

func generateNextInstructionSet(input common.SyntaxTreeNode) ([]string, error) {
	if input.InnerToken.Token == "noop" {
		// a noop operation, which could be due to declaring and not initialising an identifier
		return []string{}, nil
	}
	switch input.InnerToken.TokenKind {
	case common.TokenAssignment:
		// v = E
		codesFromChild, outputVariable, err := generateForExpression(input.ChildNodes[1])
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
		return generateForIfStatement(input)

	case common.TokenWhile:
		// while R { I }

	case common.TokenOutput:
		// output E
		codesFromChild, outputVariable, err := generateForExpression(input.ChildNodes[0])
		if err != nil {
			return []string{}, err
		}
		codesFromChild = append(
			codesFromChild,
			fmt.Sprintf("print %v", outputVariable),
		)
		return codesFromChild, nil

	default:
		return []string{}, intermediateCodeGeneratorInternalError(
			fmt.Sprintf(
				"Unexpected token received: %v",
				common.NameMapWithTokenKind[input.InnerToken.TokenKind],
			),
		)
	}
	return []string{}, &common.UnderConstructionError{
		PointOfFailure: "Intermediate Code Generator (Instruction)",
		Message:        "",
	}
}

func generateForIfStatement(input common.SyntaxTreeNode) ([]string, error) {
	// if the condition is true goto L1
	// goto L2
	// L1: --- if block ---
	// goto LEND
	// L2: if the condition is true goto L3 \\ in the case of else if
	// goto L4
	// ...
	_, _, err := generateForRelation(input.ChildNodes[0])
	if err != nil {
		return []string{}, err
	}
	return []string{}, &common.UnderConstructionError{
		PointOfFailure: "Intermediate Code Generator (Instruction)",
		Message:        "",
	}
}

func generateForRelation(input common.SyntaxTreeNode) ([]string, string, error) {
	switch input.InnerToken.TokenKind {
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
		firstChildCodes, firstOutputVariable, err := generateForExpression(input.ChildNodes[0])
		if err != nil {
			return []string{}, "", err
		}
		secondChildCodes, secondOutputVariable, err := generateForExpression(input.ChildNodes[1])
		if err != nil {
			return []string{}, "", err
		}

		identifier := getNextIdentifier()

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
				"Unexpected token %v at a relation",
				common.NameMapWithTokenKind[input.InnerToken.TokenKind],
			),
		)
	}
}

func generateForExpression(input common.SyntaxTreeNode) ([]string, string, error) {
	switch input.InnerToken.TokenKind {
	case common.TokenLiteralInt:
		return []string{}, input.InnerToken.Token, nil
	case common.TokenIdent:
		return []string{}, input.InnerToken.Token, nil
	case common.TokenInput:
		identifier := getNextIdentifier()
		codes := []string{
			fmt.Sprintf(
				"%v = input",
				identifier,
			),
		}
		return codes, identifier, nil

	case common.TokenExpressionSub:
		if len(input.ChildNodes) == 1 {
			// t = - t
			childCodes, outputVariable, err := generateForExpression(input.ChildNodes[0])
			if err != nil {
				return []string{}, "", err
			}

			nextIdentifier := getNextIdentifier()
			childCodes = append(
				childCodes,
				fmt.Sprintf("%v = - %v", nextIdentifier, outputVariable),
			)
			return childCodes, nextIdentifier, nil
		}
		fallthrough
	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionMul:
		fallthrough
	case common.TokenExpressionDiv:
		fallthrough
	case common.TokenExpressionModulo:
		// t = t op t
		firstChildCodes, firstOutputVariable, err := generateForExpression(input.ChildNodes[0])
		if err != nil {
			return []string{}, "", err
		}
		secondChildCodes, secondOutputVariable, err := generateForExpression(input.ChildNodes[1])
		if err != nil {
			return []string{}, "", err
		}

		nextIdentifier := getNextIdentifier()
		codes := []string{}
		codes = append(codes, firstChildCodes...)
		codes = append(codes, secondChildCodes...)
		codes = append(
			codes,
			fmt.Sprintf(
				"%v = %v %v %v",
				nextIdentifier,
				firstOutputVariable,
				common.Operators[input.InnerToken.TokenKind],
				secondOutputVariable,
			),
		)
		return codes, nextIdentifier, nil

	default:
		return []string{}, "", intermediateCodeGeneratorInternalError(
			fmt.Sprintf(
				"unknown token at expression %v",
				common.NameMapWithTokenKind[input.InnerToken.TokenKind],
			),
		)
	}
}

func getNextIdentifier() string {
	numberOfIdentifiers++
	return fmt.Sprintf("t%d", numberOfIdentifiers)
}

func intermediateCodeGeneratorInternalError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Intermediate Code Generator",
		Message:        message,
	}
}
