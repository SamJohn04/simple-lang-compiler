package common

type Token struct {
	TokenKind TokenKind
	Token     string
}

type TokenKind int

const (
	TokenIdent TokenKind = iota + 1
	TokenLabel

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
	TokenThen
	TokenElse

	TokenGoto

	TokenInput
	TokenOutput

	TokenLineEnd

	TokenEmpty
	TokenError
)

var NameMapWithTokenKind = map[TokenKind]string{
	TokenIdent: "Identifier",
	TokenLabel: "Label",

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
	TokenThen: "then",
	TokenElse: "else",

	TokenGoto: "goto",

	TokenInput:  "input",
	TokenOutput: "output",

	TokenLineEnd: "Line End",

	TokenEmpty: "Empty Token",
	TokenError: "Error Token",
}
