package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var currPointer common.Token

// Parsing is done using LL(1) method.
func Parser(input <-chan common.Token) (common.ParseTreeNode, error) {
	movePointerToNextToken(input)

	output, err := parseProgram(input)
	if err != nil && currPointer.TokenKind == common.TokenError {
		return common.ParseTreeNode{}, &common.CompilationError{
			PointOfFailure: "Parser",
			Message:        fmt.Sprintf("Error token causing havoc: \n\t%v", err),
		}
	} else if err != nil {
		return common.ParseTreeNode{}, err
	}
	return output, nil
}

func parseProgram(input <-chan common.Token) (common.ParseTreeNode, error) {
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
			return common.ParseTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenLineEnd {
			return common.ParseTreeNode{}, parserError("end of line (;) expected")
		}

		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I>I1;I",
			},
			ChildNodes: []common.ParseTreeNode{
				childI1,
				childI,
			},
		}, err

	case common.TokenEOF:
		fallthrough
	case common.TokenCloseCurly:
		// I -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token")
	}
}

func parseNextInstruction(input <-chan common.Token) (common.ParseTreeNode, error) {
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
		return common.ParseTreeNode{}, parserError("unexpected parse token in I1")
	}
}

func parseReassignment(input <-chan common.Token) (common.ParseTreeNode, error) {
	// I1 -> vA=R
	childIdent := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childArrayUsage, err := parseArrayUsage(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenAssignment {
		return common.ParseTreeNode{}, parserError("'=' expected")
	}
	childEquals := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>v=R",
		},
		ChildNodes: []common.ParseTreeNode{
			childIdent,
			childArrayUsage,
			childEquals,
			childR,
		},
	}, err
}

func parseArrayUsage(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOpenSquareBraces:
		// A -> [E]A
		childOpenSquareBraces := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childE, err := parseE(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenCloseSquareBraces {
			return common.ParseTreeNode{}, parserError("square bracket not closed")
		}
		childCloseSquareBraces := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childArrayUsage, err := parseArrayUsage(input)

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "A",
			},
			ChildNodes: []common.ParseTreeNode{
				childOpenSquareBraces,
				childE,
				childCloseSquareBraces,
				childArrayUsage,
			},
		}, nil

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
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		fallthrough
	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionSub:
		fallthrough
	case common.TokenExpressionMul:
		fallthrough
	case common.TokenExpressionDiv:
		fallthrough
	case common.TokenExpressionModulo:
		fallthrough
	case common.TokenAssignment:
		// A -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "A",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token after variable")
	}
}

func parseAssignment(input <-chan common.Token) (common.ParseTreeNode, error) {
	// I1 -> let I6
	childLet := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: currPointer.TokenKind,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childI6, err := parseAssignmentAfterLet(input)
	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>let I6",
		},
		ChildNodes: []common.ParseTreeNode{
			childLet,
			childI6,
		},
	}, err
}

func parseAssignmentAfterLet(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIdent:
		// I6 -> v=R
		childIdent := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return common.ParseTreeNode{}, parserError("'=' expected")
		}
		childEquals := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I6>v=R",
			},
			ChildNodes: []common.ParseTreeNode{
				childIdent,
				childEquals,
				childR,
			},
		}, err

	case common.TokenMutable:
		// I6 -> mut v I8
		childMut := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenIdent {
			return common.ParseTreeNode{}, parserError("variable expected")
		}

		childIdent := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childI8, err := parseMutableAssignment(input)

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I6>mut v I8",
			},
			ChildNodes: []common.ParseTreeNode{
				childMut,
				childIdent,
				childI8,
			},
		}, err

	default:
		return common.ParseTreeNode{}, parserError("variable or 'mut' expected")
	}
}

func parseMutableAssignment(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenLineEnd:
		// I8 -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I8",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	case common.TokenAssignment:
		// I8 -> = R
		childEquals := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I8>=R",
			},
			ChildNodes: []common.ParseTreeNode{
				childEquals,
				childR,
			},
		}, err

	default:
		return common.ParseTreeNode{}, parserError("'=' or ';' expected")
	}
}

func parseIf(input <-chan common.Token) (common.ParseTreeNode, error) {
	// I1 -> if R { I } I4
	childIf := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenIf,
			Token:     "if",
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenOpenCurly {
		return common.ParseTreeNode{}, parserError("'{' expected")
	}
	childOpenCurly := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind:  currPointer.TokenKind,
			Token:      currPointer.Token,
			LineNumber: currPointer.LineNumber,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childI, err := parseProgram(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenCloseCurly {
		return common.ParseTreeNode{}, parserError("'}' expected")
	}
	childCloseCurly := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind:  currPointer.TokenKind,
			Token:      currPointer.Token,
			LineNumber: currPointer.LineNumber,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childI4, err := parseElseCondition(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>if R {I} I4",
		},
		ChildNodes: []common.ParseTreeNode{
			childIf,
			childR,
			childOpenCurly,
			childI,
			childCloseCurly,
			childI4,
		},
	}, nil
}

func parseElseCondition(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenElse:
		// I4 -> else I7
		childElse := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenElse,
				Token:     "else",
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childI7, err := parseElseIf(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I4>else I7",
			},
			ChildNodes: []common.ParseTreeNode{
				childElse,
				childI7,
			},
		}, nil

	case common.TokenLineEnd:
		// I4 -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I4",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("'else' or ';' expected")
	}
}

func parseElseIf(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenIf:
		// I7 -> if R { I } I4
		childIf := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return common.ParseTreeNode{}, parserError("'{' expected")
		}
		childOpenCurly := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return common.ParseTreeNode{}, parserError("'}' expected")
		}
		childCloseCurly := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childI4, err := parseElseCondition(input)

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I7>if R {I} I4",
			},
			ChildNodes: []common.ParseTreeNode{
				childIf,
				childR,
				childOpenCurly,
				childI,
				childCloseCurly,
				childI4,
			},
		}, err

	case common.TokenOpenCurly:
		// I7 -> { I }
		childOpenCurly := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childI, err := parseProgram(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenCloseCurly {
			return common.ParseTreeNode{}, parserError("'}' expected")
		}
		childCloseCurly := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "I7>{I}",
			},
			ChildNodes: []common.ParseTreeNode{
				childOpenCurly,
				childI,
				childCloseCurly,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("'if' or '{' expected")
	}
}

func parseWhile(input <-chan common.Token) (common.ParseTreeNode, error) {
	// I1 -> while R { I }
	childWhile := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenWhile,
			Token:     "while",
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childR, err := parseR(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}
	if currPointer.TokenKind != common.TokenOpenCurly {
		return common.ParseTreeNode{}, parserError("'{' expected")
	}
	childOpenCurly := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind:  currPointer.TokenKind,
			Token:      currPointer.Token,
			LineNumber: currPointer.LineNumber,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childI, err := parseProgram(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	if currPointer.TokenKind != common.TokenCloseCurly {
		return common.ParseTreeNode{}, parserError("'}' expected")
	}
	childCloseCurly := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind:  currPointer.TokenKind,
			Token:      currPointer.Token,
			LineNumber: currPointer.LineNumber,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)

	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>while R {I}",
		},
		ChildNodes: []common.ParseTreeNode{
			childWhile,
			childR,
			childOpenCurly,
			childI,
			childCloseCurly,
		},
	}, nil
}

func parseOutput(input <-chan common.Token) (common.ParseTreeNode, error) {
	// I1 -> output str C
	childOutput := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenOutput,
			Token:     "output",
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	if currPointer.TokenKind != common.TokenLiteralString {
		return common.ParseTreeNode{}, parserError("string literal expected after output")
	}
	childStr := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenLiteralString,
			Token:     currPointer.Token,
		},
		ChildNodes: []common.ParseTreeNode{},
	}

	movePointerToNextToken(input)
	childC, err := parseOutputContinuation(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	outputBlock := common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "I1>output str C",
		},
		ChildNodes: []common.ParseTreeNode{
			childOutput,
			childStr,
			childC,
		},
	}
	return outputBlock, nil
}

func parseOutputContinuation(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenLineEnd:
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "C",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	case common.TokenComma:
		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		childC, err := parseOutputContinuation(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "C>,R C",
			},
			ChildNodes: []common.ParseTreeNode{
				childR,
				childC,
			},
		}, err

	default:
		return common.ParseTreeNode{}, parserError("',' or ';' expected")
	}
}

func parseR(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOpenSquareBraces:
		fallthrough
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
			return common.ParseTreeNode{}, err
		}

		childRz, err := parseRz(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R>Ra Rz",
			},
			ChildNodes: []common.ParseTreeNode{
				childRa,
				childRz,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRz(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOr:
		childOr := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenOr,
				Token:     "||",
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childRa, err := parseRa(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		childRz, err := parseRz(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rz>|| Ra Rz",
			},
			ChildNodes: []common.ParseTreeNode{
				childOr,
				childRa,
				childRz,
			},
		}, nil

	case common.TokenOpenCurly:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rz",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRa(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOpenSquareBraces:
		fallthrough
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
			return common.ParseTreeNode{}, err
		}

		childRy, err := parseRy(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ra>Rb Ry",
			},
			ChildNodes: []common.ParseTreeNode{
				childRb,
				childRy,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRy(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenAnd:
		childAnd := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenAnd,
				Token:     "&&",
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childRb, err := parseRb(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		childRy, err := parseRy(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ry>&& Rb Ry",
			},
			ChildNodes: []common.ParseTreeNode{
				childAnd,
				childRb,
				childRy,
			},
		}, nil

	case common.TokenOpenCurly:
		fallthrough
	case common.TokenCloseParanthesis:
		fallthrough
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		fallthrough
	case common.TokenOr:
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Ry",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseRb(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenNot:
		childNot := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenNot,
				Token:     "!",
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenOpenParanthesis {
			return common.ParseTreeNode{}, parserError("There should be a paranthesis set after !")
		}

		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenCloseParanthesis {
			return common.ParseTreeNode{}, parserError("The paranthesis is not closed")
		}

		movePointerToNextToken(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rb>!(R)",
			},
			ChildNodes: []common.ParseTreeNode{
				childNot,
				childR,
			},
		}, nil

	case common.TokenOpenSquareBraces:
		fallthrough
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
			return common.ParseTreeNode{}, err
		}
		childR1, err := parseR1(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "Rb>ER1",
			},
			ChildNodes: []common.ParseTreeNode{
				childE,
				childR1,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseR1(input <-chan common.Token) (common.ParseTreeNode, error) {
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
		childOperator := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childE, err := parseE(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R1>opE",
			},
			ChildNodes: []common.ParseTreeNode{
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
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "R1",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in relation")
	}
}

func parseE(input <-chan common.Token) (common.ParseTreeNode, error) {
	childT, err := parseT(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}
	childE1, err := parseE1(input)
	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "E",
		},
		ChildNodes: []common.ParseTreeNode{
			childT,
			childE1,
		},
	}, err
}

func parseE1(input <-chan common.Token) (common.ParseTreeNode, error) {
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
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		// E1 -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "E1",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionSub:
		// E1 -> +TE1 | -TE1
		childArithmeticOperator := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childT, err := parseT(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		childE1, err := parseE1(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "E1>opTE1",
			},
			ChildNodes: []common.ParseTreeNode{
				childArithmeticOperator,
				childT,
				childE1,
			},
		}, err

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in expression")
	}
}

func parseT(input <-chan common.Token) (common.ParseTreeNode, error) {
	childF, err := parseF(input)
	if err != nil {
		return common.ParseTreeNode{}, err
	}

	childT1, err := parseT1(input)
	return common.ParseTreeNode{
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "T",
		},
		ChildNodes: []common.ParseTreeNode{
			childF,
			childT1,
		},
	}, err
}

func parseT1(input <-chan common.Token) (common.ParseTreeNode, error) {
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
	case common.TokenCloseSquareBraces:
		fallthrough
	case common.TokenComma:
		fallthrough
	case common.TokenLineEnd:
		fallthrough
	case common.TokenExpressionAdd:
		fallthrough
	case common.TokenExpressionSub:
		// T1 -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "T1",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	case common.TokenExpressionMul:
		fallthrough
	case common.TokenExpressionDiv:
		fallthrough
	case common.TokenExpressionModulo:
		// T1 -> *FT1 | /FT1 | %FT1
		childArithmeticOperator := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childF, err := parseF(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		childT1, err := parseT1(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "T1>opFT1",
			},
			ChildNodes: []common.ParseTreeNode{
				childArithmeticOperator,
				childF,
				childT1,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in expression")
	}
}

func parseF(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenOpenSquareBraces:
		childOpenSquareBraces := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childArrayExpression, err := parseArrayExpression(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}

		if currPointer.TokenKind != common.TokenCloseSquareBraces {
			return common.ParseTreeNode{}, parserError("']' expected")
		}
		childCloseSquareBraces := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}
		movePointerToNextToken(input)

		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  common.TokenBlock,
				Token:      "F>[L]",
				LineNumber: childOpenSquareBraces.InnerToken.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{
				childOpenSquareBraces,
				childArrayExpression,
				childCloseSquareBraces,
			},
		}, nil

	case common.TokenIdent:
		childIdentifier := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  currPointer.TokenKind,
				Token:      currPointer.Token,
				LineNumber: currPointer.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{},
		}
		movePointerToNextToken(input)
		childArrayUsage, err := parseArrayUsage(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind:  common.TokenBlock,
				Token:      "F",
				LineNumber: childIdentifier.InnerToken.LineNumber,
			},
			ChildNodes: []common.ParseTreeNode{
				childIdentifier,
				childArrayUsage,
			},
		}, nil

	case common.TokenLiteralInt:
		fallthrough
	case common.TokenLiteralBool:
		fallthrough
	case common.TokenLiteralChar:
		fallthrough
	case common.TokenLiteralFloat:
		fallthrough
	case common.TokenInput:
		child := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>id",
			},
			ChildNodes: []common.ParseTreeNode{
				child,
			},
		}, nil

	case common.TokenOpenParanthesis:
		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		if currPointer.TokenKind != common.TokenCloseParanthesis {
			return common.ParseTreeNode{}, parserError("')' expected")
		}

		movePointerToNextToken(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>(R)",
			},
			ChildNodes: []common.ParseTreeNode{
				childR,
			},
		}, nil

	case common.TokenExpressionSub:
		childSub := common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ParseTreeNode{},
		}

		movePointerToNextToken(input)
		childF, err := parseF(input)
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "F>-F",
			},
			ChildNodes: []common.ParseTreeNode{
				childSub,
				childF,
			},
		}, err

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in expression")
	}
}

func parseArrayExpression(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenCloseSquareBraces:
		// L -> epsilon
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "L",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	case common.TokenOpenSquareBraces:
		fallthrough
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
		// L -> R L1
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		childArrayContinuation, err := parseArrayContinuation(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "L",
			},
			ChildNodes: []common.ParseTreeNode{
				childR,
				childArrayContinuation,
			},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in array")
	}
}

func parseArrayContinuation(input <-chan common.Token) (common.ParseTreeNode, error) {
	switch currPointer.TokenKind {
	case common.TokenComma:
		movePointerToNextToken(input)
		childR, err := parseR(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		childArrayContinuation, err := parseArrayContinuation(input)
		if err != nil {
			return common.ParseTreeNode{}, err
		}
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "L1",
			},
			ChildNodes: []common.ParseTreeNode{
				childR,
				childArrayContinuation,
			},
		}, nil

	case common.TokenCloseSquareBraces:
		return common.ParseTreeNode{
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "L1",
			},
			ChildNodes: []common.ParseTreeNode{},
		}, nil

	default:
		return common.ParseTreeNode{}, parserError("unexpected token in array")
	}
}

func movePointerToNextToken(input <-chan common.Token) {
	// ok is defined earlier so that := will not create currPointer as well
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
