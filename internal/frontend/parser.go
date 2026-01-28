package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var currPointer common.Token

// Parsing is done using LL(1) method.
func Parser(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	movePointerToNextToken(input)

	output, err := parseProgram(input)
	if err != nil && currPointer.TokenKind == common.TokenError {
		return common.SyntaxTreeNode{}, &common.CompilationError{
			PointOfFailure: "Parser",
			Message:        fmt.Sprintf("Error token causing havoc: \n\t%v", err),
		}
	} else if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	return output, nil
}

func parseProgram(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		fallthrough
	case common.TokenLet:
		fallthrough
	case common.TokenIf:
		fallthrough
	case common.TokenWhile:
		fallthrough
	case common.TokenOutput:
		// I -> I1;I
		childI1, err := parseNextInstruction(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenLineEnd {
			return common.SyntaxTreeNode{}, parserError("end of line (;) expected")
		}

		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I>I1;I",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childI1,
				childI,
			},
		}, err

	case common.TokenEOF:
		fallthrough
	case common.TokenCloseCurly:
		// I -> epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token")
	}
}

func parseNextInstruction(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		// I1 -> v=R
		return parseReassignment(input)

	case common.TokenLet:
		// I1 -> let I6
		return parseAssignment(input)

	case common.TokenIf:
		// I1 -> if R { I } I4
		return parseIf(input)

	case common.TokenWhile:
		// I1 -> while R { I }
		return parseWhile(input)

	case common.TokenOutput:
		// I1 -> output str C
		return parseOutput(input)

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected parse token in I1")
	}
}

func parseReassignment(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	// I1 -> v=R
	childIdent := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	if currPointer.TokenKind != common.TokenAssignment {
		return common.SyntaxTreeNode{}, parserError("'=' expected")
	}
	childEquals := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>v=R",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childIdent,
			childEquals,
			childR,
		},
	}, err
}

func parseAssignment(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	// I1 -> let I6
	childLet := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	childI6, err := parseAssignmentAfterLet(input)
	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>let I6",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childLet,
			childI6,
		},
	}, err
}

func parseAssignmentAfterLet(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		// I6 -> v=R
		childIdent := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return common.SyntaxTreeNode{}, parserError("'=' expected")
		}
		childEquals := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I6>v=R",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childIdent,
				childEquals,
				childR,
			},
		}, err

	case common.TokenMutable:
		// I6 -> mut v I8
		childMut := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenIdent {
			return common.SyntaxTreeNode{}, parserError("variable expected")
		}

		childIdent := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childI8, err := parseMutableAssignment(input)

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I6>mut v I8",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childMut,
				childIdent,
				childI8,
			},
		}, err

	default:
		return common.SyntaxTreeNode{}, parserError("variable or 'mut' expected")
	}
}

func parseMutableAssignment(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenLineEnd:
		// I8 -> epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I8",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case common.TokenAssignment:
		// I8 -> = R
		childEquals := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I8>=R",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childEquals,
				childR,
			},
		}, err

	default:
		return common.SyntaxTreeNode{}, parserError("'=' or ';' expected")
	}
}

func parseIf(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	// I1 -> if R { I } I4
	childIf := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenIf,
			Token:     "if",
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenOpenCurly {
		return common.SyntaxTreeNode{}, parserError("'{' expected")
	}

	movePointerToNextToken(input)
	childI, err := parseProgram(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenCloseCurly {
		return common.SyntaxTreeNode{}, parserError("'}' expected")
	}

	movePointerToNextToken(input)
	childI4, err := parseElseCondition(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>if R {I} I4",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childIf,
			childR,
			childI,
			childI4,
		},
	}, nil
}

func parseElseCondition(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenElse:
		// I4 -> else I7
		childElse := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenElse,
				Token:     "else",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childI7, err := parseElseIf(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I4>else I7",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childElse,
				childI7,
			},
		}, nil

	case common.TokenLineEnd:
		// I4 -> epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I4",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("'else' or ';' expected")
	}
}

func parseElseIf(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIf:
		// I7 -> if R { I } I4
		childIf := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return common.SyntaxTreeNode{}, parserError("'{' expected")
		}

		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return common.SyntaxTreeNode{}, parserError("'}' expected")
		}

		movePointerToNextToken(input)
		childI4, err := parseElseCondition(input)

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I7>if R {I} I4",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childIf,
				childR,
				childI,
				childI4,
			},
		}, err

	case common.TokenOpenCurly:
		// I7 -> { I }
		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenCloseCurly {
			return common.SyntaxTreeNode{}, parserError("'}' expected")
		}
		movePointerToNextToken(input)

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I7>{I}",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childI,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("'if' or '{' expected")
	}
}

func parseWhile(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	// I1 -> while R { I }
	childWhile := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenWhile,
			Token:     "while",
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	if currPointer.TokenKind != common.TokenOpenCurly {
		return common.SyntaxTreeNode{}, parserError("'{' expected")
	}

	movePointerToNextToken(input)
	childI, err := parseProgram(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenCloseCurly {
		return common.SyntaxTreeNode{}, parserError("'}' expected")
	}

	movePointerToNextToken(input)

	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>while R {I}",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childWhile,
			childR,
			childI,
		},
	}, nil
}

func parseOutput(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	// I1 -> output str C
	childOutput := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenOutput,
			Token:     "output",
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	if currPointer.TokenKind != common.TokenLiteralString {
		return common.SyntaxTreeNode{}, parserError("string literal expected after output")
	}
	childStr := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenLiteralString,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}

	movePointerToNextToken(input)
	childC, err := parseOutputContinuation(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	outputBlock := common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>output str C",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childOutput,
			childStr,
			childC,
		},
	}
	return outputBlock, nil
}

func parseOutputContinuation(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenLineEnd:
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "C",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case common.TokenComma:
		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childC, err := parseOutputContinuation(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "C>,E C",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childR,
				childC,
			},
		}, err

	default:
		return common.SyntaxTreeNode{}, parserError("',' or ';' expected")
	}
}

func parseR(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		fallthrough
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenOpenParanthesis:
		fallthrough
	case common.TokenInput:
		fallthrough
	case common.TokenNot:
		fallthrough
	case common.TokenExpressionSub:
		childRa, err := parseRa(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childRz, err := parseRz(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R>Ra Rz",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childRa,
				childRz,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRz(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOr:
		childOr := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenOr,
				Token:     "||",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childRa, err := parseRa(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childRz, err := parseRz(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rz>|| Ra Rz",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childOr,
				childRa,
				childRz,
			},
		}, nil

	case common.TokenOpenCurly:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rz",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRa(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		fallthrough
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenOpenParanthesis:
		fallthrough
	case common.TokenInput:
		fallthrough
	case common.TokenNot:
		fallthrough
	case common.TokenExpressionSub:
		childRb, err := parseRb(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childRy, err := parseRy(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ra>Rb Ry",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childRb,
				childRy,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRy(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenAnd:
		childAnd := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenAnd,
				Token:     "&&",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childRb, err := parseRb(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childRy, err := parseRy(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ry>&& Rb Ry",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childAnd,
				childRb,
				childRy,
			},
		}, nil

	case common.TokenOpenCurly:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		fallthrough
	case common.TokenOr:
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ry",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRb(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenNot:
		childNot := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenNot,
				Token:     "!",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenOpenParanthesis {
			return common.SyntaxTreeNode{}, parserError("There should be a paranthessis set after !")
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenCloseParanthesis {
			return common.SyntaxTreeNode{}, parserError("The paranthessis is not closed")
		}

		movePointerToNextToken(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rb>!(R)",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childNot,
				childR,
			},
		}, nil

	case common.TokenIdent:
		fallthrough
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenInput:
		fallthrough
	case common.TokenExpressionSub:
		fallthrough
	case common.TokenOpenParanthesis:
		childE, err := parseE(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		childR1, err := parseR1(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rb>ER1",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childE,
				childR1,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseR1(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenRelationalLesserThan:
		fallthrough
	case common.TokenRelationalGreaterThan:
		fallthrough
	case common.TokenRelationalEquals:
		fallthrough
	case common.TokenRelationalLesserThanOrEquals:
		fallthrough
	case common.TokenRelationalGreaterThanOrEquals:
		fallthrough
	case common.TokenRelationalNotEquals:
		childOperator := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childE, err := parseE(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R1>opE",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childOperator,
				childE,
			},
		}, err

	case common.TokenOr:
		fallthrough
	case common.TokenAnd:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenOpenCurly:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R1",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseE(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	childT, err := parseT(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}
	childE1, err := parseE1(input)
	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "E",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childT,
			childE1,
		},
	}, err
}

func parseE1(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenAnd:
		fallthrough
	case common.TokenOr:
		fallthrough
	case common.TokenRelationalLesserThan:
		fallthrough
	case common.TokenRelationalGreaterThan:
		fallthrough
	case common.TokenRelationalEquals:
		fallthrough
	case common.TokenRelationalLesserThanOrEquals:
		fallthrough
	case common.TokenRelationalGreaterThanOrEquals:
		fallthrough
	case common.TokenRelationalNotEquals:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenOpenCurly:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		// E1 -> epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "E1",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionSub:
		// E1 -> +TE1 | -TE1
		childArithmeticOperator := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childT, err := parseT(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childE1, err := parseE1(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "E1>opTE1",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childArithmeticOperator,
				childT,
				childE1,
			},
		}, err

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in expression")
	}
}

func parseT(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	childF, err := parseF(input)
	if err != nil {
		return common.SyntaxTreeNode{}, err
	}

	childT1, err := parseT1(input)
	return common.SyntaxTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "T",
		},
		ChildNodes: []common.SyntaxTreeNode{
			childF,
			childT1,
		},
	}, err
}

func parseT1(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenAnd:
		fallthrough
	case common.TokenOr:
		fallthrough
	case common.TokenRelationalLesserThan:
		fallthrough
	case common.TokenRelationalGreaterThan:
		fallthrough
	case common.TokenRelationalEquals:
		fallthrough
	case common.TokenRelationalLesserThanOrEquals:
		fallthrough
	case common.TokenRelationalGreaterThanOrEquals:
		fallthrough
	case common.TokenRelationalNotEquals:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenOpenCurly:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		fallthrough
	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionSub:
		// T1 -> epsilon
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "T1",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}, nil

	case common.TokenExpressionMul:
		fallthrough
	case common.TokenExpressionDiv:
		fallthrough
	case common.TokenExpressionModulo:
		// T1 -> *FT1 | /FT1 | %FT1
		childArithmeticOperator := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childF, err := parseF(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}

		childT1, err := parseT1(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "T1>opFT1",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childArithmeticOperator,
				childF,
				childT1,
			},
		}, nil

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in expression")
	}
}

func parseF(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		fallthrough
	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenInput:
		child := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>id",
			},
			ChildNodes: []common.SyntaxTreeNode{
				child,
			},
		}, nil

	case common.TokenOpenParanthesis:
		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.SyntaxTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenCloseParanthesis {
			return common.SyntaxTreeNode{}, parserError("')' expected")
		}

		movePointerToNextToken(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>(R)",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childR,
			},
		}, nil

	case common.TokenExpressionSub:
		childSub := common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}

		movePointerToNextToken(input)
		childF, err := parseF(input)
		return common.SyntaxTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>-F",
			},
			ChildNodes: []common.SyntaxTreeNode{
				childSub,
				childF,
			},
		}, err

	default:
		return common.SyntaxTreeNode{}, parserError("unexpected token in expression")
	}
}

func movePointerToNextToken(input <-chan common.Token) {
	// defined earlier so that := will not create currPointer as well
	var ok bool
	currPointer, ok = <-input
	if !ok {
		currPointer = common.Token{
			TokenKind: common.TokenError,
			Token:     "",
		}
	}
}

func parserError(message string) *common.CompilationError {
	return &common.CompilationError{
		PointOfFailure: "Parser",
		Message: fmt.Sprintf(
			"%v at %v (line number %v)",
			message,
			currPointer.Token,
			currPointer.LineNumber,
		),
	}
}
