package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

func TypeChecker(
	input common.SyntaxTreeNode,
) (common.SyntaxTreeNode, map[string]common.IdentifierInformation, error) {
	// golang always gives a 0 value if undefined
	// which we have maaped to the unknown type
	// as such, SyntaxTreeNode.Datatype is not necessary till here
	identifierTable := make(map[string]common.IdentifierInformation)
	output, err := checkProgram(input, identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, nil, err
	}
	return output, identifierTable, nil
}

func checkProgram(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
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
		childI1, err := checkNextInstruction(input.ChildNodes[0], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childI, err := checkProgram(input.ChildNodes[1], identifierTable)
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

func checkNextInstruction(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// v=R
		return checkReassignment(input, identifierTable)

	case common.TokenLet:
		// let I6
		return checkAssignment(input, identifierTable)

	case common.TokenIf:
		// if R { I } I4
		return checkIf(input, identifierTable)

	case common.TokenWhile:
		// while R { I }
		return checkWhile(input, identifierTable)

	case common.TokenOutput:
		// output str C
		return checkOutput(input, identifierTable)

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I1 does not match")
	}
}

func checkReassignment(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	// v=E
	childIdentifier := input.ChildNodes[0]
	childEquals := input.ChildNodes[1]
	childR, err := checkR(input.ChildNodes[2], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	information, ok := identifierTable[childIdentifier.InnerToken.Token]
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
	} else if information.DataType == common.TypedUnkown {
		identifierTable[childIdentifier.InnerToken.Token] = common.IdentifierInformation{
			DataType: childR.Datatype,
			Mutable:  true,
		}
	} else if information.DataType != childR.Datatype {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			fmt.Sprintf(
				"Datatype information (%v) and datatype to be assigned (%v) are not the same",
				common.NameMapWithType[information.DataType],
				common.NameMapWithType[childR.Datatype],
			),
		)
	}

	childEquals.ChildNodes = []common.SyntaxTreeNode{
		childIdentifier,
		childR,
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
}

func checkAssignment(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	// let I6
	// the idea is to declare the variable if it does not exist

	// we do not need let for further computations
	childI6, err := checkAssignmentAfterLet(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return childI6, nil
}

func checkAssignmentAfterLet(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// v = R
		childIdentifier := input.ChildNodes[0]
		childEquals := input.ChildNodes[1]
		childR, err := checkR(input.ChildNodes[2], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		_, ok := identifierTable[childIdentifier.InnerToken.Token]
		if ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf("redeclaring existing variable \"%v\"", childIdentifier.InnerToken.Token),
			)
		}
		identifierTable[childIdentifier.InnerToken.Token] = common.IdentifierInformation{
			DataType: childR.Datatype,
			Mutable:  false,
		}

		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childR,
		}

		return childEquals, nil

	case common.TokenMutable:
		// mut v I8
		// mut is not necessary for further calculations, but serves as a flag
		childIdentifier := input.ChildNodes[1]

		_, ok := identifierTable[childIdentifier.InnerToken.Token]
		if ok {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				fmt.Sprintf("redeclaring existing variable \"%v\"", childIdentifier.InnerToken.Token),
			)
		}
		identifierTable[childIdentifier.InnerToken.Token] = common.IdentifierInformation{
			DataType: common.TypedUnkown,
			Mutable:  true,
		}

		return checkMutableAssignment(input.ChildNodes[2], childIdentifier, identifierTable)

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I6 does not match")
	}
}

func checkMutableAssignment(
	input,
	childIdentifier common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	switch len(input.ChildNodes) {
	case 0:
		// declare v here
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenDeclare,
				Token:     "declare",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childIdentifier,
			},
		}, nil

	case 2:
		childEquals := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		identifierTable[childIdentifier.InnerToken.Token] = common.IdentifierInformation{
			DataType: childR.Datatype,
			Mutable:  true,
		}

		childEquals.ChildNodes = []common.SyntaxTreeNode{
			childIdentifier,
			childR,
		}
		return childEquals, nil

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I8 does not match")
	}
}

func checkIf(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	// if R { I } I4
	childIf := input.ChildNodes[0]
	childR, err := checkR(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	if childR.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"if has non-boolean expression as condition",
		)
	}
	childI, err := checkProgram(input.ChildNodes[2], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childI4, err := checkElseCondition(input.ChildNodes[3], identifierTable)
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
			| \       \
			if if      else
			|\ \ \      \
			R I R I ...  I
	*/

	// extra if block for code generation
	grandChildIf := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenIf,
			Token:     "if (extra)",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childR,
			childI,
		},
	}

	childIf.ChildNodes = []common.SyntaxTreeNode{
		grandChildIf,
	}
	childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)
	return childIf, nil
}

func checkElseCondition(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
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
		return checkElseIf(input.ChildNodes[1], identifierTable)

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I7 does not have 0 or 2 children")
	}
}

func checkElseIf(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIf:
		childIf := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		if childR.Datatype != common.TypedBool {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				"else if has non-boolean expression as condition",
			)
		}
		childI, err := checkProgram(input.ChildNodes[2], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childI4, err := checkElseCondition(input.ChildNodes[3], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		// extra if block for code generation
		grandChildIf := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenIf,
				Token:     "if (extra)",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childR,
				childI,
			},
		}

		childIf.ChildNodes = []common.SyntaxTreeNode{
			grandChildIf,
		}
		childIf.ChildNodes = append(childIf.ChildNodes, childI4.ChildNodes...)

		return childIf, nil

	case common.TokenBlock:
		childI, err := checkProgram(input.ChildNodes[0], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		// extra else for code generation
		childElse := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenElse,
				Token:     "else (extra)",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childI,
			},
		}

		// Preventing removal of block when expanded
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "end of if",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childElse,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, typeCheckerInternalError("I7 does not match")
	}
}

func checkWhile(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	// while R { I }
	childWhile := input.ChildNodes[0]
	childR, err := checkR(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	if childR.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"non-boolean expression in while",
		)
	}
	childI, err := checkProgram(input.ChildNodes[2], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childWhile.ChildNodes = []common.SyntaxTreeNode{
		childR,
		childI,
	}
	return childWhile, nil
}

func checkOutput(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	// output str C
	childOutput := input.ChildNodes[0]
	childStr := input.ChildNodes[1]
	childC, err := checkOutputContinuation(input.ChildNodes[2], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOutput.ChildNodes = []common.SyntaxTreeNode{
		childStr,
	}
	childOutput.ChildNodes = append(childOutput.ChildNodes, childC.ChildNodes...)
	return childOutput, nil
}

func checkOutputContinuation(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return input, nil
	}
	childOutput, err := checkOutputContinuation(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childR, err := checkR(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	input.ChildNodes = []common.SyntaxTreeNode{
		childR,
	}
	input.ChildNodes = append(input.ChildNodes, childOutput.ChildNodes...)
	return input, nil
}

func checkR(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	childRa, err := checkRa(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkRz(input.ChildNodes[1], childRa, identifierTable)
}

func checkRz(
	input,
	calculationsUntilNow common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if calculationsUntilNow.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"|| relations must have boolean children",
		)
	}

	childOr := input.ChildNodes[0]
	childRa, err := checkRa(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	if childRa.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"|| relations must have boolean children",
		)
	}

	childOr.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childRa,
	}
	childOr.Datatype = common.TypedBool
	return checkRz(input.ChildNodes[2], childOr, identifierTable)
}

func checkRa(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	childRb, err := checkRb(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkRy(input.ChildNodes[1], childRb, identifierTable)
}

func checkRy(
	input,
	calculationsUntilNow common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if calculationsUntilNow.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"&& relations must have boolean children",
		)
	}

	childAnd := input.ChildNodes[0]
	childRb, err := checkRb(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	if childRb.Datatype != common.TypedBool {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"&& relations must have boolean children",
		)
	}

	childAnd.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childRb,
	}
	childAnd.Datatype = common.TypedBool
	return checkRy(input.ChildNodes[2], childAnd, identifierTable)
}

func checkRb(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if input.ChildNodes[0].InnerToken.TokenKind == common.TokenNot {
		childNot := input.ChildNodes[0]
		childR, err := checkR(input.ChildNodes[1], identifierTable)
		if childR.Datatype != common.TypedBool {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				"! relation must have a boolean child",
			)
		}

		childNot.ChildNodes = []common.SyntaxTreeNode{
			childR,
		}
		childNot.Datatype = common.TypedBool
		return childNot, err
	}
	if len(input.ChildNodes[1].ChildNodes) == 0 {
		// R1 is empty
		return checkE(input.ChildNodes[0], identifierTable)
	}
	/*
		Converts
			R
			|\
			E R1
			  | \
			  op E

		to this
			op
			|\
			E E
	*/
	firstChildE, err := checkE(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childOperator := input.ChildNodes[1].ChildNodes[0]
	secondChildE, err := checkE(input.ChildNodes[1].ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childOperator.ChildNodes = []common.SyntaxTreeNode{
		firstChildE,
		secondChildE,
	}
	childOperator.Datatype = common.TypedBool
	return childOperator, nil
}

func checkE(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	childT, err := checkT(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkE1(input.ChildNodes[1], childT, identifierTable)
}

func checkE1(
	input,
	calculationsUntilNow common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}

	if calculationsUntilNow.Datatype != common.TypedInt &&
		calculationsUntilNow.Datatype != common.TypedChar &&
		calculationsUntilNow.Datatype != common.TypedFloat {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"only int, float, and char types can be used in mathematical expressions",
		)
	}

	childOperator := input.ChildNodes[0]
	childT, err := checkT(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	if childT.Datatype != common.TypedInt &&
		childT.Datatype != common.TypedChar &&
		childT.Datatype != common.TypedFloat {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"only int, float, and char types can be used in mathematical expressions",
		)
	}

	childOperator.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childT,
	}
	if calculationsUntilNow.Datatype == common.TypedFloat ||
		childT.Datatype == common.TypedFloat {
		childOperator.Datatype = common.TypedFloat
	} else {
		childOperator.Datatype = common.TypedInt
	}
	return checkE1(input.ChildNodes[2], childOperator, identifierTable)
}

func checkT(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	childF, err := checkF(input.ChildNodes[0], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return checkT1(input.ChildNodes[1], childF, identifierTable)
}

func checkT1(
	input,
	calculationsUntilNow common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}

	if calculationsUntilNow.Datatype != common.TypedInt &&
		calculationsUntilNow.Datatype != common.TypedChar &&
		calculationsUntilNow.Datatype != common.TypedFloat {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"only int, float, and char types can be used in mathematical expressions",
		)
	}

	childOperator := input.ChildNodes[0]
	childF, err := checkF(input.ChildNodes[1], identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	if childF.Datatype != common.TypedInt &&
		childF.Datatype != common.TypedChar &&
		childF.Datatype != common.TypedFloat {
		return common.SyntaxTreeNode{}, typeCheckerCompilationError(
			"only int, float, and char types can be used in mathematical expressions",
		)
	}

	childOperator.ChildNodes = []common.SyntaxTreeNode{
		calculationsUntilNow,
		childF,
	}
	if calculationsUntilNow.Datatype == common.TypedFloat ||
		childF.Datatype == common.TypedFloat {
		childOperator.Datatype = common.TypedFloat
	} else {
		childOperator.Datatype = common.TypedInt
	}

	childT1, err := checkT1(input.ChildNodes[2], childOperator, identifierTable)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return childT1, nil
}

func checkF(
	input common.SyntaxTreeNode,
	identifierTable map[string]common.IdentifierInformation,
) (common.SyntaxTreeNode, error) {
	if len(input.ChildNodes) == 0 {
		return common.SyntaxTreeNode{}, typeCheckerInternalError(
			fmt.Sprintf("children of F (%v) is 0; expects at least 1", input.InnerToken.Token),
		)
	}
	switch input.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenExpressionSub:
		childSub := input.ChildNodes[0]
		childF, err := checkF(input.ChildNodes[1], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		/*
			- takes in
				int
				char
				float
			and returns
				int (int, char, bool)
				float
		*/

		if childF.Datatype != common.TypedInt &&
			childF.Datatype != common.TypedChar &&
			childF.Datatype != common.TypedFloat {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError(
				"only int, char, and float types work with negation",
			)
		}

		childSub.ChildNodes = []common.SyntaxTreeNode{
			childF,
		}

		if childF.Datatype == common.TypedChar {
			childSub.Datatype = common.TypedInt
		} else {
			childSub.Datatype = childF.Datatype
		}
		return childSub, nil

	case common.TokenBlock:
		childE, err := checkE(input.ChildNodes[0], identifierTable)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		return childE, nil

	case common.TokenIdent:
		information, ok := identifierTable[input.ChildNodes[0].InnerToken.Token]
		if !ok || information.DataType == common.TypedUnkown {
			return common.SyntaxTreeNode{}, typeCheckerCompilationError("use of identifier without declaring")
		}

		childIdent := input.ChildNodes[0]
		childIdent.Datatype = information.DataType
		return childIdent, nil

	case common.TokenLiteralInt:
		childInt := input.ChildNodes[0]
		childInt.Datatype = common.TypedInt
		return childInt, nil

	case common.TokenLiteralBool:
		childBool := input.ChildNodes[0]
		childBool.Datatype = common.TypedBool
		return childBool, nil

	case common.TokenLiteralChar:
		childChar := input.ChildNodes[0]
		childChar.Datatype = common.TypedChar
		return childChar, nil

	case common.TokenLiteralFloat:
		childFloat := input.ChildNodes[0]
		childFloat.Datatype = common.TypedFloat
		return childFloat, nil

	case common.TokenInput:
		childInput := input.ChildNodes[0]
		// change this based on input type
		childInput.Datatype = common.TypedInt
		return childInput, nil

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
