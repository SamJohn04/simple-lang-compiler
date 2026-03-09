package common

// so that more datatypes may be added without worry

type Datatype interface {
	IsDatatype(datatype Datatype) bool
	PerformUnaryOperation(operator UnaryOperatorNode) (Datatype, error)
	PerformBinaryOperation(operator BinaryOperatorNode, with Datatype) (Datatype, error)
	ToString() (string, int, error)
	ToRepresentation() string
}

// this is the type returned by a void function
type VoidDatatype struct{}

func (v VoidDatatype) IsDatatype(datatype Datatype) bool {
	_, ok := datatype.(VoidDatatype)
	return ok
}

func (v VoidDatatype) PerformUnaryOperation(operator UnaryOperatorNode) (Datatype, error) {
	return nil, compilationError("void cannot be an operand")
}

func (v VoidDatatype) PerformBinaryOperation(
	operator BinaryOperatorNode, with Datatype,
) (Datatype, error) {
	return nil, compilationError("void cannot be an operand")
}

func (v VoidDatatype) ToString() (string, int, error) {
	return "", 0, internalError("void cannot be assigned as a datatype")
}

func (v VoidDatatype) ToRepresentation() string {
	return "void"
}

// the basic datatypes: int, bool, char, float, and unknown
type PrimitiveDatatype int

const (
	// when a mutable variable is declared but not initialized
	TypedUnknown PrimitiveDatatype = iota

	TypedInt
	TypedBool
	TypedChar
	TypedFloat
)

func (p PrimitiveDatatype) IsDatatype(datatype Datatype) bool {
	// this check is to prevent another enum type in the future from matching
	primitiveDatatype, ok := datatype.(PrimitiveDatatype)

	return ok && p == primitiveDatatype
}

func (p PrimitiveDatatype) PerformUnaryOperation(operator UnaryOperatorNode) (Datatype, error) {
	switch operator {
	case UnaryMinus:
		switch p {
		case TypedInt:
			return TypedInt, nil
		case TypedFloat:
			return TypedFloat, nil
		default:
			return nil, compilationError("unknown type in unary minus")
		}

	case UnaryNot:
		if p != TypedBool {
			return nil, compilationError("unknown type in unary not")
		}
		return TypedBool, nil

	default:
		return nil, internalError("unknown operation")
	}
}

func (p PrimitiveDatatype) PerformBinaryOperation(
	operator BinaryOperatorNode, with Datatype,
) (Datatype, error) {
	switch datatype := with.(type) {
	case VoidDatatype:
		return nil, compilationError("void cannot be an operand")

	case PrimitiveDatatype:
		return operationOnPrimitives(operator, p, datatype)

	case ArrayDatatype:
		return nil, compilationError("arrays cannot be an operand")

	case StringDatatype:
		return nil, compilationError("string cannot be an operand with a non-string")

	default:
		return nil, internalError("unknown operand datatype")
	}
}

func (p PrimitiveDatatype) ToString() (string, int, error) {
	switch p {
	case TypedUnknown:
		// more correct type checking
		return "int", 0, nil

	case TypedBool:
		return "bool", 1, nil

	case TypedChar:
		return "char", 1, nil

	case TypedInt:
		return "long long", 1, nil

	case TypedFloat:
		return "double", 1, nil

	default:
		return "", 0, internalError("unknown type")
	}
}

func (p PrimitiveDatatype) ToRepresentation() string {
	switch p {
	case TypedUnknown:
		return "u"

	case TypedBool:
		return "b"

	case TypedChar:
		return "c"

	case TypedInt:
		return "l"

	case TypedFloat:
		return "d"

	default:
		return ""
	}
}

func operationOnPrimitives(
	operator BinaryOperatorNode, firstPrimitive PrimitiveDatatype, secondPrimitive PrimitiveDatatype,
) (Datatype, error) {
	switch operator {
	case BinaryPlus:
		fallthrough
	case BinaryMinus:
		fallthrough
	case BinaryMul:
		fallthrough
	case BinaryDiv:
		fallthrough
	case BinaryModulo:
		if firstPrimitive != TypedInt &&
			firstPrimitive != TypedChar &&
			firstPrimitive != TypedFloat {
			return nil, compilationError("unexpected type in mathematical expression")
		} else if secondPrimitive != TypedInt &&
			secondPrimitive != TypedChar &&
			secondPrimitive != TypedFloat {
			return nil, compilationError("unexpected type in mathematical expression")
		}
		if firstPrimitive == TypedFloat || secondPrimitive == TypedFloat {
			return TypedFloat, nil
		}
		return TypedInt, nil

	case BinaryRelationalEquals:
		fallthrough
	case BinaryRelationalNotEquals:
		fallthrough
	case BinaryRelationalGreaterThan:
		fallthrough
	case BinaryRelationalGreaterThanOrEquals:
		fallthrough
	case BinaryRelationalLesserThan:
		fallthrough
	case BinaryRelationalLesserThanOrEquals:
		return TypedBool, nil

	case BinaryAnd:
		fallthrough
	case BinaryOr:
		if firstPrimitive != TypedBool || secondPrimitive != TypedBool {
			return nil, compilationError("unsupported type in logical expression")
		}
		return TypedBool, nil

	default:
		return nil, internalError("unsupported binary expression")
	}
}

type ArrayDatatype struct {
	ElementType      Datatype
	NumberOfElements int
}

func (a ArrayDatatype) IsDatatype(datatype Datatype) bool {
	// an array is said to be the same datatype of another array if and only if
	// - the datatype of its elements are matching
	// - the number of elements are matching
	arrayDatatype, ok := datatype.(ArrayDatatype)

	return ok &&
		a.ElementType.IsDatatype(arrayDatatype.ElementType) &&
		a.NumberOfElements == arrayDatatype.NumberOfElements
}

func (a ArrayDatatype) PerformUnaryOperation(operator UnaryOperatorNode) (Datatype, error) {
	return nil, compilationError("unsupported operation on arrays")
}

func (a ArrayDatatype) PerformBinaryOperation(
	operator BinaryOperatorNode, with Datatype,
) (Datatype, error) {
	return nil, compilationError("unsupported operation on arrays")
}

func (a ArrayDatatype) ToString() (string, int, error) {
	s, childLength, err := a.ElementType.ToString()
	if err != nil {
		return "", 1, err
	}
	length := a.NumberOfElements * childLength
	return s, length, nil
}

func (a ArrayDatatype) ToRepresentation() string {
	return a.ElementType.ToRepresentation()
}

type StringDatatype struct {
	HasKnownLength bool
	CharacterCount int
}

func (s StringDatatype) IsDatatype(datatype Datatype) bool {
	_, ok := datatype.(StringDatatype)
	return ok
}

func (s StringDatatype) PerformUnaryOperation(operator UnaryOperatorNode) (Datatype, error) {
	return nil, compilationError("unsupported operation on strings")
}

func (s StringDatatype) PerformBinaryOperation(
	operator BinaryOperatorNode, with Datatype,
) (Datatype, error) {
	secondStringDatatype, ok := with.(StringDatatype)
	if !ok {
		return nil, compilationError("unsupported operation of string with another type")
	}
	if operator != BinaryPlus {
		return nil, compilationError("unsupported operator with string")
	}
	if !s.HasKnownLength || !secondStringDatatype.HasKnownLength {
		return StringDatatype{
			HasKnownLength: false,
			CharacterCount: -1,
		}, nil
	}
	if s.CharacterCount < 0 || secondStringDatatype.CharacterCount < 0 {
		return nil, internalError("strings have known lengths but they are lesser than 0")
	}
	return StringDatatype{
		HasKnownLength: true,
		CharacterCount: s.CharacterCount + secondStringDatatype.CharacterCount,
	}, nil
}

func (s StringDatatype) ToString() (string, int, error) {
	return "char*", 1, nil
}

func (s StringDatatype) ToRepresentation() string {
	return "str"
}

func compilationError(message string) *CompilationError {
	return &CompilationError{
		PointOfFailure: "types",
		Message:        message,
	}
}

func internalError(message string) *InternalError {
	return &InternalError{
		PointOfFailure: "types",
		Message:        message,
	}
}
