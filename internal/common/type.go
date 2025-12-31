package common

import "fmt"

type TokenKind int

const (
	TokenIdent TokenKind = iota + 1

	TokenLiteralInt

	TokenLiteralBool   // WARN not implemented yet
	TokenLiteralChar   // WARN not implemented yet
	TokenLiteralFloat  // WARN not implemented yet
	TokenLiteralString // WARN not implemented yet

	TokenAssignment

	TokenOr  // WARN not implemented yet
	TokenAnd // WARN not implemented yet
	TokenNot // WARN not implemented yet

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

	// used by the parser, usually as a grouper
	// specifics about the information in the Token string
	TokenBlock
)

var NameMapWithTokenKind = map[TokenKind]string{
	TokenIdent: "Identifier",

	TokenLiteralInt: "Literal Int",

	TokenLiteralBool:   "Literal Bool",   // WARN not implemented yet
	TokenLiteralChar:   "Literal Char",   // WARN not implemented yet
	TokenLiteralFloat:  "Literal Float",  // WARN not implemented yet
	TokenLiteralString: "Literal String", // WARN not implemented yet

	TokenAssignment: "Assignment",

	TokenOr:  "OR",  // WARN not implemented yet
	TokenAnd: "AND", // WARN not implemented yet
	TokenNot: "NOT", // WARN not implemented yet

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

var Operators = map[TokenKind]string{
	TokenExpressionAdd:    "+",
	TokenExpressionSub:    "-",
	TokenExpressionMul:    "*",
	TokenExpressionDiv:    "/",
	TokenExpressionModulo: "%",
}

type DataTypeOfIdentifier int

const (
	TypedInt    DataTypeOfIdentifier = iota + 1
	TypedBool                        // WARN not implemented yet
	TypedChar                        // WARN not implemented yet
	TypedFloat                       // WARN not implemented yet
	TypedString                      // WARN not implemented yet
	TypedVoid                        // WARN not implemented yet

	TypedUnkown // What mut variables get typed as, before assignment
)

type Token struct {
	TokenKind TokenKind
	Token     string
}

type SyntaxTreeNode struct {
	InnerToken Token
	ChildNodes []SyntaxTreeNode
}

func (n SyntaxTreeNode) ShallowCopy() SyntaxTreeNode {
	return SyntaxTreeNode{
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

type UnderConstructionError struct {
	PointOfFailure string
	Message        string
}

func (e *UnderConstructionError) Error() string {
	errMsg := e.PointOfFailure + " is still under construction"
	if e.Message != "" {
		errMsg += ": " + e.Message
	}
	return errMsg
}

// This error is returned when something goes wrong with the compiler.
// This usually points to a mistake from my side.
type InternalError struct {
	PointOfFailure string
	Message        string
}

func (e *InternalError) Error() string {
	return "There is something wrong with the compiler.\n" +
		"It is recommended to contact someone in the language creation team regarding this issue.\n" +
		"\t" + e.Message + "\n" +
		"in " + e.PointOfFailure
}

// This error is returned when something goes wrong with the user input.
// This usually points to a mistake on the compiled code's side
type CompilationError struct {
	PointOfFailure string
	Message        string
}

func (e *CompilationError) Error() string {
	return e.PointOfFailure + ": " + e.Message
}
