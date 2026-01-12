package backend

import (
	"fmt"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// we are compiling to C
func CodeGenerator(
	input []string,
	identifierTable map[string]common.IdentifierInformation,
) (string, error) {
	codes := strings.Builder{}
	buffer := []string{}
	var err error

	codes.WriteString("#include <stdio.h>\n")
	codes.WriteString("#include <stdbool.h>\n")
	codes.WriteString("int main(){\n")

	for identifier, information := range identifierTable {
		if information.DataType == common.TypedUnkown {
			// this can only be introduced into the program
			// by declaring and not initialising a value
			continue
		}
		fmt.Fprintf(
			&codes,
			"%v %v;",
			common.NameMapWithType[information.DataType],
			identifier,
		)
	}

	for _, line := range input {
		buffer, err = writeCodeForLine(&codes, line, buffer)
		if err != nil {
			return "", err
		}
	}

	codes.WriteString("}")
	return codes.String(), nil
}

func writeCodeForLine(codes *strings.Builder, line string, buffer []string) ([]string, error) {
	if line[len(line)-1] == ':' {
		// label
		// TODO add buffer check
		codes.WriteString(line)
		return []string{}, nil
	}
	words := strings.Split(line, " ")
	switch words[0] {
	case "goto":
		// jump
		codes.WriteString(line)
		codes.WriteString(";")
		return []string{}, nil

	case "if":
		// if R goto L
		fmt.Fprintf(codes, "if (%v) { %v; }", words[1], strings.Join(words[2:], " "))
		return []string{}, nil

	case "param":
		// param t
		buffer = append(buffer, strings.Join(words[1:], " ")) // for strings with spaces
		return buffer, nil

	case "call":
		// call func_name
		fmt.Fprintf(codes, "%v(", words[1])
		for i, b := range buffer {
			codes.WriteString(b)
			if i < len(buffer)-1 {
				codes.WriteString(", ")
			}
		}
		codes.WriteString(");")
		return []string{}, nil
	}

	if words[1] != "=" {
		return []string{}, codeGeneratorError(
			fmt.Sprintf("expected v = ..., found %v instead of =", words[1]),
		)
	}

	if words[2] == "input" {
		// i = input
		fmt.Fprintf(codes, "scanf(\"%%d\", &%v);", words[0])
		return []string{}, nil
	}

	fmt.Fprintf(codes, "%v;", line)
	return []string{}, nil
}

func codeGeneratorError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Code Generator",
		Message:        message,
	}
}
