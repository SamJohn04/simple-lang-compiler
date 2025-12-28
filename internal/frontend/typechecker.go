package frontend

import (
	"fmt"

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
	switch len(input.ChildNodes) {
	case 0:
		// epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case 2:
		// I1;I
		childI1, err := checkI1(input.ChildNodes[0])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childI, err := checkI(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		output := common.SyntaxTreeNode{
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

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError(
			fmt.Sprintf("I expects I>I1;I or I, got %v", input.InnerToken.Token),
		)
	}
}

func checkI1(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// v=E
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

	case common.TokenLet:
		// let I6
		// let is unnecessary for further calculations
		childI6, err := checkI6(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		return childI6, nil

	case common.TokenIf:
		// if R { I } I4
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

		childIf.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)
		return childIf, nil

	case common.TokenWhile:
		// while R { I }
		childWhile := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI, err := checkI(input.ChildNodes[2])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childWhile.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		return childWhile, nil

	case common.TokenOutput:
		// output E
		childOutput := input.ChildNodes[0]
		childE, err := checkE(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

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
	switch len(input.ChildNodes) {
	case 0:
		// epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "noop",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case 2:
		// else I7
		if input.ChildNodes[0].InnerToken.TokenKind != common.TokenElse {
			return common.SyntaxTreeNode{}, typeCheckerInternalError(
				fmt.Sprintf(
					"I7 did not have else; had %v",
					common.NameMapWithTokenKind[input.ChildNodes[0].InnerToken.TokenKind],
				),
			)
		}
		return checkI7(input.ChildNodes[1])

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I7 does not have 0 or 2 children")
	}
}

func checkI6(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// v = E
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

		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childE,
		}

		return childEquals, nil

	case common.TokenMutable:
		// mut v I8
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
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIf:
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

		childIf.ChildNodes = []common.SyntaxTreeNode{
			childR,
			childI,
		}
		childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)

		return childIf, nil

	case common.TokenBlock:
		childI, err := checkI(input.ChildNodes[0])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		// Preventing removal of I block when expended
		return common.SyntaxTreeNode{
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
	switch len(input.ChildNodes) {
	case 0:
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "noop",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case 2:
		childEquals := input.ChildNodes[0]
		childE, err := checkE(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		IdentTable[childIdentifier.InnerToken.Token] = IdentifierInformation{
			DataType: common.TypedInt,
			Mutable:  true,
		}

		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childE,
		}
		return childEquals, nil

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I8 does not match")
	}
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

	childOperator.ChildNodes = []common.SyntaxTreeNode{
		firstChildE,
		secondChildE,
	}
	return childOperator, nil
}

func checkE(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	childT, err := checkT(input.ChildNodes[0])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkE1(input.ChildNodes[1], childT)
}

func checkE1(input, calculationsUntilNow common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}

	childOperator := input.ChildNodes[0]
	childT, err := checkT(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

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
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}

	childOperator := input.ChildNodes[0]
	childF, err := checkF(input.ChildNodes[1])
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOperator.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childF,
	}
	return checkT1(input.ChildNodes[2], childOperator)
}

func checkF(input common.SyntaxTreeNode) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenExpressionSub:
		childSub := input.ChildNodes[0]
		childF, err := checkF(input.ChildNodes[1])
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childSub.ChildNodes = []common.SyntaxTreeNode{
			childF,
		}
		return childSub, nil

	case common.TokenBlock:
		return checkE(input.ChildNodes[0])

	case common.TokenIdent:
		information, ok := IdentTable[input.ChildNodes[0].InnerToken.Token]
		if !ok || information.DataType != common.TypedInt { // TODO update typechecking
			return common.SyntaxTreeNode{}, typeCheckerCompilationError("use of identifier without declaring")
		}
		fallthrough
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenInput:
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
