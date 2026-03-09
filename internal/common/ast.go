package common

import (
	"errors"
	"fmt"
)

type ProgramAST struct {
	Instructions []InstructionAST
}

func (p ProgramAST) PerformAllChecks(identifiers []IdentifierInformation) error {
	for _, instruction := range p.Instructions {
		err := instruction.PerformChecks(identifiers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p ProgramAST) ThreeAddressCode(
	identifiers []IdentifierInformation, numberOfGotos *int,
) ([]string, []IdentifierInformation, error) {
	threeAddressCodes := []string{}
	for _, instruction := range p.Instructions {
		codes, ident, err := instruction.ThreeAddressCode(identifiers, numberOfGotos)
		if err != nil {
			return []string{}, identifiers, err
		}
		threeAddressCodes = append(threeAddressCodes, codes...)
		identifiers = ident
	}
	return threeAddressCodes, identifiers, nil
}

func (p ProgramAST) Display() {
	for _, instruction := range p.Instructions {
		fmt.Println(instruction)
	}
}

type InstructionAST interface {
	PerformChecks(identifiers []IdentifierInformation) error
	ThreeAddressCode(
		identifiers []IdentifierInformation,
		numberOfGotos *int,
	) ([]string, []IdentifierInformation, error)
}

type AssignmentAST struct {
	AssignToIdentifier int
	ArrayValues        []ExpressionAST
	AssignValue        ExpressionAST
}

func (a AssignmentAST) PerformChecks(identifiers []IdentifierInformation) error {
	assignedDatatype, err := a.AssignValue.GetDatatype(identifiers)
	if err != nil {
		return err
	}

	identifierDatatype := identifiers[a.AssignToIdentifier].Datatype
	if len(a.ArrayValues) > 0 &&
		(identifierDatatype == nil || identifierDatatype.IsDatatype(TypedUnknown)) {
		return errors.New("undeclared identifier cannot have array accesses")
	}
	if identifierDatatype == nil || identifierDatatype.IsDatatype(TypedUnknown) {
		identifiers[a.AssignToIdentifier].Datatype = assignedDatatype
		return nil
	}

	arrayDatatype, ok := identifierDatatype.(ArrayDatatype)
	for range a.ArrayValues {
		if !ok {
			return errors.New("more array accesses than nested arrays")
		}
		identifierDatatype = arrayDatatype.ElementType
		arrayDatatype, ok = identifierDatatype.(ArrayDatatype)
	}

	if ok {
		return errors.New("assignment of an array in an array is not possible")
	}
	if !identifierDatatype.IsDatatype(assignedDatatype) {
		return errors.New("identifier datatype and operand datatype do not match")
	}
	return nil
}

func (a AssignmentAST) ThreeAddressCode(
	identifiers []IdentifierInformation,
	numberOfGotos *int,
) ([]string, []IdentifierInformation, error) {
	result, codes, identifiers, err := a.AssignValue.ThreeAddressCode(identifiers)
	if err != nil {
		return codes, identifiers, err
	}

	if len(a.ArrayValues) == 0 {
		codes = append(
			codes,
			fmt.Sprintf("%v = %v", identifierFromIndex(a.AssignToIdentifier), result),
		)
		return codes, identifiers, nil
	}

	arrayDatatype, arrayOk := identifiers[a.AssignToIdentifier].Datatype.(ArrayDatatype)
	arrayResult := "0"
	for _, access := range a.ArrayValues {
		if !arrayOk {
			return codes, identifiers, errors.New("non-array where array expected")
		}

		assignTo, arrayCodes, ident, err := access.ThreeAddressCode(identifiers)
		if err != nil {
			return codes, identifiers, err
		}

		codes = append(codes, arrayCodes...)

		variable1, ident := nextIdentifier(ident, TypedInt)
		variable2, ident := nextIdentifier(ident, TypedInt)
		codes = append(
			codes,
			fmt.Sprintf("%v = %v * %v", variable1, arrayResult, arrayDatatype.NumberOfElements),
			fmt.Sprintf("%v = %v + %v", variable2, variable1, assignTo),
		)

		arrayResult = variable2
		identifiers = ident
		arrayDatatype, arrayOk = arrayDatatype.ElementType.(ArrayDatatype)
	}
	codes = append(codes, fmt.Sprintf(
		"%v [] %v = %v",
		identifierFromIndex(a.AssignToIdentifier),
		arrayResult,
		result,
	))
	return codes, identifiers, nil
}

type IfStatementAST struct {
	IfExpressions []IfExpression
}

type IfExpression struct {
	Condition ExpressionAST
	Program   ProgramAST
}

func (i IfStatementAST) PerformChecks(identifiers []IdentifierInformation) error {
	for _, expression := range i.IfExpressions {
		conditionDatatype, err := expression.Condition.GetDatatype(identifiers)
		if err != nil {
			return err
		}
		if !conditionDatatype.IsDatatype(TypedBool) {
			return errors.New("non-boolean value in an if condition")
		}
		err = expression.Program.PerformAllChecks(identifiers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i IfStatementAST) ThreeAddressCode(
	identifiers []IdentifierInformation,
	numberOfGotos *int,
) ([]string, []IdentifierInformation, error) {
	codes := []string{}
	ifEndGoto := getNextGoto(numberOfGotos)

	for _, ifExpression := range i.IfExpressions {
		variable, conditionCodes, id, err := ifExpression.Condition.ThreeAddressCode(
			identifiers,
		)
		if err != nil {
			return []string{}, identifiers, err
		}

		programCodes, id, err := ifExpression.Program.ThreeAddressCode(id, numberOfGotos)
		if err != nil {
			return []string{}, identifiers, err
		}

		holdGoto := getNextGoto(numberOfGotos)
		nextGoto := getNextGoto(numberOfGotos)

		codes = append(codes, conditionCodes...)
		codes = append(
			codes,
			fmt.Sprintf("if %v goto %v", variable, holdGoto),
			fmt.Sprintf("goto %v", nextGoto),
			fmt.Sprintf("%v:", holdGoto),
		)
		codes = append(codes, programCodes...)
		codes = append(
			codes,
			fmt.Sprintf("goto %v", ifEndGoto),
			fmt.Sprintf("%v:", nextGoto),
		)
		identifiers = id
	}
	codes = append(codes, fmt.Sprintf("%v:", ifEndGoto))

	return codes, identifiers, nil
}

type WhileStatementAST struct {
	Condition ExpressionAST
	Program   ProgramAST
}

func (w WhileStatementAST) PerformChecks(identifiers []IdentifierInformation) error {
	conditionDatatype, err := w.Condition.GetDatatype(identifiers)
	if err != nil {
		return err
	}
	if !conditionDatatype.IsDatatype(TypedBool) {
		return errors.New("non-boolean value in a while condition")
	}
	err = w.Program.PerformAllChecks(identifiers)
	if err != nil {
		return err
	}
	return nil
}

func (w WhileStatementAST) ThreeAddressCode(
	identifiers []IdentifierInformation,
	numberOfGotos *int,
) ([]string, []IdentifierInformation, error) {
	whileGoto := getNextGoto(numberOfGotos)
	holdGoto := getNextGoto(numberOfGotos)
	nextGoto := getNextGoto(numberOfGotos)
	threeAddressCodes := []string{
		fmt.Sprintf("%v:", whileGoto),
	}

	relation, relationCodes, identifiers, err := w.Condition.ThreeAddressCode(identifiers)
	if err != nil {
		return []string{}, identifiers, err
	}

	programCodes, identifiers, err := w.Program.ThreeAddressCode(identifiers, numberOfGotos)
	if err != nil {
		return []string{}, identifiers, err
	}

	threeAddressCodes = append(threeAddressCodes, relationCodes...)
	threeAddressCodes = append(
		threeAddressCodes,
		fmt.Sprintf("if %v goto %v", relation, holdGoto),
		fmt.Sprintf("goto %v", nextGoto),
		fmt.Sprintf("%v:", holdGoto),
	)
	threeAddressCodes = append(threeAddressCodes, programCodes...)
	threeAddressCodes = append(
		threeAddressCodes,
		fmt.Sprintf("goto %v", whileGoto),
		fmt.Sprintf("%v:", nextGoto),
	)
	return threeAddressCodes, identifiers, nil
}

type OutputStatementAST struct {
	Arguments []ExpressionAST
}

func (o OutputStatementAST) PerformChecks(identifiers []IdentifierInformation) error {
	if len(o.Arguments) == 0 {
		return errors.New("output should always have the first argument")
	}
	datatype, err := o.Arguments[0].GetDatatype(identifiers)
	if err != nil {
		return err
	}
	if !datatype.IsDatatype(StringDatatype{}) {
		return errors.New("output should always have the first argument as a string")
	}
	for _, output := range o.Arguments[1:] {
		_, err = output.GetDatatype(identifiers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o OutputStatementAST) ThreeAddressCode(
	identifiers []IdentifierInformation,
	numberOfGotos *int,
) ([]string, []IdentifierInformation, error) {
	threeAddressCodes := []string{}
	parameters := []string{}

	for _, argument := range o.Arguments {
		param, codes, ids, err := argument.ThreeAddressCode(identifiers)
		if err != nil {
			return []string{}, identifiers, err
		}
		identifiers = ids
		threeAddressCodes = append(threeAddressCodes, codes...)
		parameters = append(parameters, fmt.Sprintf("param %v", param))
	}

	threeAddressCodes = append(threeAddressCodes, parameters...)
	threeAddressCodes = append(threeAddressCodes, "call printf")
	return threeAddressCodes, identifiers, nil
}

type ExpressionAST interface {
	GetDatatype(identifiers []IdentifierInformation) (Datatype, error)
	ThreeAddressCode(
		identifiers []IdentifierInformation,
	) (string, []string, []IdentifierInformation, error)
}

type UnaryExpression struct {
	Operator UnaryOperatorNode
	Operand  ExpressionAST
}

type UnaryOperatorNode int

const (
	UnaryMinus UnaryOperatorNode = iota + 1
	UnaryNot
)

var nameWithUnaryOperator = map[UnaryOperatorNode]string{
	UnaryMinus: "-",
	UnaryNot:   "!",
}

func (u UnaryExpression) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	operandDatatype, err := u.Operand.GetDatatype(identifiers)
	if err != nil {
		return nil, err
	}
	return operandDatatype.PerformUnaryOperation(u.Operator)
}

func (u UnaryExpression) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	var label string

	datatype, err := u.GetDatatype(identifiers)
	if err != nil {
		return "", []string{}, identifiers, err
	}

	label, identifiers = nextIdentifier(identifiers, datatype)
	result, threeAddressCodes, identifiers, err := u.Operand.ThreeAddressCode(
		identifiers,
	)
	if err != nil {
		return "", []string{}, identifiers, err
	}
	threeAddressCodes = append(threeAddressCodes, fmt.Sprintf(
		"%v = %v %v",
		label,
		nameWithUnaryOperator[u.Operator],
		result,
	))

	return label, threeAddressCodes, identifiers, nil
}

type BinaryExpression struct {
	Operator      BinaryOperatorNode
	FirstOperand  ExpressionAST
	SecondOperand ExpressionAST
}

type BinaryOperatorNode int

const (
	BinaryPlus BinaryOperatorNode = iota + 1
	BinaryMinus
	BinaryMul
	BinaryDiv
	BinaryModulo

	BinaryRelationalEquals
	BinaryRelationalGreaterThan
	BinaryRelationalGreaterThanOrEquals
	BinaryRelationalNotEquals
	BinaryRelationalLesserThan
	BinaryRelationalLesserThanOrEquals

	BinaryOr
	BinaryAnd
)

var nameWithBinaryOperator = map[BinaryOperatorNode]string{
	BinaryPlus:   "+",
	BinaryMinus:  "-",
	BinaryMul:    "*",
	BinaryDiv:    "/",
	BinaryModulo: "%",

	BinaryRelationalEquals:              "==",
	BinaryRelationalGreaterThan:         ">",
	BinaryRelationalGreaterThanOrEquals: ">=",
	BinaryRelationalNotEquals:           "!=",
	BinaryRelationalLesserThan:          "<",
	BinaryRelationalLesserThanOrEquals:  "<=",

	BinaryOr:  "||",
	BinaryAnd: "&&",
}

func (b BinaryExpression) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	firstOperandDatatype, err := b.FirstOperand.GetDatatype(identifiers)
	if err != nil {
		return nil, err
	}
	secondOperandDatatype, err := b.SecondOperand.GetDatatype(identifiers)
	if err != nil {
		return nil, err
	}
	return firstOperandDatatype.PerformBinaryOperation(b.Operator, secondOperandDatatype)
}

func (b BinaryExpression) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	datatype, err := b.GetDatatype(identifiers)
	if err != nil {
		return "", []string{}, identifiers, err
	}
	label, identifiers := nextIdentifier(identifiers, datatype)
	threeAddressCodes := []string{}

	firstResult, firstOperandCodes, identifiers, err := b.FirstOperand.ThreeAddressCode(
		identifiers,
	)
	if err != nil {
		return "", []string{}, identifiers, err
	}
	threeAddressCodes = append(threeAddressCodes, firstOperandCodes...)

	secondResult, secondOperandCodes, identifiers, err := b.SecondOperand.ThreeAddressCode(
		identifiers,
	)
	if err != nil {
		return "", []string{}, identifiers, err
	}
	threeAddressCodes = append(threeAddressCodes, secondOperandCodes...)

	threeAddressCodes = append(threeAddressCodes, fmt.Sprintf(
		"%v = %v %v %v",
		label,
		firstResult,
		nameWithBinaryOperator[b.Operator],
		secondResult,
	))

	return label, threeAddressCodes, identifiers, nil
}

type InputExpression struct{}

func (i InputExpression) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	return TypedInt, nil
}

func (i InputExpression) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	label, identifiers := nextIdentifier(identifiers, TypedInt)
	return label, []string{
		fmt.Sprintf("%v = input", label),
	}, identifiers, nil
}

type ArrayExpression struct {
	Elements []ExpressionAST
}

func (a ArrayExpression) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	if len(a.Elements) == 0 {
		return ArrayDatatype{
			ElementType:      TypedUnknown,
			NumberOfElements: 0,
		}, nil
	}
	baseDatatype, err := a.Elements[0].GetDatatype(identifiers)
	if err != nil {
		return nil, err
	}
	for _, element := range a.Elements {
		datatype, err := element.GetDatatype(identifiers)
		if err != nil {
			return nil, err
		}
		if !datatype.IsDatatype(baseDatatype) {
			return nil, errors.New("unmatching datatypes in array")
		}
	}
	return ArrayDatatype{
		ElementType:      baseDatatype,
		NumberOfElements: len(a.Elements),
	}, nil
}

func (a ArrayExpression) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	datatype, err := a.GetDatatype(identifiers)
	if err != nil {
		return "", []string{}, identifiers, err
	}

	elements, err := a.flattenElements(identifiers)
	if err != nil {
		return "", []string{}, identifiers, err
	}

	label, identifiers := nextIdentifier(identifiers, datatype)
	threeAddressCodes := []string{}
	for index, element := range elements {
		result, t, ids, err := element.ThreeAddressCode(identifiers)
		if err != nil {
			return "", []string{}, identifiers, err
		}
		threeAddressCodes = append(threeAddressCodes, t...)
		threeAddressCodes = append(
			threeAddressCodes,
			fmt.Sprintf("%v [] %v = %v", label, index, result),
		)
		identifiers = ids
	}
	return label, threeAddressCodes, identifiers, nil
}

func (a ArrayExpression) flattenElements(
	identifiers []IdentifierInformation,
) ([]ExpressionAST, error) {
	if len(a.Elements) == 0 {
		return []ExpressionAST{}, nil
	}

	datatype, err := a.Elements[0].GetDatatype(identifiers)
	if err != nil {
		return []ExpressionAST{}, err
	}
	arrayDatatype, ok := datatype.(ArrayDatatype)

	elements := a.Elements
	for ok {
		elements, err = a.unwrap(elements)
		if err != nil {
			return []ExpressionAST{}, err
		}
		arrayDatatype, ok = arrayDatatype.ElementType.(ArrayDatatype)
	}
	return elements, nil
}

func (a ArrayExpression) unwrap(elements []ExpressionAST) ([]ExpressionAST, error) {
	result := []ExpressionAST{}
	for _, element := range elements {
		arrayElement, ok := element.(ArrayExpression)
		if !ok {
			return result, errors.New("non-array element in array unwrapping")
		}
		result = append(result, arrayElement.Elements...)
	}
	return result, nil
}

type Identifier struct {
	Id          int
	ArrayValues []ExpressionAST
}

func (i Identifier) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	// e.g., let v = {{1, 2}, {3, 4}, {5, 6}}; v[1] has {1} as ArrayValues and array{int, 2} as Datatype
	if i.Id < 0 || i.Id >= len(identifiers) {
		return nil, errors.New("identifer out-of-bounds")
	}
	baseDatatype := identifiers[i.Id].Datatype

	if len(i.ArrayValues) == 0 {
		return baseDatatype, nil
	}

	arrayDatatype, ok := baseDatatype.(ArrayDatatype)
	if !ok {
		return nil, errors.New("array accesses on a non-array datatype")
	}
	for range i.ArrayValues {
		if arrayDatatype, ok = baseDatatype.(ArrayDatatype); !ok {
			return nil, errors.New("array accesses greater than number of nested arrays")
		}
		baseDatatype = arrayDatatype.ElementType
	}
	return baseDatatype, nil
}

func (i Identifier) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	codes := []string{}
	label := identifierFromIndex(i.Id)
	offset := "0"

	if len(i.ArrayValues) == 0 {
		return label, []string{}, identifiers, nil
	}

	datatype := identifiers[i.Id].Datatype
	arrayDatatype, arrayOk := datatype.(ArrayDatatype)

	for _, access := range i.ArrayValues {
		if !arrayOk {
			return "", []string{}, identifiers, errors.New("mismatching types")
		}
		result, threeAddressCode, identifiersCopy, err := access.ThreeAddressCode(identifiers)
		if err != nil {
			return "", []string{}, identifiers, err
		}
		codes = append(
			codes,
			threeAddressCode...,
		)
		next, identifiersCopy := nextIdentifier(identifiersCopy, TypedInt)
		codes = append(codes, fmt.Sprintf("%v = %v * %v", next, offset, arrayDatatype.NumberOfElements))

		offset, identifiers = nextIdentifier(identifiersCopy, TypedInt)
		codes = append(codes, fmt.Sprintf("%v = %v + %v", offset, next, result))

		datatype = arrayDatatype.ElementType
		arrayDatatype, arrayOk = datatype.(ArrayDatatype)
	}

	result, identifiers := nextIdentifier(identifiers, datatype)
	codes = append(codes, fmt.Sprintf("%v = %v [] %v", result, label, offset))

	return result, codes, identifiers, nil
}

type Literal struct {
	Value    string
	Datatype Datatype
}

func (l Literal) GetDatatype(identifiers []IdentifierInformation) (Datatype, error) {
	return l.Datatype, nil
}

func (l Literal) ThreeAddressCode(
	identifiers []IdentifierInformation,
) (string, []string, []IdentifierInformation, error) {
	return l.Value, []string{}, identifiers, nil
}

func getNextGoto(numberOfGotos *int) string {
	(*numberOfGotos)++
	return fmt.Sprintf("L%d", *numberOfGotos)
}

func nextIdentifier(
	identifiers []IdentifierInformation, datatype Datatype,
) (string, []IdentifierInformation) {
	name := identifierFromIndex(len(identifiers))
	identifiers = append(identifiers, IdentifierInformation{
		IdentifierName: name,
		Datatype:       datatype,
	})
	// since name is the same as identifier code
	return name, identifiers
}

func identifierFromIndex(i int) string {
	return fmt.Sprintf("_t%d", i)
}
