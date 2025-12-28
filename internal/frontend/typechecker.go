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
		return common.SyntaxTreeNode{}, typeCheckerInternalError("string Token of input is 0")
	}
	if inputTokenStringAsBlocks[0] != "I" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("not I when expecting I")
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
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I not containing 0 or 2 children")
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
		return common.SyntaxTreeNode{}, typeCheckerInternalError(
			fmt.Sprintf(
				"blocksize of I1 is %v, not 2",
				len(inputTokenStringAsBlocks),
			),
		)
	}

	if inputTokenStringAsBlocks[0] != "I1" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("not I1 when expecting I1")
	}

	switch inputTokenStringAsBlocks[1] {
	case "v=E":
		childIdentifier := input.ChildNodes[0]
		childEquals := input.ChildNodes[1]
		childE, err := checkE(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		information, ok := IdentTable[childIdentifier.InnerToken.Token]
		if !ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf(
					"Identifier %v has not been declared",
					childIdentifier.InnerToken.Token,
				),
			)
		} else if !information.Mutable {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf(
					"Identifier %v was not marked as mutable (mut keyword) and was reassigned",
					childIdentifier.InnerToken.Token,
				),
			)
		}

		if information.DataType == common.TypedUnkown {
			IdentTable[childIdentifier.InnerToken.Token] = IdentifierInformation{
				DataType: common.TypedInt, // TODO check E
				Mutable:  true,
			}
		}

		childEquals.IsLeaf = false
		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childE,
		}

		/*
			Essentially going from
				I1
				| \  \
				id =  E
			To
				=
				| \
				id E
		*/
		return childEquals, nil

	case "let I6":
		if len(input.ChildNodes) != 1 {
			return common.SyntaxTreeNode{}, typeCheckerInternalError("I1>let I6 does not have 1 child")
		}

		childI6, err := checkI6(input.ChildNodes[0])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		return childI6, nil

	case "if R {I} I4":
		childIf := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI, err := checkI(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI4, err := checkI4(input.ChildNodes[3])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		/*
			Converts
				I1
				| \ \ \
				if R I I4
					   |    \
					   else I7
							...
			to
				if
				|\ \ \
				R I R I ...
		*/

		childIf.IsLeaf = false
		childIf.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)
		return childIf, nil

	case "while R {I}":
		childWhile := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI, err := checkI(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childWhile.IsLeaf = false
		childWhile.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		return childWhile, nil

	case "output E":
		childOutput := input.ChildNodes[0]
		childE, err := checkE(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childOutput.IsLeaf = false
		childOutput.ChildNodes = []common.SyntaxTreeNode{
			childE,
		}
		return childOutput, nil

	default:
		return input, &common.UnderConstructionError{
			PointOfFailure: "Type Checker",
			Message:        fmt.Sprintf("I1 when \"%v\"", input.InnerToken.Token),
		}
	}
}

func checkI4(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if inputTokenStringAsBlocks[0] != "I4" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("Expecting I4")
	}

	if len(inputTokenStringAsBlocks) == 1 {
		return common.SyntaxTreeNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "noop",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil
	}

	if inputTokenStringAsBlocks[1] != "else I7" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I4 does not match")
	}

	return checkI7(input.ChildNodes[1])
}

func checkI6(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if len(inputTokenStringAsBlocks) != 2 {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I6 does not have 2 blocks")
	}

	if inputTokenStringAsBlocks[0] != "I6" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("Expecting I6")
	}

	switch inputTokenStringAsBlocks[1] {
	case "v=E":
		childIdentifier := input.ChildNodes[0]
		childEquals := input.ChildNodes[1]
		childE, err := checkE(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		_, ok := IdentTable[childIdentifier.InnerToken.Token]
		if ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf("redeclaring existing variable \"%v\"", childIdentifier.InnerToken.Token),
			)
		}
		IdentTable[childIdentifier.InnerToken.Token] = IdentifierInformation{
			DataType: common.TypedInt, // TODO change to check E
			Mutable:  false,
		}

		childEquals.IsLeaf = false
		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childE,
		}

		return childEquals, nil

	case "mut v I8":
		// mut is not necessary for further calculations, but serves as a flag
		childIdentifier := input.ChildNodes[1]

		_, ok := IdentTable[childIdentifier.InnerToken.Token]
		if ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf("redeclaring existing variable \"%v\"", childIdentifier.InnerToken.Token),
			)
		}
		IdentTable[childIdentifier.InnerToken.Token] = IdentifierInformation{
			DataType: common.TypedUnkown,
			Mutable:  true,
		}

		return checkI8(input.ChildNodes[2], childIdentifier)

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I6 does not match")
	}
}

func checkI7(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	switch inputTokenStringAsBlocks[1] {
	case "if R {I} I4":
		childIf := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1])
		childI, err := checkI(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI4, err := checkI4(input.ChildNodes[3])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childIf.IsLeaf = false
		childIf.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)

		return childIf, nil

	case "{I}":
		childI, err := checkI(input.ChildNodes[0])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		// Preventing removal of I block when expended
		return common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "end of if",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childI,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I7 does not match")
	}
}

func checkI8(input, childIdentifier common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if inputTokenStringAsBlocks[0] != "I8" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("Expecting I8")
	}

	if len(inputTokenStringAsBlocks) == 1 {
		return common.SyntaxTreeNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "noop",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil
	}

	if len(inputTokenStringAsBlocks) > 2 || inputTokenStringAsBlocks[1] != "=E" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I8 error")
	}

	childEquals := input.ChildNodes[0]
	childE, err := checkE(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	IdentTable[childIdentifier.InnerToken.Token] = IdentifierInformation{
		DataType: common.TypedInt,
		Mutable:  true,
	}

	childEquals.IsLeaf = false
	childEquals.ChildNodes = []common.SyntaxTreeNode{
		childIdentifier,
		childE,
	}
	return childEquals, nil
}

func checkR(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	/*
		Converts
			R
			|\  \
			E R1 E

		to this
			R1
			|\
			E E
	*/
	firstChildE, err := checkE(input.ChildNodes[0])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childOperator := input.ChildNodes[1]
	secondChildE, err := checkE(input.ChildNodes[2])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOperator.IsLeaf = false
	childOperator.ChildNodes = []common.SyntaxTreeNode{
		firstChildE,
		secondChildE,
	}
	return childOperator, nil
}

func checkE(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	inputTokenStringAsBlocks := strings.Split(input.InnerToken.Token, ">")
	if len(inputTokenStringAsBlocks) != 1 {
		return common.SyntaxTreeNode{}, typeCheckerInternalError(
			fmt.Sprintf(
				"blocksize of E is %v, not 1",
				len(inputTokenStringAsBlocks),
			),
		)
	}
	if inputTokenStringAsBlocks[0] != "E" {
		return common.SyntaxTreeNode{}, typeCheckerInternalError(
			fmt.Sprintf(
				"expecting E, received %v",
				inputTokenStringAsBlocks[0],
			),
		)
	}

	childT, err := checkT(input.ChildNodes[0])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkE1(input.ChildNodes[1], childT)
}

func checkE1(input, calculationsUntilNow common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	if input.InnerToken.Token == "E1" {
		return calculationsUntilNow, nil
	}

	childOperator := input.ChildNodes[0]
	childT, err := checkT(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOperator.IsLeaf = false
	childOperator.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childT,
	}
	return checkE1(input.ChildNodes[2], childOperator)
}

func checkT(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	childF, err := checkF(input.ChildNodes[0])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkT1(input.ChildNodes[1], childF)
}

func checkT1(input, calculationsUntilNow common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	if input.InnerToken.Token == "T1" {
		return calculationsUntilNow, nil
	}

	childOperator := input.ChildNodes[0]
	childF, err := checkF(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOperator.IsLeaf = false
	childOperator.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childF,
	}
	return checkT1(input.ChildNodes[2], childOperator)
}

func checkF(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	switch input.InnerToken.Token {
	case "F>-F":
		childSub := input.ChildNodes[0]
		childF, err := checkF(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childSub.IsLeaf = false
		childSub.ChildNodes = []common.SyntaxTreeNode{
			childF,
		}
		return childSub, nil
	case "F>(E)":
		return checkE(input.ChildNodes[0])
	case "F>id":
		if input.ChildNodes[0].InnerToken.TokenKind == common.TokenIdent {
			information, ok := IdentTable[input.ChildNodes[0].InnerToken.Token]
			if !ok || information.DataType != common.TypedInt { // TODO update typechecking
				return common.SyntaxTreeNode{}, typeCheckerCompilationError("use of identifier without declaring")
			}
		}
		return input.ChildNodes[0], nil
	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("F does not match")
	}
}

func typeCheckerCompilationError(message string) *common.CompilationError {
	return &common.CompilationError{
		PointOfFailure: "Type Checker",
		Message:        message,
	}
}

func typeCheckerInternalError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Type Checker",
		Message:        message,
	}
}
