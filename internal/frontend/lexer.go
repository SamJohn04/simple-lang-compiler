package frontend

import (
	"bufio"
	"io"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// The basic idea of this function is to accept both stdin and file input as parameters.
func Lexer(reader io.Reader, output chan<- common.Token) {
	scanner := bufio.NewScanner(reader)
	defer close(output)

	for scanner.Scan() {
		line := scanner.Text()
		lexLine(line, output)
	}

	// To denote the end of scanner
	output <- common.Token{
		TokenKind: common.TokenEOF,
		Token:     "end of file",
	}
}

func lexLine(line string, output chan<- common.Token) {
	// Until the length of line is 0, keep calling lexSegment
	for len(line) > 0 {
		op, remainingLine := lexSegment(line)
		if op.TokenKind != common.TokenEmpty {
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
			// since no other cases exist
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
			// since no other cases exist
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
	if isCharacterFromVariable(segment[0]) && !(segment[0] >= '0' && segment[0] <= '9') {
		index := isVariableCharactersUntil(segment)
		return common.Token{
			TokenKind: common.TokenIdent,
			Token:     segment[:index],
		}, segment[index:]
	}

	// string check
	if segment[0] == '"' {
		end := 0
		for i, c := range segment[1:] {
			if c == '"' {
				end = i + 2
				break
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
			Token:     segment[:end],
		}, segment[end:]
	}

	// character check
	if segment[0] == '\'' {
		if len(segment) < 3 {
			return common.Token{
				TokenKind: common.TokenError,
				Token:     segment,
			}, ""
		}
		switch segment[1] {
		case '\'':
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

	// number check
	if segment[0] >= '0' && segment[0] <= '9' {
		index, fullstopIndex := isNumericCharactersUntil(segment)
		if fullstopIndex == -1 {
			return common.Token{
				TokenKind: common.TokenLiteralInt,
				Token:     segment[:index],
			}, segment[index:]
		}
		return common.Token{
			TokenKind: common.TokenLiteralFloat,
			Token:     segment[:index],
		}, segment[index:]
	}

	return common.Token{
		TokenKind: common.TokenError,
		Token:     segment,
	}, ""
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

// Returns the index of the first non-numeric character and the '.' character (if it exists).
// Returns on the second full stop as well as on any non-numeric characters.
// If you need an integer, check the second value. If it is -1, use the first value, otherwise the second value.
func isNumericCharactersUntil(segment string) (int, int) {
	fullstopIndex := -1
	for i, c := range segment {
		if c >= '0' && c <= '9' {
			continue
		} else if c == '.' && fullstopIndex == -1 {
			fullstopIndex = i
			continue
		}
		return i, fullstopIndex
	}
	return len(segment), fullstopIndex
}
