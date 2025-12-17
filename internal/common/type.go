package common

import "fmt"

type TokenKind int

const (
	TokenIdent TokenKind = iota + 1

	TokenLiteralInt

	TokenAssignment

	TokenRelationalEquals
	TokenRelationalNotEquals
	TokenRelationalGreaterThan
	TokenRelationalGreaterThanOrEquals
	TokenRelationalLesserThan
	TokenRelationalLesserThanOrEquals

	TokenExpressionAdd
	TokenExpressionSub
	TokenExpressionMul
	TokenExpressionDiv
	TokenExpressionModulo

	TokenOpenParanthesis
	TokenCloseParanthesis

	TokenIf
	TokenElse

	TokenWhile

	TokenLet
	TokenMutable

	TokenOpenCurly
	TokenCloseCurly

	TokenInput
	TokenOutput

	TokenLineEnd

	TokenEOF

	TokenEmpty
	TokenError

	// parser
	TokenBlock
)

var NameMapWithTokenKind = map[TokenKind]string{
	TokenIdent: "Identifier",

	TokenLiteralInt: "Literal Int",

	TokenAssignment: "Assignment",

	TokenRelationalEquals:              "Relational Equals",
	TokenRelationalNotEquals:           "Relational Not Equals",
	TokenRelationalGreaterThan:         "Relational Greater Than",
	TokenRelationalGreaterThanOrEquals: "Relational Greater Than Or Equals",
	TokenRelationalLesserThan:          "Relational Lesser Than",
	TokenRelationalLesserThanOrEquals:  "Relational Lesser Than Or Equals",

	TokenExpressionAdd:    "Add",
	TokenExpressionSub:    "Sub",
	TokenExpressionMul:    "Mul",
	TokenExpressionDiv:    "Div",
	TokenExpressionModulo: "Modulo",

	TokenOpenParanthesis:  "Open Paranthesis",
	TokenCloseParanthesis: "Close Paranthesis",

	TokenIf:   "if",
	TokenElse: "else",

	TokenWhile: "while",

	TokenLet:     "let",
	TokenMutable: "mut",

	TokenOpenCurly:  "Open Curly Braces",
	TokenCloseCurly: "Close Curly Braces",

	TokenInput:  "input",
	TokenOutput: "output",

	TokenLineEnd: "Line End",

	TokenEOF: "End of File",

	TokenEmpty: "Empty Token",
	TokenError: "Error Token",

	TokenBlock: "Code Block",
}

type Token struct {
	TokenKind TokenKind
	Token     string
}

type GeneratorOutput struct {
	Result string
	Err    error
}

type SyntaxTreeNode struct {
	IsLeaf     bool
	InnerToken Token
	ChildNodes []SyntaxTreeNode
}

func (n SyntaxTreeNode) ShallowCopy() SyntaxTreeNode {
	return SyntaxTreeNode{
		IsLeaf: n.IsLeaf,
		InnerToken: Token{
			TokenKind: n.InnerToken.TokenKind,
			Token:     n.InnerToken.Token,
		},
		ChildNodes: n.ChildNodes,
	}
}

func (n SyntaxTreeNode) Display(start string) {
	fmt.Println(start, NameMapWithTokenKind[n.InnerToken.TokenKind], n.InnerToken.Token)
	for _, t := range n.ChildNodes {
		t.Display(start + "+")
	}
}
