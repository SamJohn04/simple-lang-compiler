package frontend

import (
	"bufio"
	"io"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// The basic idea of this function's first parameter as io.Reader
// is to accept both stdin and file input as parameters.
func Lexer(reader io.Reader, output chan<- common.Token) {
	defer close(output)
	scanner := bufio.NewScanner(reader)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber += 1
		line := scanner.Text()
		lexLine(line, lineNumber, output)
	}

	// To denote the end of scanner
	output <- common.Token{
		TokenKind: common.TokenEOF,
		Token:     "Lexer: end of file",
	}
}

func lexLine(line string, lineNumber int, output chan<- common.Token) {
	for len(line) > 0 {
		op, remainingLine := lexSegment(line)
		// TokenEmpty is sent in case the remaining string has no meaningful components
		// We do not wish to propogate such cases any further
		if op.TokenKind != common.TokenEmpty {
			// set LineNumber here
			op.LineNumber = lineNumber
			output <- op
		}
		line = remainingLine
	}
}

func lexSegment(segment string) (common.Token, string) {
	// Remove all leading and ending \t
	segment = strings.Trim(segment, "\t")
	// Remove all leading and ending ' '
	segment = strings.Trim(segment, " ")

	if len(segment) == 0 || len(segment) >= 2 && segment[:2] == "//" {
		return common.Token{
			TokenKind: common.TokenEmpty,
			Token:     "",
		}, ""
	}

	switch segment[0] {
	case ';':
		return common.Token{
			TokenKind: common.TokenLineEnd,
			Token:     ";",
		}, segment[1:]

	case ',':
		return common.Token{
			TokenKind: common.TokenComma,
			Token:     ",",
		}, segment[1:]

	case '(':
		return common.Token{
			TokenKind: common.TokenOpenParanthesis,
			Token:     "(",
		}, segment[1:]

	case ')':
		return common.Token{
			TokenKind: common.TokenCloseParanthesis,
			Token:     ")",
		}, segment[1:]

	case '{':
		return common.Token{
			TokenKind: common.TokenOpenCurly,
			Token:     "{",
		}, segment[1:]

	case '}':
		return common.Token{
			TokenKind: common.TokenCloseCurly,
			Token:     "}",
		}, segment[1:]

	case '[':
		return common.Token{
			TokenKind: common.TokenOpenSquareBraces,
			Token:     "[",
		}, segment[1:]

	case ']':
		return common.Token{
			TokenKind: common.TokenCloseSquareBraces,
			Token:     "]",
		}, segment[1:]

	case '+':
		return common.Token{
			TokenKind: common.TokenExpressionAdd,
			Token:     "+",
		}, segment[1:]

	case '-':
		return common.Token{
			TokenKind: common.TokenExpressionSub,
			Token:     "-",
		}, segment[1:]

	case '*':
		return common.Token{
			TokenKind: common.TokenExpressionMul,
			Token:     "*",
		}, segment[1:]

	case '/':
		return common.Token{
			TokenKind: common.TokenExpressionDiv,
			Token:     "/",
		}, segment[1:]

	case '%':
		return common.Token{
			TokenKind: common.TokenExpressionModulo,
			Token:     "%",
		}, segment[1:]

	case '=':
		if len(segment) < 2 || segment[1] != '=' {
			return common.Token{
				TokenKind: common.TokenAssignment,
				Token:     "=",
			}, segment[1:]
		}
		return common.Token{
			TokenKind: common.TokenRelationalEquals,
			Token:     "==",
		}, segment[2:]

	case '<':
		if len(segment) < 2 || segment[1] != '=' {
			return common.Token{
				TokenKind: common.TokenRelationalLesserThan,
				Token:     "<",
			}, segment[1:]
		}
		return common.Token{
			TokenKind: common.TokenRelationalLesserThanOrEquals,
			Token:     "<=",
		}, segment[2:]

	case '>':
		if len(segment) < 2 || segment[1] != '=' {
			return common.Token{
				TokenKind: common.TokenRelationalGreaterThan,
				Token:     ">",
			}, segment[1:]
		}
		return common.Token{
			TokenKind: common.TokenRelationalGreaterThanOrEquals,
			Token:     ">=",
		}, segment[2:]

	case '!':
		if len(segment) < 2 || segment[1] != '=' {
			return common.Token{
				TokenKind: common.TokenNot,
				Token:     "!",
			}, segment[1:]
		}
		return common.Token{
			TokenKind: common.TokenRelationalNotEquals,
			Token:     "!=",
		}, segment[2:]

	case '&':
		if len(segment) < 2 || segment[1] != '&' {
			// since a single & symbol is meaningless
			return common.Token{
				TokenKind: common.TokenError,
				Token:     segment,
			}, ""
		}
		return common.Token{
			TokenKind: common.TokenAnd,
			Token:     "&&",
		}, segment[2:]

	case '|':
		if len(segment) < 2 || segment[1] != '|' {
			// since a single | symbol is meaningless
			return common.Token{
				TokenKind: common.TokenError,
				Token:     segment,
			}, ""
		}
		return common.Token{
			TokenKind: common.TokenOr,
			Token:     "||",
		}, segment[2:]

	case 'i':
		if isWordToken(segment, "if") {
			return common.Token{
				TokenKind: common.TokenIf,
				Token:     "if",
			}, segment[2:]
		} else if isWordToken(segment, "input") {
			return common.Token{
				TokenKind: common.TokenInput,
				Token:     "input",
			}, segment[5:]
		}

	case 'e':
		if isWordToken(segment, "else") {
			return common.Token{
				TokenKind: common.TokenElse,
				Token:     "else",
			}, segment[4:]
		}

	case 'w':
		if isWordToken(segment, "while") {
			return common.Token{
				TokenKind: common.TokenWhile,
				Token:     "while",
			}, segment[5:]
		}

	case 'o':
		if isWordToken(segment, "output") {
			return common.Token{
				TokenKind: common.TokenOutput,
				Token:     "output",
			}, segment[6:]
		}

	case 'l':
		if isWordToken(segment, "let") {
			return common.Token{
				TokenKind: common.TokenLet,
				Token:     "let",
			}, segment[3:]
		}

	case 'm':
		if isWordToken(segment, "mut") {
			return common.Token{
				TokenKind: common.TokenMutable,
				Token:     "mut",
			}, segment[3:]
		}

	case 't':
		if isWordToken(segment, "true") {
			return common.Token{
				TokenKind: common.TokenLiteralBool,
				Token:     "true",
			}, segment[4:]
		}

	case 'f':
		if isWordToken(segment, "false") {
			return common.Token{
				TokenKind: common.TokenLiteralBool,
				Token:     "false",
			}, segment[5:]
		}
	}

	// variable check
	if segment[0] >= 'A' && segment[0] <= 'Z' ||
		segment[0] >= 'a' && segment[0] <= 'z' ||
		segment[0] == '_' {
		return lexVariable(segment)
	}

	// string check
	if segment[0] == '"' {
		return lexString(segment)
	}

	// character check
	if segment[0] == '\'' {
		return lexChar(segment)
	}

	// number check
	if segment[0] >= '0' && segment[0] <= '9' {
		return lexNumber(segment)
	}

	return common.Token{
		TokenKind: common.TokenError,
		Token:     segment,
	}, ""
}

func lexVariable(segment string) (common.Token, string) {
	index := isVariableCharactersUntil(segment)
	return common.Token{
		TokenKind: common.TokenIdent,
		Token:     segment[:index],
	}, segment[index:]
}

func lexString(segment string) (common.Token, string) {
	end := 0
	escapeFromNextCharacter := false
	for i, c := range segment[1:] {
		if c == '"' && !escapeFromNextCharacter {
			// + 1 because we are starting the loop from segment[1]
			end = i + 1
			break
		} else if c == '\\' && !escapeFromNextCharacter {
			escapeFromNextCharacter = true
		} else {
			escapeFromNextCharacter = false
		}
	}
	if end == 0 {
		return common.Token{
			TokenKind: common.TokenError,
			Token:     "\" is not closed",
		}, ""
	}
	return common.Token{
		TokenKind: common.TokenLiteralString,
		Token:     segment[:end+1],
	}, segment[end+1:]
}

func lexChar(segment string) (common.Token, string) {
	if len(segment) < 3 {
		return common.Token{
			TokenKind: common.TokenError,
			Token:     segment,
		}, ""
	}
	switch segment[1] {
	case '\'':
		// '' (0 characters) is invalid
		return common.Token{
			TokenKind: common.TokenError,
			Token:     segment,
		}, ""

	case '\\':
		// verify the closing quotes
		if len(segment) < 4 || segment[3] != '\'' {
			return common.Token{
				TokenKind: common.TokenError,
				Token:     segment,
			}, ""
		}
		return common.Token{
			TokenKind: common.TokenLiteralChar,
			Token:     segment[:4],
		}, segment[4:]

	default:
		if segment[2] != '\'' {
			return common.Token{
				TokenKind: common.TokenError,
				Token:     segment,
			}, ""
		}
		return common.Token{
			TokenKind: common.TokenLiteralChar,
			Token:     segment[:3],
		}, segment[3:]
	}
}

func lexNumber(segment string) (common.Token, string) {
	index := isNumberUntil(segment)
	if index == len(segment) || segment[index] != '.' {
		return common.Token{
			TokenKind: common.TokenLiteralInt,
			Token:     segment[:index],
		}, segment[index:]
	}
	floatingPointIndex := isNumberUntil(segment[index+1:])
	if floatingPointIndex == 0 {
		// no number after floating point
		return common.Token{
			TokenKind: common.TokenError,
			Token:     segment,
		}, ""
	}
	return common.Token{
		TokenKind: common.TokenLiteralFloat,
		Token:     segment[:floatingPointIndex+index+1],
	}, segment[floatingPointIndex+index+1:]
}

// Intended for multiple character tokens like if, mut, etc.
// If the whole string matches, checks to see if the next character is not a variable token.
func isWordToken(segment, token string) bool {
	if !strings.HasPrefix(segment, token) {
		return false
	}
	return len(segment) == len(token) || !isCharacterFromVariable(segment[len(token)])
}

// Returns the index of the first non-variable character.
// Does not grant any special consideration to the first character,
// i.e., the variable may be starting with a number.
// Add a c < 0 || c > 9 check for the first character if necessary.
func isVariableCharactersUntil(segment string) int {
	for i := range len(segment) { // we cannot use i, c := range segment because of rune vs byte issues
		if !isCharacterFromVariable(segment[i]) {
			return i
		}
	}
	return len(segment)
}

func isCharacterFromVariable(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_'
}

// Returns the index of the first non-numeric character
func isNumberUntil(segment string) int {
	for i, c := range segment {
		if c < '0' || c > '9' {
			return i
		}
	}
	return len(segment)
}
