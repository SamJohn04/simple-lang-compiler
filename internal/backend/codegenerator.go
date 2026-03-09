package backend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// we are compiling to C
func CodeGenerator(
	input []string,
	identifiers []common.IdentifierInformation,
) (string, error) {
	codes := strings.Builder{}
	buffer := []string{}
	var err error

	writeStart(&codes)

	codes.WriteString("int main() {\n\t")

	for index, information := range identifiers {
		if information.Datatype == nil || information.Datatype.IsDatatype(common.TypedUnknown) {
			// this can only be introduced into the program
			// by declaring and not initialising a value
			continue
		}
		datatype, length, err := information.Datatype.ToString()
		if err != nil {
			fmt.Printf("WARN: %v", err)
			continue
		}
		if length <= 1 {
			fmt.Fprintf(&codes, "%v _t%v;\n\t", datatype, index)
			continue
		}
		fmt.Fprintf(
			&codes,
			"%v _arr%v[%v];\n\t",
			datatype,
			index,
			length,
		)
		fmt.Fprintf(&codes, "%v* _t%v = _arr%v;", datatype, index, index)
		codes.WriteString("\n\t")
	}

	for _, line := range input {
		buffer, err = writeCodeForLine(&codes, line, buffer, identifiers)
		codes.WriteString("\n\t")
		if err != nil {
			return "", err
		}
	}

	codes.WriteString("\n}")
	return codes.String(), nil
}

func writeCodeForLine(
	codes *strings.Builder,
	line string,
	buffer []string,
	identifiers []common.IdentifierInformation,
) ([]string, error) {
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
			codes.Write([]byte(b))
			if i < len(buffer)-1 {
				fmt.Fprint(codes, ", ")
			}
		}
		fmt.Fprint(codes, ");")
		return []string{}, nil
	}

	if words[1] == "[]" {
		fmt.Fprintf(codes, "%v[%v] = %v;", words[0], words[2], strings.Join(words[4:], " "))
		return []string{}, nil
	}

	if words[1] != "=" {
		return []string{}, codeGeneratorError(
			fmt.Sprintf("expected v = ..., found %v instead of =", words[1]),
		)
	}

	if words[2] == "input" {
		// i = input
		fmt.Fprintf(codes, "scanf(\"%%lld\", &%v);", words[0])
		return []string{}, nil
	}

	_, length, err := identifiers[indexFromIdentifier(words[0])].Datatype.ToString()
	if err != nil {
		return []string{}, err
	}
	if length > 1 {
		datatype := identifiers[indexFromIdentifier(words[0])].Datatype
		destination := words[0]
		source := words[2]
		if len(words) > 3 && words[3] == "[]" {
			source = fmt.Sprintf("(%v + %v)", source, words[4])
		}
		fmt.Fprintf(
			codes,
			"copy__%v(%v, %v, %v);",
			datatype.ToRepresentation(),
			destination,
			source,
			length,
		)
		return []string{}, nil
	}

	if len(words) > 3 && words[3] == "[]" {
		fmt.Fprintf(codes, "%v = %v[%v];", words[0], words[2], words[4])
		return []string{}, nil
	}

	fmt.Fprintf(codes, "%v;", line)
	return []string{}, nil
}

func writeStart(codes *strings.Builder) {
	codes.WriteString("#include <stdio.h>\n")
	codes.WriteString("#include <stdbool.h>\n")
	codes.WriteString("#include <string.h>\n")

	codes.WriteString(`
void copy__str(char** dest, char** src, long long length) {
	for (long long i = 0; i < length; i++) {
		strncpy(dest[i], src[i], 1022);
		dest[i][1023] = '\0';
	}
}
	`)

	representations := map[string]string{
		"bool":      "b",
		"char":      "c",
		"long long": "l",
		"double":    "d",
	}

	for datatype, representation := range representations {
		fmt.Fprintf(
			codes,
			`
void copy__%v(%v* dest, %v* src, long long length) {
	for (long long i = 0; i < length; i++) {
		dest[i] = src[i];
	}
}`, representation, datatype, datatype,
		)
		codes.WriteString("\n")
	}
}

func indexFromIdentifier(identifier string) int {
	i, _ := strconv.Atoi(identifier[2:])
	return i
}

func codeGeneratorError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Code Generator",
		Message:        message,
	}
}
