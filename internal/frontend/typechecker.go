package frontend

import (
	"fmt"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

type IdentifierInformation struct {
	DataType common.DataTypeOfIdentifier
	Mutable  bool
}

var IdentTable map[string]IdentifierInformation

func TypeChecker(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	IdentTable = make(map[string]IdentifierInformation)
	output, err := checkI(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return output, nil
}

func checkI(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if len(inputTokenStringAsBlocks) == 0 {
		return common.SyntaxTreeNode{}, TypeCheckerInternalError("string Token of input is 0")
	}
	if inputTokenStringAsBlocks[0] != "I" {
		return common.SyntaxTreeNode{}, TypeCheckerInternalError("not I when expecting I")
	}
	if input.IsLeaf {
		return common.SyntaxTreeNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil
	}
	if len(input.ChildNodes) != 2 {
		return common.SyntaxTreeNode{}, TypeCheckerInternalError("I not containing 0 or 2 children")
	}

	childI1, err := checkI1(input.ChildNodes[0])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childI, err := checkI(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	output := common.SyntaxTreeNode{
		IsLeaf: false,
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I>I1*",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childI1,
		},
	}
	output.ChildNodes = append(output.ChildNodes, childI.ChildNodes...)
	return output, nil
}

func checkI1(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if len(inputTokenStringAsBlocks) != 2 {
		return common.SyntaxTreeNode{}, TypeCheckerInternalError(
			fmt.Sprintf(
				"blocksize of I1 is %v, not 2",
				len(inputTokenStringAsBlocks),
			),
		)
	}
	if inputTokenStringAsBlocks[0] != "I1" {
		return common.SyntaxTreeNode{}, TypeCheckerInternalError("not I1 when expecting I1")
	}

	switch inputTokenStringAsBlocks[1] {
	case "v=E":
		if len(input.ChildNodes) != 3 {
			return common.SyntaxTreeNode{}, TypeCheckerInternalError("I1>v=E has an incorrect number of children")
		}

		identifierChild := input.ChildNodes[0]
		if identifierChild.InnerToken.TokenKind != common.TokenIdent {
			return common.SyntaxTreeNode{}, TypeCheckerInternalError(
				fmt.Sprintf(
					"I>v=E does not have an identifier as the first child; has %v",
					common.NameMapWithTokenKind[identifierChild.InnerToken.TokenKind],
				),
			)
		}

		information, ok := IdentTable[identifierChild.InnerToken.Token]
		if !ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf(
					"Identifier %v has not been declared",
					identifierChild.InnerToken.Token,
				),
			)
		} else if !information.Mutable {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf(
					"Identifier %v was not marked as mutable (mut keyword).\nE.g.: let mut id = 3;",
					identifierChild.InnerToken.Token,
				),
			)
		}
		return common.SyntaxTreeNode{}, &common.UnderConstructionError{
			PointOfFailure: "Type Checker",
			Message:        "",
		}
	}
	return common.SyntaxTreeNode{}, &common.UnderConstructionError{
		PointOfFailure: "Type Checker",
		Message:        "",
	}
}

func typeCheckerCompilationError(message string) *common.CompilationError {
	return &common.CompilationError{
		PointOfFailure: "Type Checker",
		Message:        message,
	}
}

func TypeCheckerInternalError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Type Checker",
		Message:        message,
	}
}
