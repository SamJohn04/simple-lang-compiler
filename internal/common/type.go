package common

import "fmt"

type TokenKind int

const (
	TokenIdent TokenKind = iota + 1

	TokenLiteralInt
	TokenLiteralString

	TokenLiteralBool
	TokenLiteralChar
	TokenLiteralFloat

	TokenAssignment

	TokenOr
	TokenAnd
	TokenNot

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

	TokenOpenSquareBraces
	TokenCloseSquareBraces

	TokenInput
	TokenOutput

	TokenComma

	TokenLineEnd

	TokenEOF

	TokenEmpty
	TokenError

	// used by the parser, usually as a grouper
	// specifics about the information in the Token string
	TokenBlock
	// used by typechecker to show that a variable is declared and not initialized here
	TokenDeclare
)

var NameMapWithTokenKind = map[TokenKind]string{
	TokenIdent: "Identifier",

	TokenLiteralInt:    "Literal Int",
	TokenLiteralString: "Literal String",

	TokenLiteralBool:  "Literal Bool",
	TokenLiteralChar:  "Literal Char",
	TokenLiteralFloat: "Literal Float",

	TokenAssignment: "Assignment",

	TokenOr:  "OR",
	TokenAnd: "AND",
	TokenNot: "NOT",

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

	TokenOpenSquareBraces:  "Open Square Braces",
	TokenCloseSquareBraces: "Close Square Braces",

	TokenInput:  "input",
	TokenOutput: "output",

	TokenComma: "comma",

	TokenLineEnd: "Line End",

	TokenEOF: "End of File",

	TokenEmpty: "Empty Token",
	TokenError: "Error Token",

	TokenBlock:   "Code Block",
	TokenDeclare: "declare",
}

var Operators = map[TokenKind]string{
	TokenExpressionAdd:    "+",
	TokenExpressionSub:    "-",
	TokenExpressionMul:    "*",
	TokenExpressionDiv:    "/",
	TokenExpressionModulo: "%",

	TokenRelationalEquals:              "==",
	TokenRelationalGreaterThan:         ">",
	TokenRelationalGreaterThanOrEquals: ">=",
	TokenRelationalNotEquals:           "!=",
	TokenRelationalLesserThan:          "<",
	TokenRelationalLesserThanOrEquals:  "<=",

	TokenOr:  "||",
	TokenAnd: "&&",
	TokenNot: "!",
}

type DataTypeOfIdentifier int

const (
	TypedUnkown DataTypeOfIdentifier = iota // What mut variables get typed as, before assignment

	TypedInt
	TypedBool
	TypedChar
	TypedFloat
	TypedString // WARN not implemented yet
	TypedVoid   // WARN not implemented yet
)

var NameMapWithType = map[DataTypeOfIdentifier]string{
	TypedUnkown: "Unknown Type",

	TypedInt:    "int",
	TypedBool:   "bool",
	TypedChar:   "char",
	TypedFloat:  "float",
	TypedString: "string",
	TypedVoid:   "void",
}

type Token struct {
	TokenKind TokenKind
	Token     string
}

type SyntaxTreeNode struct {
	InnerToken Token
	ChildNodes []SyntaxTreeNode
	Datatype   DataTypeOfIdentifier
}

func (n SyntaxTreeNode) ShallowCopy() SyntaxTreeNode {
	return SyntaxTreeNode{
		InnerToken: Token{
			TokenKind: n.InnerToken.TokenKind,
			Token:     n.InnerToken.Token,
		},
		ChildNodes: n.ChildNodes,
		Datatype:   n.Datatype,
	}
}

func (n SyntaxTreeNode) Display(start string) {
	fmt.Println(start, NameMapWithTokenKind[n.InnerToken.TokenKind], n.InnerToken.Token, NameMapWithType[n.Datatype])
	for _, t := range n.ChildNodes {
		t.Display(start + start)
	}
}

type IdentifierInformation struct {
	DataType DataTypeOfIdentifier
	Mutable  bool
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
