package frontend

import (
	"errors"
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

/*
type identRow struct {
	mutable bool
}
*/

var (
	IdentTable  map[string]bool
	currPointer common.Token
)

// Parsing is done using LL(1) method.
func Parser(input <-chan common.Token) error {
	movePointerToNextToken(input)
	IdentTable = make(map[string]bool)

	output := common.ASTNode{
		IsLeaf: false,
		InnerToken: common.Token{
			TokenKind: common.TokenBlock,
			Token:     "",
		},
		ChildNodes: []common.ASTNode{},
	}
	err := parseI(output, input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(output)
	return nil
}

func parseI(output common.ASTNode, input <-chan common.Token) error {
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

func parseI1(output common.ASTNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent {
		// I1 -> v=E
		value, ok := IdentTable[currPointer.Token]
		if !ok || !value {
			// variable should be initialized with mut keyword first
			// if its value will change.
			return errors.New("variable not declared or is not mutable")
		}
		output.ChildNodes = append(output.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.ASTNode{
				{
					IsLeaf: true,
					InnerToken: common.Token{
						TokenKind: currPointer.TokenKind,
						Token:     currPointer.Token,
					},
					ChildNodes: []common.ASTNode{},
				},
			},
		})
		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' expected")
		}
		movePointerToNextToken(input)
		return parseE(output.ChildNodes[len(output.ChildNodes)-1], input)
	} else if currPointer.TokenKind == common.TokenIf {
		// I1 -> if R { I } I4
		ifBlock := common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ASTNode{},
		}
		movePointerToNextToken(input)

		// The idea is to have:
		// 				if
		// |-------------------------|
		// |  |         |   |        |
		// R Block [...(R Block)] [Block]
		err := parseR(ifBlock, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return errors.New("'{' expected")
		}

		ifBlock.ChildNodes = append(ifBlock.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.ASTNode{},
		})
		movePointerToNextToken(input)
		err = parseI(ifBlock.ChildNodes[len(ifBlock.ChildNodes)-1], input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenCloseCurly {
			return errors.New("'}' expected")
		}
		movePointerToNextToken(input)
		err = parseI4(ifBlock, input)
		if err != nil {
			return err
		}
		output.ChildNodes = append(output.ChildNodes, ifBlock)
		return nil
	} else if currPointer.TokenKind == common.TokenWhile {
		// I1 -> while R { I }
		whileBlock := common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ASTNode{},
		}
		movePointerToNextToken(input)
		err := parseR(whileBlock, input)
		if err != nil {
			return err
		}
		if currPointer.TokenKind != common.TokenOpenCurly {
			return errors.New("'{' expected")
		}
		whileBlock.ChildNodes = append(whileBlock.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.ASTNode{},
		})
		movePointerToNextToken(input)
		err = parseI(whileBlock.ChildNodes[len(whileBlock.ChildNodes)-1], input)
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
		outputBlock := common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ASTNode{},
		}
		movePointerToNextToken(input)
		err := parseE(outputBlock, input)
		if err != nil {
			return err
		}
		output.ChildNodes = append(output.ChildNodes, outputBlock)
		return nil
	}
	return errors.New("unexpected parse token in I1")
}

func parseI4(output common.ASTNode, input <-chan common.Token) error {
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

func parseI6(output common.ASTNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent {
		// I6 -> v=E
		_, ok := IdentTable[currPointer.Token]
		if ok {
			return errors.New("variable already declared")
		}
		IdentTable[currPointer.Token] = false
		output.ChildNodes = append(output.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.ASTNode{
				{
					IsLeaf: true,
					InnerToken: common.Token{
						TokenKind: currPointer.TokenKind,
						Token:     currPointer.Token,
					},
					ChildNodes: []common.ASTNode{},
				},
			},
		})
		movePointerToNextToken(input)
		if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' expected")
		}
		movePointerToNextToken(input)
		return parseE(output.ChildNodes[len(output.ChildNodes)-1], input)
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
		identifier := common.ASTNode{
			IsLeaf: true,
			InnerToken: common.Token{
				TokenKind: currPointer.TokenKind,
				Token:     currPointer.Token,
			},
			ChildNodes: []common.ASTNode{},
		}
		movePointerToNextToken(input)
		if currPointer.TokenKind == common.TokenLineEnd {
			return nil
		} else if currPointer.TokenKind != common.TokenAssignment {
			return errors.New("'=' or ';' expected")
		}
		output.ChildNodes = append(output.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			},
			ChildNodes: []common.ASTNode{
				identifier,
			},
		})
		movePointerToNextToken(input)
		return parseE(output.ChildNodes[len(output.ChildNodes)-1], input)
	}
	return errors.New("unexpected parse token in I6")
}

func parseI7(output common.ASTNode, input <-chan common.Token) error {
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
		output.ChildNodes = append(output.ChildNodes, common.ASTNode{
			IsLeaf: false,
			InnerToken: common.Token{
				TokenKind: common.TokenBlock,
				Token:     "",
			},
			ChildNodes: []common.ASTNode{},
		})
		movePointerToNextToken(input)
		err = parseI(output.ChildNodes[len(output.ChildNodes)-1], input)
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
		output.ChildNodes = append(output.ChildNodes, common.ASTNode{
			IsLeaf:     false,
			InnerToken: common.Token{},
		})
		movePointerToNextToken(input)
		err := parseI(output.ChildNodes[len(output.ChildNodes)-1], input)
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

func parseR(output common.ASTNode, input <-chan common.Token) error {
	if currPointer.TokenKind == common.TokenIdent ||
		currPointer.TokenKind == common.TokenLiteralInt ||
		currPointer.TokenKind == common.TokenOpenParanthesis ||
		currPointer.TokenKind == common.TokenInput {
		err := parseE(output, input)
		if err != nil {
			return err
		}
		return parseR1(output, input)
	}
	return errors.New("unexpected parse token in R")
}

func parseR1(output common.ASTNode, input <-chan common.Token) error {
	return nil
}

func parseE(output common.ASTNode, input <-chan common.Token) error {
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
