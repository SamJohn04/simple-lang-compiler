package frontend

import (
	"errors"
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

var (
	IdentTable  map[string]bool
	currPointer common.Token
)

// Parsing is done using LL(1) method.
// TODO Separate concerns:
// // Does both the syntax analyzer and semantic analyzer functions
func Parser(input <-chan common.Token) (common.SyntaxTreeNode, error) {
	movePointerToNextToken(input)
	IdentTable = make(map[string]bool)

	output := common.SyntaxTreeNode{
		IsLeaf: false,
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "",
		},
		ChildNodes: []common.SyntaxTreeNode{},
	}
	err := parseI(&output, input)
	if err != nil {
		fmt.Printf("error found at %v: %v", currPointer.Token, err)
		return common.SyntaxTreeNode{}, err
	}
	return output, nil
}

func parseI(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent ||
		currPointer.TokenKind == common.TokenIf ||
		currPointer.TokenKind == common.TokenWhile ||
		currPointer.TokenKind == common.TokenLet ||
		currPointer.TokenKind == common.TokenOutput {
		// I -> I1;I
		err := parseI1(output, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenLineEnd {
			return errors.New("end of line (;) expected.")
		}
		movePointerToNextToken(input)
		return parseI(output, input)
	} else if currPointer.TokenKind == common.TokenEOF ||
		currPointer.TokenKind == common.TokenCloseCurly {
		// I -> epsilon
		return nil
	}
	return errors.New("unexpected parse token in I")
}

func parseI1(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent {
		// I1 -> v=E
		value, ok := IdentTable[currPointer.Token]
		if !ok || !value {
			// variable should be initialized with mut keyword first
			// if its value will change.
			return errors.New("variable not declared or is not mutable")
		}
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.SyntaxTreeNode{
				{
					IsLeaf: true,
					InnerToken: common.Token{
						TokenKind: currPointer.TokenKind,
						Token:     currPointer.Token,
					},
					ChildNodes: []common.SyntaxTreeNode{},
				},
			},
		})
		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' expected")
		}
		movePointerToNextToken(input)
		return parseE(&output.ChildNodes[len(output.ChildNodes)-1], input)
	} else if currPointer.TokenKind == common.TokenIf {
		// I1 -> if R { I } I4
		ifBlock := common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}
		movePointerToNextToken(input)

		// The idea is to have:
		// 				if
		// |-------------------------|
		// |  |         |   |        |
		// R Block [...(R Block)] [Block]
		err := parseR(&ifBlock, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return errors.New("'{' expected")
		}

		ifBlock.ChildNodes = append(ifBlock.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		})
		movePointerToNextToken(input)
		err = parseI(&ifBlock.ChildNodes[len(ifBlock.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return errors.New("'}' expected")
		}
		movePointerToNextToken(input)
		err = parseI4(&ifBlock, input)
		if err != nil {
			return err
		}
		output.ChildNodes = append(output.ChildNodes, ifBlock)
		return nil
	} else if currPointer.TokenKind == common.TokenWhile {
		// I1 -> while R { I }
		whileBlock := common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}
		movePointerToNextToken(input)
		err := parseR(&whileBlock, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return errors.New("'{' expected")
		}
		whileBlock.ChildNodes = append(whileBlock.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		})
		movePointerToNextToken(input)
		err = parseI(&whileBlock.ChildNodes[len(whileBlock.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return errors.New("'}' expected")
		}
		output.ChildNodes = append(output.ChildNodes, whileBlock)
		movePointerToNextToken(input)
		return nil
	} else if currPointer.TokenKind == common.TokenLet {
		// I1 -> let I6
		movePointerToNextToken(input)
		return parseI6(output, input)
	} else if currPointer.TokenKind == common.TokenOutput {
		// I1 -> output E
		outputBlock := common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}
		movePointerToNextToken(input)
		err := parseE(&outputBlock, input)
		if err != nil {
			return err
		}
		output.ChildNodes = append(output.ChildNodes, outputBlock)
		return nil
	}
	return errors.New("unexpected parse token in I1")
}

func parseI4(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenElse {
		// I4 -> else I7
		movePointerToNextToken(input)
		return parseI7(output, input)
	} else if currPointer.TokenKind == common.TokenLineEnd {
		// I4 -> epsilon
		return nil
	}
	return errors.New("unexpected parse token in I4; expecting else or ;")
}

func parseI6(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent {
		// I6 -> v=E
		_, ok := IdentTable[currPointer.Token]
		if ok {
			return errors.New("variable already declared")
		}
		IdentTable[currPointer.Token] = false
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.SyntaxTreeNode{
				{
					IsLeaf: true,
					InnerToken: common.Token{
						TokenKind: currPointer.TokenKind,
						Token:     currPointer.Token,
					},
					ChildNodes: []common.SyntaxTreeNode{},
				},
			},
		})
		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' expected")
		}
		movePointerToNextToken(input)
		return parseE(&output.ChildNodes[len(output.ChildNodes)-1], input)
	} else if currPointer.TokenKind == common.TokenMutable {
		// mut v [I3] // I3 continued here
		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenIdent {
			return errors.New("variable expected")
		}
		_, ok := IdentTable[currPointer.Token]
		if ok {
			return errors.New("variable already declared")
		}
		IdentTable[currPointer.Token] = true
		identifier := common.SyntaxTreeNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		}
		movePointerToNextToken(input)
		if currPointer.TokenKind == common.TokenLineEnd {
			return nil
		} else if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' or ';' expected")
		}
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.SyntaxTreeNode{
				identifier,
			},
		})
		movePointerToNextToken(input)
		return parseE(&output.ChildNodes[len(output.ChildNodes)-1], input)
	}
	return errors.New("unexpected parse token in I6")
}

func parseI7(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIf {
		// I7 -> if R { I } I4
		movePointerToNextToken(input)
		err := parseR(output, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return errors.New("'{' expected")
		}
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.SyntaxTreeNode{},
		})
		movePointerToNextToken(input)
		err = parseI(&output.ChildNodes[len(output.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return errors.New("'}' expected")
		}
		movePointerToNextToken(input)
		return parseI4(output, input)
	} else if currPointer.TokenKind == common.TokenOpenCurly {
		// I7 -> { I }
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf:     false,
			InnerToken: common.Token{},
		})
		movePointerToNextToken(input)
		err := parseI(&output.ChildNodes[len(output.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return errors.New("'}' expected")
		}
		movePointerToNextToken(input)
		return nil
	}
	return errors.New("unexpected parse token in I7")
}

func parseR(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent ||
		currPointer.TokenKind == common.TokenLiteralInt ||
		currPointer.TokenKind == common.TokenOpenParanthesis ||
		currPointer.TokenKind == common.TokenInput {
		// R -> ER1
		err := parseE(output, input)
		if err != nil {
			return err
		}
		return parseR1(output, input)
	}
	return errors.New("unexpected parse token in R")
}

func parseR1(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if len(output.ChildNodes) == 0 {
		// this should NOT happen!
		return errors.New("empty set in R1")
	}
	resultE1 := output.ChildNodes[len(output.ChildNodes)-1].ShallowCopy()
	if currPointer.TokenKind != common.TokenRelationalLesserThan &&
		currPointer.TokenKind != common.TokenRelationalGreaterThan &&
		currPointer.TokenKind != common.TokenRelationalEquals &&
		currPointer.TokenKind != common.TokenRelationalLesserThanOrEquals &&
		currPointer.TokenKind != common.TokenRelationalGreaterThanOrEquals &&
		currPointer.TokenKind != common.TokenRelationalNotEquals {
		return errors.New("unexpected parse token in R1")
	}
	// we go from (-> E) to (-> R -> E)
	output.ChildNodes[len(output.ChildNodes)-1].IsLeaf = false
	output.ChildNodes[len(output.ChildNodes)-1].InnerToken.TokenKind = currPointer.TokenKind
	output.ChildNodes[len(output.ChildNodes)-1].InnerToken.Token = currPointer.Token
	output.ChildNodes[len(output.ChildNodes)-1].ChildNodes = []common.SyntaxTreeNode{
		resultE1,
	}
	movePointerToNextToken(input)
	return parseE(&output.ChildNodes[len(output.ChildNodes)-1], input)
}

func parseE(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	err := parseT(output, input)
	if err != nil {
		return err
	}
	return parseE1(output, input)
}

func parseE1(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenRelationalLesserThan ||
		currPointer.TokenKind == common.TokenRelationalGreaterThan ||
		currPointer.TokenKind == common.TokenRelationalEquals ||
		currPointer.TokenKind == common.TokenRelationalLesserThanOrEquals ||
		currPointer.TokenKind == common.TokenRelationalGreaterThanOrEquals ||
		currPointer.TokenKind == common.TokenRelationalNotEquals ||
		currPointer.TokenKind == common.TokenCloseParanthesis ||
		currPointer.TokenKind == common.TokenOpenCurly ||
		currPointer.TokenKind == common.TokenLineEnd {
		// E1 -> epsilon
		return nil
	}
	if currPointer.TokenKind == common.TokenExpressionAdd ||
		currPointer.TokenKind == common.TokenExpressionSub {
		// E1 -> +TE1 | -TE1
		if len(output.ChildNodes) == 0 {
			// this should NOT happen!
			return errors.New("empty set in E1")
		}
		// from (-> T) to (-> (+/-) -> T)
		resultT := output.ChildNodes[len(output.ChildNodes)-1].ShallowCopy()
		output.ChildNodes[len(output.ChildNodes)-1].IsLeaf = false
		output.ChildNodes[len(output.ChildNodes)-1].InnerToken.Token = currPointer.Token
		output.ChildNodes[len(output.ChildNodes)-1].InnerToken.TokenKind = currPointer.TokenKind
		output.ChildNodes[len(output.ChildNodes)-1].ChildNodes = []common.SyntaxTreeNode{
			resultT,
		}
		movePointerToNextToken(input)
		err := parseT(&output.ChildNodes[len(output.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		return parseE1(output, input)
	}
	return errors.New("unexpected parse token in E1")
}

func parseT(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	err := parseF(output, input)
	if err != nil {
		return err
	}
	return parseT1(output, input)
}

func parseT1(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenRelationalLesserThan ||
		currPointer.TokenKind == common.TokenRelationalGreaterThan ||
		currPointer.TokenKind == common.TokenRelationalEquals ||
		currPointer.TokenKind == common.TokenRelationalLesserThanOrEquals ||
		currPointer.TokenKind == common.TokenRelationalGreaterThanOrEquals ||
		currPointer.TokenKind == common.TokenRelationalNotEquals ||
		currPointer.TokenKind == common.TokenExpressionAdd ||
		currPointer.TokenKind == common.TokenExpressionSub ||
		currPointer.TokenKind == common.TokenCloseParanthesis ||
		currPointer.TokenKind == common.TokenOpenCurly ||
		currPointer.TokenKind == common.TokenLineEnd {
		// T1 -> epsilon
		return nil
	}
	if currPointer.TokenKind == common.TokenExpressionMul ||
		currPointer.TokenKind == common.TokenExpressionDiv ||
		currPointer.TokenKind == common.TokenExpressionModulo {
		// T1 -> *FT1 | /FT1 | %FT1
		if len(output.ChildNodes) == 0 {
			// this should NOT happen!
			return errors.New("empty set in T1")
		}
		// -> F to -> (* / '/' / %) -> F
		resultF := output.ChildNodes[len(output.ChildNodes)-1].ShallowCopy()
		output.ChildNodes[len(output.ChildNodes)-1].IsLeaf = false
		output.ChildNodes[len(output.ChildNodes)-1].InnerToken.Token = currPointer.Token
		output.ChildNodes[len(output.ChildNodes)-1].InnerToken.TokenKind = currPointer.TokenKind
		output.ChildNodes[len(output.ChildNodes)-1].ChildNodes = []common.SyntaxTreeNode{
			resultF,
		}
		movePointerToNextToken(input)
		err := parseF(&output.ChildNodes[len(output.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		return parseT1(output, input)
	}
	return errors.New("unexpected parse token in T1")
}

func parseF(output *common.SyntaxTreeNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent ||
		currPointer.TokenKind == common.TokenLiteralInt ||
		currPointer.TokenKind == common.TokenInput {
		output.ChildNodes = append(output.ChildNodes, common.SyntaxTreeNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.SyntaxTreeNode{},
		})
		movePointerToNextToken(input)
		return nil
	}
	if currPointer.TokenKind != common.TokenOpenParanthesis {
		return errors.New("unexpected parse token in F")
	}
	movePointerToNextToken(input)
	err := parseE(output, input)
	if err != nil {
		return err
	}
	if currPointer.TokenKind != common.TokenCloseParanthesis {
		return errors.New("unexpected parse token in F")
	}
	movePointerToNextToken(input)
	return nil
}

func movePointerToNextToken(input <-chan common.Token) {
	var ok bool
	currPointer, ok = <-input
	if !ok {
		currPointer = common.Token{
			TokenKind: common.TokenError,
			Token:     "",
		}
	}
}
