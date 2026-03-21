package common

type Token struct {
	LineNumber int
	TokenKind  TokenKind
	Token      string
}

type TokenKind int

const (
	// Identifier
	// may contain a literal (such as an int or a string), array, or a function name
	TokenIdent TokenKind = iota + 1

	// integers
	TokenLiteralInt
	// strings
	TokenLiteralString
	// booleans
	TokenLiteralBool
	// characters
	TokenLiteralChar
	// floating point numbers
	TokenLiteralFloat

	// = symbol
	TokenAssignment

	// usable on two boolean expressions to give a boolean value

	TokenOr
	TokenAnd

	// useable on a boolean expression to give a boolean value
	TokenNot

	// usable on any two values and gives a boolean value

	TokenRelationalEquals
	TokenRelationalNotEquals

	// TODO clarify if/how they work when comparing two values of different types
	// theory: int, float, char values work normally with each other; characters taking ascii values
	//	boolean values take 0 for false and 1 for true
	//	other types do not work

	TokenRelationalGreaterThan
	TokenRelationalGreaterThanOrEquals
	TokenRelationalLesserThan
	TokenRelationalLesserThanOrEquals

	// mathematical expressions
	// only works on int, float, and char; characters are given their corresponding ascii values

	TokenExpressionAdd
	TokenExpressionSub
	TokenExpressionMul
	TokenExpressionDiv
	TokenExpressionModulo

	TokenOpenParanthesis
	TokenCloseParanthesis

	// if condition
	TokenIf
	// else condition
	// else if cases are marked by an else token followed by an if token
	TokenElse

	// while condition
	TokenWhile

	TokenLet
	TokenMutable

	TokenOpenCurly
	TokenCloseCurly

	TokenOpenSquareBraces
	TokenCloseSquareBraces

	// TODO turn these into functions instead
	// get user input
	TokenInput
	// give user output
	TokenOutput

	// , used to separate lists, arguments, etc.
	TokenComma

	// ; showing the end of the current instruction
	TokenLineEnd

	// Token showing the end of the file
	TokenEOF

	// used by the lexer

	// return nothing in the Lexer
	// if, for example, the rest of the line is a comment
	TokenEmpty
	// in case something goes wrong
	// e.g. $ outside quotes
	TokenError

	// used by the parser, usually as a grouper
	// specifics about the information in the Token string
	TokenBlock
)

var NameMapWithTokenKind = map[TokenKind]string{
	TokenIdent: "Identifier",

	TokenLiteralInt:    "Literal Int",
	TokenLiteralString: "Literal String",
	TokenLiteralBool:   "Literal Bool",
	TokenLiteralChar:   "Literal Char",
	TokenLiteralFloat:  "Literal Float",

	TokenAssignment: "Assignment =",

	TokenOr:  "OR",
	TokenAnd: "AND",
	TokenNot: "NOT",

	TokenRelationalEquals:              "Relational Equals ==",
	TokenRelationalNotEquals:           "Relational Not Equals !=",
	TokenRelationalGreaterThan:         "Relational Greater Than >",
	TokenRelationalGreaterThanOrEquals: "Relational Greater Than Or Equals >=",
	TokenRelationalLesserThan:          "Relational Lesser Than <",
	TokenRelationalLesserThanOrEquals:  "Relational Lesser Than Or Equals <=",

	TokenExpressionAdd:    "Add",
	TokenExpressionSub:    "Sub",
	TokenExpressionMul:    "Mul",
	TokenExpressionDiv:    "Div",
	TokenExpressionModulo: "Modulo",

	TokenOpenParanthesis:  "Open Paranthesis (",
	TokenCloseParanthesis: "Close Paranthesis )",

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

	TokenBlock: "Code Block",
}

type IdentifierInformation struct {
	IdentifierName string
	Datatype       Datatype
	Mutable        bool
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
