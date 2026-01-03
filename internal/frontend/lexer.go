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
	} else if segment[0] == ';' {
		return common.Token{
			TokenKind: common.TokenLineEnd,
			Token:     ";",
		}, segment[1:]
	} else if segment[0] == ',' {
		return common.Token{
			TokenKind: common.TokenComma,
			Token:     ",",
		}, segment[1:]
	} else if segment[0] == '(' {
		return common.Token{
			TokenKind: common.TokenOpenParanthesis,
			Token:     "(",
		}, segment[1:]
	} else if segment[0] == ')' {
		return common.Token{
			TokenKind: common.TokenCloseParanthesis,
			Token:     ")",
		}, segment[1:]
	} else if segment[0] == '{' {
		return common.Token{
			TokenKind: common.TokenOpenCurly,
			Token:     "{",
		}, segment[1:]
	} else if segment[0] == '}' {
		return common.Token{
			TokenKind: common.TokenCloseCurly,
			Token:     "}",
		}, segment[1:]
	} else if segment[0] == '+' {
		return common.Token{
			TokenKind: common.TokenExpressionAdd,
			Token:     "+",
		}, segment[1:]
	} else if segment[0] == '-' {
		return common.Token{
			TokenKind: common.TokenExpressionSub,
			Token:     "-",
		}, segment[1:]
	} else if segment[0] == '*' {
		return common.Token{
			TokenKind: common.TokenExpressionMul,
			Token:     "*",
		}, segment[1:]
	} else if segment[0] == '/' {
		return common.Token{
			TokenKind: common.TokenExpressionDiv,
			Token:     "/",
		}, segment[1:]
	} else if segment[0] == '%' {
		return common.Token{
			TokenKind: common.TokenExpressionModulo,
			Token:     "%",
		}, segment[1:]
	} else if segment[0] == '!' {
		return common.Token{
			TokenKind: common.TokenNot,
			Token:     "!",
		}, segment[1:]
	} else if len(segment) >= 2 && segment[:2] == "&&" {
		return common.Token{
			TokenKind: common.TokenAnd,
			Token:     "&&",
		}, segment[2:]
	} else if len(segment) >= 2 && segment[:2] == "||" {
		return common.Token{
			TokenKind: common.TokenOr,
			Token:     "||",
		}, segment[2:]
	} else if len(segment) >= 2 && segment[:2] == "==" {
		return common.Token{
			TokenKind: common.TokenRelationalEquals,
			Token:     "==",
		}, segment[2:]
	} else if len(segment) >= 2 && segment[:2] == "!=" {
		return common.Token{
			TokenKind: common.TokenRelationalNotEquals,
			Token:     "!=",
		}, segment[2:]
	} else if len(segment) >= 2 && segment[:2] == ">=" {
		return common.Token{
			TokenKind: common.TokenRelationalGreaterThanOrEquals,
			Token:     ">=",
		}, segment[2:]
	} else if len(segment) >= 2 && segment[:2] == "<=" {
		return common.Token{
			TokenKind: common.TokenRelationalLesserThanOrEquals,
			Token:     "<=",
		}, segment[2:]
	} else if segment[0] == '>' {
		return common.Token{
			TokenKind: common.TokenRelationalGreaterThan,
			Token:     ">",
		}, segment[1:]
	} else if segment[0] == '<' {
		return common.Token{
			TokenKind: common.TokenRelationalLesserThan,
			Token:     "<",
		}, segment[1:]
	} else if segment[0] == '=' {
		return common.Token{
			TokenKind: common.TokenAssignment,
			Token:     "=",
		}, segment[1:]
	} else if len(segment) >= 2 && segment[:2] == "if" &&
		(len(segment) == 2 || !isCharacterFromVariable(segment[2])) {
		return common.Token{
			TokenKind: common.TokenIf,
			Token:     "if",
		}, segment[2:]
	} else if len(segment) >= 4 && segment[:4] == "else" &&
		(len(segment) == 4 || !isCharacterFromVariable(segment[4])) {
		return common.Token{
			TokenKind: common.TokenElse,
			Token:     "else",
		}, segment[4:]
	} else if len(segment) >= 5 && segment[:5] == "while" &&
		(len(segment) == 5 || !isCharacterFromVariable(segment[5])) {
		return common.Token{
			TokenKind: common.TokenWhile,
			Token:     "while",
		}, segment[5:]
	} else if len(segment) >= 5 && segment[:5] == "input" &&
		(len(segment) == 5 || !isCharacterFromVariable(segment[5])) {
		return common.Token{
			TokenKind: common.TokenInput,
			Token:     "input",
		}, segment[5:]
	} else if len(segment) >= 6 &&
		segment[:6] == "output" && (len(segment) == 6 || !isCharacterFromVariable(segment[6])) {
		return common.Token{
			TokenKind: common.TokenOutput,
			Token:     "output",
		}, segment[6:]
	} else if len(segment) >= 3 &&
		segment[:3] == "let" && (len(segment) == 3 || !isCharacterFromVariable(segment[3])) {
		return common.Token{
			TokenKind: common.TokenLet,
			Token:     "let",
		}, segment[3:]
	} else if len(segment) >= 3 &&
		segment[:3] == "mut" && (len(segment) == 3 || !isCharacterFromVariable(segment[3])) {
		return common.Token{
			TokenKind: common.TokenMutable,
			Token:     "mut",
		}, segment[3:]
	}

	if isCharacterFromVariable(segment[0]) && !(segment[0] >= '0' && segment[0] <= '9') {
		for i, c := range segment {
			if c >= 'A' && c <= 'Z' {
				continue
			} else if c >= 'a' && c <= 'z' {
				continue
			} else if c >= '0' && c <= '9' {
				continue
			} else if c == '_' {
				continue
			}
			return common.Token{
				TokenKind: common.TokenIdent,
				Token:     segment[:i],
			}, segment[i:]
		}
		return common.Token{
			TokenKind: common.TokenIdent,
			Token:     segment,
		}, ""
	}

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
			Token:     segment[0:end],
		}, segment[end:]
	}

	if segment[0] >= '0' && segment[0] <= '9' {
		for i, c := range segment {
			if c >= '0' && c <= '9' {
				continue
			}
			return common.Token{
				TokenKind: common.TokenLiteralInt,
				Token:     segment[:i],
			}, segment[i:]
		}
		return common.Token{
			TokenKind: common.TokenLiteralInt,
			Token:     segment,
		}, ""
	}

	return common.Token{
		TokenKind: common.TokenError,
		Token:     segment,
	}, ""
}

func isCharacterFromVariable(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_'
}
