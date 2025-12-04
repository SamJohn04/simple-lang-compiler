package common

type Token struct {
	TokenKind TokenKind
	Token     string
}

type ASTNode struct {
	IsLeaf     bool
	InnerToken Token
	ChildNodes []ASTNode
}

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
